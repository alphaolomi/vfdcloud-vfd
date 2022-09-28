package vfd

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/vfdcloud/vfd/models"
)

var ErrRegistrationFailed = errors.New("registration failed")

type (
	RegistrationResponse struct {
		ACKCODE     string   `xml:"ACKCODE"`
		ACKMSG      string   `xml:"ACKMSG"`
		REGID       string   `xml:"REGID"`
		SERIAL      string   `xml:"SERIAL"`
		UIN         string   `xml:"UIN"`
		TIN         string   `xml:"TIN"`
		VRN         string   `xml:"VRN"`
		MOBILE      string   `xml:"MOBILE"`
		ADDRESS     string   `xml:"ADDRESS"`
		STREET      string   `xml:"STREET"`
		CITY        string   `xml:"CITY"`
		COUNTRY     string   `xml:"COUNTRY"`
		NAME        string   `xml:"NAME"`
		RECEIPTCODE string   `xml:"RECEIPTCODE"`
		REGION      string   `xml:"REGION"`
		ROUTINGKEY  string   `xml:"ROUTINGKEY"`
		GC          int64    `xml:"GC"`
		TAXOFFICE   string   `xml:"TAXOFFICE"`
		USERNAME    string   `xml:"USERNAME"`
		PASSWORD    string   `xml:"PASSWORD"`
		TOKENPATH   string   `xml:"TOKENPATH"`
		TAXCODES    TAXCODES `xml:"TAXCODES"`
	}

	TAXCODES struct {
		XMLName xml.Name `xml:"TAXCODES"`
		Text    string   `xml:",chardata"`
		CODEA   string   `xml:"CODEA"`
		CODEB   string   `xml:"CODEB"`
		CODEC   string   `xml:"CODEC"`
		CODED   string   `xml:"CODED"`
	}
	RegistrationRequest struct {
		ContentType string
		CertSerial  string
		Tin         string
		CertKey     string
	}

	Registrar func(ctx context.Context, url string, privateKey *rsa.PrivateKey,
		request *RegistrationRequest) (*RegistrationResponse, error)

	RegistrationMiddleware func(next Registrar) Registrar
)

func responseFormat(response *models.REGDATARESP) *RegistrationResponse {
	return &RegistrationResponse{
		ACKCODE:     response.ACKCODE,
		ACKMSG:      response.ACKMSG,
		REGID:       response.REGID,
		SERIAL:      response.SERIAL,
		UIN:         response.UIN,
		TIN:         response.TIN,
		VRN:         response.VRN,
		MOBILE:      response.MOBILE,
		ADDRESS:     response.ADDRESS,
		STREET:      response.STREET,
		CITY:        response.CITY,
		COUNTRY:     response.COUNTRY,
		NAME:        response.NAME,
		RECEIPTCODE: response.RECEIPTCODE,
		REGION:      response.REGION,
		ROUTINGKEY:  response.ROUTINGKEY,
		GC:          response.GC,
		TAXOFFICE:   response.TAXOFFICE,
		USERNAME:    response.USERNAME,
		PASSWORD:    response.PASSWORD,
		TOKENPATH:   response.TOKENPATH,
		TAXCODES: TAXCODES{
			CODEA: response.TAXCODES.CODEA,
			CODEB: response.TAXCODES.CODEB,
			CODEC: response.TAXCODES.CODEC,
			CODED: response.TAXCODES.CODED,
		},
	}
}

func wrapRegistrationMiddleware(registrar Registrar, mw ...RegistrationMiddleware) Registrar {
	// Loop backwards through the middleware invoking each one. Replace the
	// registrar with the new wrapped registrar. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		u := mw[i]
		if u != nil {
			registrar = u(registrar)
		}
	}
	return registrar
}

// Register send the registration for a Virtual Fiscal Device to the VFD server. The
// registration request is signed with the private key of the certificate used to
// authenticate the client.
func Register(ctx context.Context, requestURL string, privateKey *rsa.PrivateKey,
	request *RegistrationRequest, mw []RegistrationMiddleware,
) (*RegistrationResponse, error) {
	registrar := func(ctx context.Context, url string, privateKey *rsa.PrivateKey,
		request *RegistrationRequest,
	) (*RegistrationResponse, error) {
		client := getHttpClientInstance().client
		return register(ctx, client, url, privateKey, request)
	}

	registrar = wrapRegistrationMiddleware(registrar, mw...)

	return registrar(ctx, requestURL, privateKey, request)
}

func register(ctx context.Context, client *http.Client, requestURL string, privateKey *rsa.PrivateKey,
	request *RegistrationRequest,
) (*RegistrationResponse, error) {
	var (
		taxIdNumber = request.Tin
		certKey     = request.CertKey
		certSerial  = EncodeBase64String(request.CertSerial)
	)

	reg := models.REGDATA{
		TIN:     taxIdNumber,
		CERTKEY: certKey,
	}

	out, err := xml.Marshal(&reg)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to marshal registration body: %w", ErrRegistrationFailed, err)
	}

	signedPayload, err := Sign(privateKey, out)
	if err != nil {
		return nil, err
	}

	signedPayloadBase64 := EncodeBase64Bytes(signedPayload)
	requestPayload := models.REGDATAEFDMS{
		REGDATA:        reg,
		EFDMSSIGNATURE: signedPayloadBase64,
	}

	out, err = xml.Marshal(&requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(out))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", ContentTypeXML)
	req.Header.Set("Cert-Serial", certSerial)
	req.Header.Set("Client", RegistrationRequestClient)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %v: %w", ErrRegistrationFailed, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "registration: could not close response body %v", err)
		}
	}(resp.Body)

	out, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrRegistrationFailed, err)
	}

	if resp.StatusCode == 500 {
		errBody := models.Error{}
		err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&errBody)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", ErrRegistrationFailed, err)
		}

		return nil, fmt.Errorf("%w: %s", ErrRegistrationFailed, errBody.Message)
	}

	responseBody := models.REGRESPACK{}
	err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrRegistrationFailed, err)
	}

	response := &responseBody.EFDMSRESP

	// check if the response code is equal to zero if not
	// return an error with code and message
	if responseCode := response.ACKCODE; responseCode != "0" {
		responseMessage := response.ACKMSG
		return nil, fmt.Errorf("%v response code: %s, message: %s", ErrRegistrationFailed, responseCode, responseMessage)
	}

	return responseFormat(response), nil
}

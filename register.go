package vfd

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
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
	RegistrationRequest struct {
		URL         string
		ContentType string
		CertSerial  string
		Client      string
		Tin         string
		CertKey     string
		PrivateKey  *rsa.PrivateKey
	}

	TaxCodes struct {
		CodeA string
		CodeB string
		CodeC string
		CodeD string
	}

	RegistrationResponse struct {
		Code        string   `json:"Code,omitempty"`
		Message     string   `json:"Message,omitempty"`
		Id          string   `json:"Id,omitempty"`
		Serial      string   `json:"Serial,omitempty"`
		Uin         string   `json:"Uin,omitempty"`
		Tin         string   `json:"Tin,omitempty"`
		Vrn         string   `json:"Vrn,omitempty"`
		Mobile      string   `json:"Mobile,omitempty"`
		Address     string   `json:"Address,omitempty"`
		Street      string   `json:"Street,omitempty"`
		City        string   `json:"City,omitempty"`
		Country     string   `json:"Country,omitempty"`
		Name        string   `json:"Name,omitempty"`
		ReceiptCode string   `json:"ReceiptCode,omitempty"`
		Region      string   `json:"Region,omitempty"`
		RoutingKey  string   `json:"RoutingKey,omitempty"`
		GC          string   `json:"GC,omitempty"`
		TaxOffice   string   `json:"TaxOffice,omitempty"`
		Username    string   `json:"Username,omitempty"`
		Password    string   `json:"Password,omitempty"`
		TokenPath   string   `json:"TokenPath,omitempty"`
		TaxCodes    TaxCodes `json:"TaxCodes"`
	}
)

func Register(ctx context.Context, request *RegistrationRequest) (*models.RegistrationResponse, error) {
	return register(ctx, getInstance().http, request)
}

func (c *client) Register(ctx context.Context, request *RegistrationRequest) (*models.RegistrationResponse, error) {
	//reg := models.RegistrationBody{
	//	TIN:     request.Tin,
	//	CERTKEY: request.CertKey,
	//}
	//
	//out, err := xml.Marshal(&reg)
	//if err != nil {
	//	return nil, fmt.Errorf("%v: failed to marshal registration body: %w", ErrRegistrationFailed, err)
	//}
	//
	//signedPayload, err := c.SignPayload(ctx, request.PrivateKey, out)
	//if err != nil {
	//	return nil, err
	//}
	//
	//signedPayloadBase64 := base64.StdEncoding.EncodeToString(signedPayload)
	//requestPayload := models.RegistrationRequest{
	//	Reg:            reg,
	//	EFDMSSIGNATURE: signedPayloadBase64,
	//}
	//
	//out, err = xml.Marshal(&requestPayload)
	//if err != nil {
	//	return nil, err
	//}
	//
	//req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(out))
	//if err != nil {
	//	return nil, err
	//}
	//req.Header.Set("Content-Type", request.ContentType)
	//req.Header.Set("Cert-Serial", request.CertSerial)
	//req.Header.Set("client", request.client)
	//
	//resp, err := c.http.StandardClient().Do(req)
	//if err != nil {
	//	return nil, fmt.Errorf("http error: %v: %w", ErrRegistrationFailed, err)
	//}
	//
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		_, _ = fmt.Fprintf(os.Stderr, "registration: could not close response body %v", err)
	//	}
	//}(resp.Body)
	//
	//out, err = io.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("%v: %w", ErrRegistrationFailed, err)
	//}
	//
	//if resp.StatusCode == 500 {
	//	errBody := models.Error{}
	//	err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&errBody)
	//	if err != nil {
	//		return nil, fmt.Errorf("%v: %w", ErrRegistrationFailed, err)
	//	}
	//
	//	return nil, fmt.Errorf("%w: %s", ErrRegistrationFailed, errBody.Message)
	//}
	//
	//responseBody := models.RegistrationAck{}
	//err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&responseBody)
	//if err != nil {
	//	return nil, fmt.Errorf("%v: %w", ErrRegistrationFailed, err)
	//}
	//
	//response := &responseBody.EFDMSRESP
	//
	//// check if the response code is equal to zero if not
	//// return an error with code and message
	//if responseCode := response.ACKCODE; responseCode != "0" {
	//	responseMessage := response.ACKMSG
	//	return nil, fmt.Errorf("%v response code: %s, message: %s", ErrRegistrationFailed, responseCode, responseMessage)
	//}
	//
	//return response, nil

	var (
		client = c.http
	)

	return register(ctx, client, request)
}

func register(ctx context.Context, client *http.Client, request *RegistrationRequest) (*models.RegistrationResponse, error) {
	var (
		requestURL  = request.URL
		taxIdNumber = request.Tin
		certKey     = request.CertKey
		privateKey  = request.PrivateKey
		apiClient   = request.Client
		certSerial  = base64.StdEncoding.EncodeToString([]byte(request.CertSerial))
		contentType = request.ContentType
	)

	reg := models.RegistrationBody{
		TIN:     taxIdNumber,
		CERTKEY: certKey,
	}

	out, err := xml.Marshal(&reg)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to marshal registration body: %w", ErrRegistrationFailed, err)
	}

	signedPayload, err := Sign(ctx, privateKey, out)
	if err != nil {
		return nil, err
	}

	signedPayloadBase64 := base64.StdEncoding.EncodeToString(signedPayload)
	requestPayload := models.RegistrationRequest{
		Reg:            reg,
		EFDMSSIGNATURE: signedPayloadBase64,
	}

	out, err = xml.Marshal(&requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(out))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cert-Serial", certSerial)
	req.Header.Set("Client", apiClient)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %v: %w", ErrRegistrationFailed, err)
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

	responseBody := models.RegistrationAck{}
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

	return response, nil
}

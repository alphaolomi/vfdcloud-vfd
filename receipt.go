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

const (
	CustomerIDTypeTIN CUSTIDTYPE = iota + 1
	CustomerIDTypeDrivingLicense
	CustomerIDTypeVoterID
	CustomerIDTypePassport
	CustomerIDTypeNIDA
	CustomerIDTypeNIL
	CustomerIDTypeMeterNumber
)

const ()

const (
	PMTTYPE_CHEQUE  = "CHEQUE"  // Cheque
	PMTTYPE_CCARD   = "CCARD"   // Credit Card
	PMTTYPE_EMONEY  = "EMONEY"  // Electronic Money
	PMTTYPE_INVOICE = "INVOICE" // Invoice
)

type (

	// CUSTIDTYPE is Buyer Used ID type.
	// CustomerIDTypeTIN (1) - Tax Identification Number,
	// CustomerIDTypeDrivingLicense (2) - Driving License,
	// CustomerIDTypeVoterID (3) - Voters Number,
	// CustomerIDTypePassport (4) - Travel Passport,
	// CustomerIDTypeNIDA (5) - National ID,
	// CustomerIDTypeNIL (6) - NIL (No Identity Used),
	// CustomerIDTypeMeterNumber (7) - Meter Number
	CUSTIDTYPE int

	// PMTTYPE is the payment type used during the payment.
	// 1. CASH 2. CHEQUE 3. CCARD 4. EMONEY and 5. INVOICE
	PMTTYPE string

	// VATRATE is the identifier of the VAT rate used during the payment.
	// Identifier of the Tax rate
	// A= 18 (Standard Rate for VAT items)
	// B= 0 (Special Rate)
	// C= 0 (Zero rated for Non-VAT items)
	// D= 0 (Special Relief for relieved items)
	// E= 0 (Exempt items)
	VATRATE string

	// ReceiptRequest is the request body for posting the receipt
	// ContentType is always Application/xml
	// RoutingKey vfdrct
	// CertSerial is the serial of the Key certificate provided. Should be base64 encoded.
	// BearerToken bearer <token_value>
	// i.e. it should start with word bearer followed by space then followed by current token value
	ReceiptRequest struct {
		URL         string
		ContentType string
		CertSerial  string
		BearerToken string
		RoutingKey  string
		PrivateKey  *rsa.PrivateKey
	}
)

func (idType CUSTIDTYPE) String() string {
	return [...]string{"1", "2", "3", "4", "5", "6", "7"}[idType-1]
}

var ErrReceiptUploadFailed = errors.New("receipt upload failed")

func (c *httpx) UploadReceipt(ctx context.Context, request *ReceiptRequest, rct *models.RCT) (*Response, error) {
	client := c.client

	return receiptUpload(ctx, client, request, rct)
}

func ReceiptUpload(ctx context.Context, request *ReceiptRequest, rct *models.RCT) (*Response, error) {
	client := http.DefaultClient
	return receiptUpload(ctx, client, request, rct)
}

func receiptUpload(ctx context.Context, client *http.Client, request *ReceiptRequest, rct *models.RCT) (*Response, error) {
	var (
		requestURL  = request.URL
		privateKey  = request.PrivateKey
		contentType = request.ContentType
		routingKey  = request.RoutingKey
		certSerial  = request.CertSerial
		bearerToken = request.BearerToken
	)

	// print certSerial
	// FIXME: remove this line
	fmt.Sprintf("\n\n\ncertSerial(receipt upload): %s\n\n\n", certSerial)

	//log.F(xio.Stdout, "VFD_GATEWAY",
	//	"TRACE", log.LevelDebug, "UploadReceipt: ", fmt.Sprintf("%+v", request),
	//)
	pctx, cancel := context.WithCancel(ctx)
	defer cancel()

	out, err := xml.Marshal(rct)
	if err != nil {
		return nil, err
	}

	//log.F(xio.Stdout, "VFD_GATEWAY", "TRACE", log.LevelDebug, fmt.Sprintf("UploadReceipt: %s", string(out)))

	signedPayload, err := Sign(ctx, privateKey, out)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}
	signedPayloadBase64 := base64.StdEncoding.EncodeToString(signedPayload)
	requestPayload := models.RCTEFDMS{
		RCT:            *rct,
		EFDMSSIGNATURE: signedPayloadBase64,
	}

	out, err = xml.Marshal(&requestPayload)
	if err != nil {
		return nil, err
	}
	//outHeader := []byte(xml.Header + string(out))

	//log.F(xio.Stdout, "VFD_GATEWAY", "TRACE", log.LevelDebug, fmt.Sprintf("UploadReceipt: %s", string(outHeader)))
	req, err := http.NewRequestWithContext(pctx, http.MethodPost, requestURL, bytes.NewBuffer(out))
	if err != nil {
		return nil, err
	}

	/// print out
	// Filename should be {DATE}-{GC}-{DC}-receipt.xml
	// They are stored under receipts folder in the home directory

	// create receipt folder if not exists in the home directory
	if _, err := os.Stat("./receipts"); os.IsNotExist(err) {
		err = os.Mkdir("./receipts", 0755)
		if err != nil {
			return nil, err
		}
	}
	fileName := fmt.Sprintf("%s-%s-%s-receipt.xml", rct.DATE, rct.GC, rct.DC)
	filePath := fmt.Sprintf("./receipts/%s", fileName)
	f, _ := os.Create(filePath)
	_, _ = f.Write(out)

	// Header setting
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Routing-Key", routingKey)
	req.Header.Set("Cert-Serial", certSerial)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	//	log.F(xio.Stdout, "VFD_GATEWAY", "TRACE", log.LevelDebug, fmt.Sprintf("\n\nHeaders: %+v\n\n", req.Header))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "registration: could not close response body %v", err)
		}
	}(resp.Body)

	out, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}

	// Log the response
	//log.F(xio.Stdout, "VFD_GATEWAY", "TRACE", log.LevelDebug, "registration response, status: ", resp.Status, "headers", resp.Header, "body", string(out))

	if resp.StatusCode == 500 {
		errBody := models.Error{}
		err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&errBody)
		if err != nil {
			return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
		}

		return nil, fmt.Errorf("registration error: %s", errBody.Message)
	}

	response := models.RCTACKEFDMS{}
	err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}

	return &Response{
		Number:  response.RCTACK.RCTNUM,
		Date:    response.RCTACK.DATE,
		Time:    response.RCTACK.TIME,
		Code:    response.RCTACK.ACKCODE,
		Message: response.RCTACK.ACKMSG,
	}, nil
}

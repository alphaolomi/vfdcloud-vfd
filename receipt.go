package vfd

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/vfdcloud/vfd/models"
	"io"
	"net/http"
	"os"
)

var ErrReceiptUploadFailed = errors.New("receipt upload failed")

type (
	// ReceiptUploader uploads receipts to the VFD server
	ReceiptUploader func(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
		receipt *models.RCT) (*Response, error)

	ReceiptUploadMiddleware func(next ReceiptUploader) ReceiptUploader
)

func VerifyUploadReceiptRequest() ReceiptUploadMiddleware {
	m := func(next ReceiptUploader) ReceiptUploader {
		u := func(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
			receipt *models.RCT) (*Response, error) {

			// Steps:
			// TODO 1. Verify the request headers
			// TODO 2. verify request URL
			// TODO 3. Verify the receipt
			return next(ctx, url, headers, privateKey, receipt)
		}
		return u
	}
	return m
}

// UploadReceipt uploads a receipt to the VFD server.
func UploadReceipt(ctx context.Context, requestURL string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	rct *models.RCT, mw ...ReceiptUploadMiddleware) (*Response, error) {
	client := httpClientInstance().client
	uploader := func(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
		receipt *models.RCT) (*Response, error) {
		return uploadReceipt(ctx, client, url, headers, privateKey, receipt)
	}
	uploader = wrapReceiptUploaderMiddleware(uploader, VerifyUploadReceiptRequest())
	uploader = wrapReceiptUploaderMiddleware(uploader, mw...)
	return uploader(ctx, requestURL, headers, privateKey, rct)
}

func wrapReceiptUploaderMiddleware(uploader ReceiptUploader, mw ...ReceiptUploadMiddleware,
) ReceiptUploader {
	// Loop backwards through the middleware invoking each one. Replace the
	// fetcher with the new wrapped fetcher. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		u := mw[i]
		if u != nil {
			uploader = u(uploader)
		}
	}

	return uploader
}

func uploadReceipt(ctx context.Context, client *http.Client, requestURL string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	rct *models.RCT) (*Response, error) {
	var (
		contentType = headers.ContentType
		routingKey  = headers.RoutingKey
		certSerial  = headers.CertSerial
		bearerToken = headers.BearerToken
	)

	newContext, cancel := context.WithCancel(ctx)
	defer cancel()

	out, err := xml.Marshal(rct)
	if err != nil {
		return nil, err
	}

	signedPayload, err := Sign(ctx, privateKey, out)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}
	signedPayloadBase64 := EncodeBase64Bytes(signedPayload)
	requestPayload := models.RCTEFDMS{
		RCT:            *rct,
		EFDMSSIGNATURE: signedPayloadBase64,
	}

	out, err = xml.Marshal(&requestPayload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(newContext, http.MethodPost, requestURL, bytes.NewBuffer(out))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Routing-Key", routingKey)
	req.Header.Set("Cert-Serial", certSerial)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

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

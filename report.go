package vfd

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/xml"
	"fmt"
	"github.com/vfdcloud/vfd/models"
	"io"
	"net/http"
	"os"
)

var ErrReportSubmitFailed = fmt.Errorf("report submit failed")

type (
	// ReportTotals contains different number of totals
	ReportTotals struct {
		DailyTotalAmount float64
		Gross            float64
		Corrections      float64
		Discounts        float64
		Surcharges       float64
		TicketsVoid      float64
		TicketsVoidTotal int64
		TicketsFiscal    int64
		TicketsNonFiscal int64
	}

	ReportSubmitter func(ctx context.Context, url string, headers *RequestHeaders,
		privateKey *rsa.PrivateKey,
		report *models.Report) (*Response, error)

	ReportSubmitMiddleware func(next ReportSubmitter) ReportSubmitter
)

// wrapReportSubmitterMiddleware wraps a ReportSubmitter with a list of ReportSubmitMiddleware.
func wrapReportSubmitterMiddleware(submitter ReportSubmitter, mw ...ReportSubmitMiddleware,
) ReportSubmitter {
	// Loop backwards through the middleware invoking each one. Replace the
	// submitter with the new wrapped submitter. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		submitter = mw[i](submitter)
	}
	return submitter
}

// submitReport submits a report to the VFD server.
func submitReport(ctx context.Context, client *http.Client, requestURL string, headers *RequestHeaders,
	privateKey *rsa.PrivateKey,
	report *models.Report) (*Response, error) {
	var (
		contentType = headers.ContentType
		routingKey  = headers.RoutingKey
		certSerial  = headers.CertSerial
		bearerToken = headers.BearerToken
	)

	newContext, cancel := context.WithCancel(ctx)
	defer cancel()

	out, err := xml.Marshal(report)
	if err != nil {
		return nil, err
	}

	signedPayload, err := Sign(ctx, privateKey, out)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}
	signedPayloadBase64 := EncodeBase64Bytes(signedPayload)
	requestPayload := models.ReportEFDMS{
		ZREPORT:        *report,
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
		return nil, fmt.Errorf("%v : %w", ErrReportSubmitFailed, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close the body %v", err)
		}
	}(resp.Body)

	out, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReportSubmitFailed, err)
	}

	if resp.StatusCode == 500 {
		errBody := models.Error{}
		err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&errBody)
		if err != nil {
			return nil, fmt.Errorf("%v : %w", ErrReportSubmitFailed, err)
		}

		return nil, fmt.Errorf("registration error: %s", errBody.Message)
	}

	response := models.ReportAckEFDMS{}
	err = xml.NewDecoder(bytes.NewBuffer(out)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReportSubmitFailed, err)
	}

	return &Response{
		Number:  response.ZACK.ZNUMBER,
		Date:    response.ZACK.DATE,
		Time:    response.ZACK.TIME,
		Code:    response.ZACK.ACKCODE,
		Message: response.ZACK.ACKMSG,
	}, nil
}

func SubmitReport(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	report *models.Report, mw ...ReportSubmitMiddleware) (*Response, error) {
	client := httpClientInstance().client
	submitter := func(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
		report *models.Report) (*Response, error) {
		return submitReport(ctx, client, url, headers, privateKey, report)
	}
	submitter = wrapReportSubmitterMiddleware(submitter, mw...)
	return submitReport(ctx, client, url, headers, privateKey, report)
}

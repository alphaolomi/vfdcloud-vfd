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
	"strings"
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
		TicketsVoid      int64
		TicketsVoidTotal float64
		TicketsFiscal    int64
		TicketsNonFiscal int64
	}

	Address struct {
		Name    string
		Street  string
		Mobile  string
		City    string
		Country string
	}

	ReportParams struct {
		Date             string
		Time             string
		VRN              string
		TIN              string
		TaxOffice        string
		RegistrationID   string
		ZNumber          string
		EFDSerial        string
		RegistrationDate string
		User             string
		SIMIMSI          string
		VATChangeNum     int64
		HeadChangeNum    int64
		FirmwareVersion  string
		FirmwareChecksum string
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

func (lines *Address) AsList() []string {
	return []string{
		strings.ToUpper(lines.Name),
		strings.ToUpper(lines.Street),
		fmt.Sprintf("MOBILE: %s", lines.Mobile),
		strings.ToUpper(fmt.Sprintf("%s,%s", lines.City, lines.Country))}
}

func GenerateZReport(params *ReportParams, address Address, vats []VatTotal, payments []Payment, totals ReportTotals) *models.ZREPORT {
	const (
		SIMIMSI       = "WEBAPI"
		FWVERSION     = "3.0"
		FWCHECKSUM    = "WEBAPI"
		VATCHANGENUM  = "0"
		HEADCHANGENUM = "0"
		ERRORS        = ""
	)

	payments1 := make([]*models.PAYMENT, len(payments))
	for _, p := range payments {
		payments1 = append(payments1, &models.PAYMENT{
			PMTTYPE:   string(p.Type),
			PMTAMOUNT: p.Amount,
		})
	}

	vats1 := make([]*models.VATTOTAL, len(vats))
	for _, v := range vats {
		taxAmount := v.Rate * v.Amount
		rate := fmt.Sprintf("%s-%.2f", v.ID, v.Rate)
		vats1 = append(vats1, &models.VATTOTAL{
			VATRATE:    rate,
			NETTAMOUNT: v.Amount,
			TAXAMOUNT:  taxAmount,
		})
	}

	TT := models.REPORTTOTALS{
		DAILYTOTALAMOUNT: totals.DailyTotalAmount,
		GROSS:            totals.Gross,
		CORRECTIONS:      totals.Corrections,
		DISCOUNTS:        totals.Discounts,
		SURCHARGES:       totals.Surcharges,
		TICKETSVOID:      totals.TicketsVoid,
		TICKETSVOIDTOTAL: totals.TicketsVoidTotal,
		TICKETSFISCAL:    totals.TicketsFiscal,
		TICKETSNONFISCAL: totals.TicketsNonFiscal,
	}
	return &models.ZREPORT{
		XMLName: xml.Name{},
		Text:    "",
		DATE:    params.Date,
		TIME:    params.Time,
		HEADER: struct {
			Text string   `xml:",chardata"`
			LINE []string `xml:"LINE"`
		}{
			LINE: address.AsList(),
		},
		VRN:              params.VRN,
		TIN:              params.TIN,
		TAXOFFICE:        params.TaxOffice,
		REGID:            params.RegistrationID,
		ZNUMBER:          params.ZNumber,
		EFDSERIAL:        params.EFDSerial,
		REGISTRATIONDATE: params.RegistrationDate,
		USER:             params.User,
		SIMIMSI:          SIMIMSI,
		TOTALS:           TT,
		VATTOTALS: models.VATTOTALS{
			VATTOTAL: vats1,
		},
		PAYMENTS: models.PAYMENTS{
			PAYMENT: payments1,
		},
		CHANGES: struct {
			Text          string `xml:",chardata"`
			VATCHANGENUM  string `xml:"VATCHANGENUM"`
			HEADCHANGENUM string `xml:"HEADCHANGENUM"`
		}{
			VATCHANGENUM:  VATCHANGENUM,
			HEADCHANGENUM: HEADCHANGENUM,
		},
		ERRORS:     ERRORS,
		FWVERSION:  FWVERSION,
		FWCHECKSUM: FWCHECKSUM,
	}
}

// ReportPayloadBytes returns the bytes of the report payload. It calls xml.Marshal on the report.
// then replace all the occurrences of <PAYMENT>, </PAYMENT>, <VATTOTAL>, </VATTOTAL> with empty string ""
// and then add the xml.Header to the beginning of the payload.
func ReportPayloadBytes(params *ReportParams, address Address, vats []VatTotal, payments []Payment,
	totals ReportTotals) ([]byte, error) {
	replaceList := []string{"<PAYMENT>", "", "</PAYMENT>", "", "<VATTOTAL>", "", "</VATTOTAL>", ""}
	replacer := strings.NewReplacer(replaceList...)
	zReport := GenerateZReport(params, address, vats, payments, totals)
	payload, err := xml.Marshal(zReport)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the report: %w", err)
	}
	payload = []byte(replacer.Replace(string(payload)))
	payload = append([]byte(xml.Header), payload...)
	return payload, nil
}

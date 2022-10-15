package vfd

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/vfdcloud/vfd/models"
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
	}

	ReportRequest struct {
		Params  *ReportParams
		Address *Address
		Totals  *ReportTotals
		VATS    []VATTOTAL
		Payment []Payment
	}

	ReportSubmitter func(ctx context.Context, url string, headers *RequestHeaders,
		privateKey *rsa.PrivateKey,
		report *ReportRequest) (*Response, error)

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
	report *ReportRequest,
) (*Response, error) {
	var (
		certSerial  = headers.CertSerial
		bearerToken = headers.BearerToken
	)

	newContext, cancel := context.WithCancel(ctx)
	defer cancel()

	payload, err := ReportBytes(
		privateKey, report.Params, *report.Address, report.VATS,
		report.Payment, *report.Totals)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the report payload: %w", err)
	}

	req, err := http.NewRequestWithContext(newContext, http.MethodPost, requestURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", ContentTypeXML)
	req.Header.Set("Routing-Key", SubmitReportRoutingKey)
	req.Header.Set("Cert-Serial", EncodeBase64String(certSerial))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReportSubmitFailed, err)
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReportSubmitFailed, err)
	}

	if resp.StatusCode == http.StatusInternalServerError {
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
	report *ReportRequest, mw ...ReportSubmitMiddleware,
) (*Response, error) {
	client := getHttpClientInstance().client
	submitter := func(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
		report *ReportRequest,
	) (*Response, error) {
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
		strings.ToUpper(fmt.Sprintf("%s,%s", lines.City, lines.Country)),
	}
}

func GenerateZReport(params *ReportParams, address Address, vats []VATTOTAL, payments []Payment, totals ReportTotals) *models.ZREPORT {
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

	PAYMENTS := models.PAYMENTS{PAYMENT: payments1}

	vats1 := make([]*models.VATTOTAL, len(vats))
	for _, v := range vats {
		rate := fmt.Sprintf("%s-%.2f", v.ID, v.Rate)
		vats1 = append(vats1, &models.VATTOTAL{
			VATRATE:    rate,
			NETTAMOUNT: v.NetAmount,
			TAXAMOUNT:  v.TaxAmount,
		})
	}

	VATTOTALS := models.VATTOTALS{VATTOTAL: vats1}

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
	report := &models.ZREPORT{
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
		USER:             "",
		SIMIMSI:          SIMIMSI,
		TOTALS:           TT,
		VATTOTALS:        VATTOTALS,
		PAYMENTS:         PAYMENTS,
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

	report.RoundOff()

	return report
}

// ReportBytes returns the bytes of the report payload. It calls xml.Marshal on the report.
// then replace all the occurrences of <PAYMENT>, </PAYMENT>, <VATTOTAL>, </VATTOTAL> with empty string ""
// and then add the xml.Header to the beginning of the payload.
func ReportBytes(privateKey *rsa.PrivateKey, params *ReportParams, address Address,
	vats []VATTOTAL, payments []Payment,
	totals ReportTotals,
) ([]byte, error) {
	zReport := GenerateZReport(params, address, vats, payments, totals)
	payload, err := xml.Marshal(zReport)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the report: %w", err)
	}

	payloadString := formatReportXmlPayload(payload, totals, vats, payments)

	signedPayload, err := SignPayload(privateKey, []byte(payloadString))
	if err != nil {
		return nil, fmt.Errorf("failed to sign the payload: %w", err)
	}
	base64PayloadSignature := base64.StdEncoding.EncodeToString(signedPayload)
	report := fmt.Sprintf("<EFDMS>%s<EFDMSSIGNATURE>%s</EFDMSSIGNATURE></EFDMS>", payloadString, base64PayloadSignature)
	report = fmt.Sprintf("%s%s", xml.Header, report)

	return []byte(report), nil
}

func formatReportXmlPayload(payload []byte, totals ReportTotals, vats []VATTOTAL, payments []Payment) string {
	replaceList := []string{"<PAYMENT>", "", "</PAYMENT>", "", "<VATTOTAL>", "", "</VATTOTAL>", ""}
	replacer := strings.NewReplacer(replaceList...)
	payloadString := replacer.Replace(string(payload))
	dailyAmountTag := fmt.Sprintf("<DAILYTOTALAMOUNT>%.2f</DAILYTOTALAMOUNT>", totals.DailyTotalAmount)
	grossAmountTag := fmt.Sprintf("<GROSS>%.2f</GROSS>", totals.Gross)

	netAmountTag := fmt.Sprintf("<NETTAMOUNT>%.2f</NETTAMOUNT>", vats[0].TaxAmount)
	pmtAmountTag := fmt.Sprintf("<PMTAMOUNT>%.2f</PMTAMOUNT>", payments[0].Amount)
	taxAmountTag := fmt.Sprintf("<TAXAMOUNT>%.2f</TAXAMOUNT>", vats[0].TaxAmount)
	regexDailyAmount := regexp.MustCompile(`<DAILYTOTALAMOUNT>.*</DAILYTOTALAMOUNT>`)
	regexGrossAmount := regexp.MustCompile(`<GROSS>.*</GROSS>`)
	regexPmtAmount := regexp.MustCompile(`<PMTAMOUNT>.*</PMTAMOUNT>`)
	regexNetAmount := regexp.MustCompile(`<NETTAMOUNT>.*</NETTAMOUNT>`)
	regexTaxAmount := regexp.MustCompile(`<TAXAMOUNT>.*</TAXAMOUNT>`)
	// replace all the occurrences of the regex with the correct value
	payloadString = regexDailyAmount.ReplaceAllString(payloadString, dailyAmountTag)
	payloadString = regexGrossAmount.ReplaceAllString(payloadString, grossAmountTag)
	payloadString = regexPmtAmount.ReplaceAllString(payloadString, pmtAmountTag)
	payloadString = regexNetAmount.ReplaceAllString(payloadString, netAmountTag)
	payloadString = regexTaxAmount.ReplaceAllString(payloadString, taxAmountTag)

	return payloadString
}

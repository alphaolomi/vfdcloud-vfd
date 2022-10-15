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
	"strings"

	"github.com/vfdcloud/base"
	"github.com/vfdcloud/vfd/internal/models"
)

var ErrReceiptUploadFailed = errors.New("receipt upload failed")

const (
	NonTaxableItemCode    = 3
	TaxableItemCode       = 1
	NonTaxableItemVatCode = "C"
	TaxableItemVatCode    = "A"
)

type (
	// ReceiptParams contains parameters icluded while sending the receipts
	ReceiptParams struct {
		Date           string
		Time           string
		TIN            string
		RegistrationID string
		EFDSerial      string
		ReceiptNum     string
		DailyCounter   int64
		GlobalCounter  int64
		ZNum           string
		ReceiptVNum    string
	}

	// Customer contains customer information
	Customer struct {
		Type   CustomerID
		ID     string
		Name   string
		Mobile string
	}

	// Item represent a purchased item. TaxCode is an integer that can take the
	// value of 1 for taxable items and 3 for non-taxable items.
	Item struct {
		ID          string
		Description string
		TaxCode     int64
		Quantity    float64
		Price       float64
		Discount    float64
	}

	ReceiptRequest struct {
		Params   ReceiptParams
		Customer Customer
		Items    []Item
		Payments []Payment
	}
	//// ReceiptSubmitter uploads receipts to the VFD server
	//ReceiptSubmitter func(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	//	receipt *ReceiptRequest) (*Response, error)
	//
	//ReceiptSubmitMiddleware func(next ReceiptSubmitter) ReceiptSubmitter
)

// SubmitReceipt uploads a receipt to the VFD server.
func SubmitReceipt(ctx context.Context, requestURL string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	receiptRequest *ReceiptRequest,
) (*Response, error) {
	client := getHttpClientInstance().client
	return submitReceipt(ctx, client, requestURL, headers, privateKey, receiptRequest)
}

//func wrapReceiptSubmitMiddlewares(uploader ReceiptSubmitter, mw ...ReceiptSubmitMiddleware) ReceiptSubmitter {
//	// Loop backwards through the middleware invoking each one. Replace the
//	// fetcher with the new wrapped fetcher. Looping backwards ensures that the
//	// first middleware of the slice is the first to be executed by requests.
//	for i := len(mw) - 1; i >= 0; i-- {
//		u := mw[i]
//		if u != nil {
//			uploader = u(uploader)
//		}
//	}
//
//	return uploader
//}

func submitReceipt(ctx context.Context, client *http.Client, requestURL string, headers *RequestHeaders,
	privateKey *rsa.PrivateKey, rct *ReceiptRequest,
) (*Response, error) {
	var (
		certSerial  = headers.CertSerial
		bearerToken = headers.BearerToken
	)

	newContext, cancel := context.WithCancel(ctx)
	defer cancel()

	payload, err := ReceiptBytes(
		privateKey, rct.Params, rct.Customer, rct.Items, rct.Payments)
	if err != nil {
		return nil, fmt.Errorf("%v : %w", ErrReceiptUploadFailed, err)
	}

	req, err := http.NewRequestWithContext(newContext, http.MethodPost, requestURL,
		bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", ContentTypeXML)
	req.Header.Set("Routing-Key", SubmitReceiptRoutingKey)
	req.Header.Set("Cert-Serial", encodeBase64String(certSerial))
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

	out, err := io.ReadAll(resp.Body)
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

func generateReceipt(params ReceiptParams, customer Customer, items []Item, payments []Payment) *models.RCT {
	rctPayments := make([]*models.PAYMENT, len(payments))
	for i, payment := range payments {
		rctPayments[i] = &models.PAYMENT{
			PMTTYPE:   string(payment.Type),
			PMTAMOUNT: fmt.Sprintf("%.2f", payment.Amount),
		}
	}

	RESULTS := processItems(items)
	ITEMS := models.ITEMS{ITEM: RESULTS.ITEMS}
	TOTALS := RESULTS.TOTALS
	VATTOTALS := models.VATTOTALS{VATTOTAL: RESULTS.VATTOTALS}
	PAYMENTS := models.PAYMENTS{PAYMENT: rctPayments}

	RECEIPT := &models.RCT{
		DATE:       params.Date,
		TIME:       params.Time,
		TIN:        params.TIN,
		REGID:      params.RegistrationID,
		EFDSERIAL:  params.EFDSerial,
		CUSTIDTYPE: int64(customer.Type),
		CUSTID:     customer.ID,
		CUSTNAME:   customer.Name,
		MOBILENUM:  customer.Mobile,
		RCTNUM:     params.ReceiptNum,
		DC:         params.DailyCounter,
		GC:         params.GlobalCounter,
		ZNUM:       params.ZNum,
		RCTVNUM:    params.ReceiptVNum,
		ITEMS:      ITEMS,
		TOTALS:     TOTALS,
		PAYMENTS:   PAYMENTS,
		VATTOTALS:  VATTOTALS,
	}

	// round off all values to 2 decimal places
	RECEIPT.RoundOff()

	return RECEIPT
}

func ReceiptBytes(privateKey *rsa.PrivateKey, params ReceiptParams, customer Customer,
	items []Item, payments []Payment,
) ([]byte, error) {
	receipt := generateReceipt(params, customer, items, payments)
	receiptBytes, err := xml.Marshal(receipt)
	if err != nil {
		return nil, fmt.Errorf("could not marshal receipt: %w", err)
	}
	replacer := strings.NewReplacer(
		"<PAYMENT>", "",
		"</PAYMENT>", "",
		"<VATTOTAL>", "",
		"</VATTOTAL>", "")

	receiptBytes = []byte(replacer.Replace(string(receiptBytes)))
	signedReceipt, err := Sign(privateKey, receiptBytes)
	if err != nil {
		return nil, fmt.Errorf("could not sign receipt: %w", err)
	}
	base64SignedReceipt := encodeBase64Bytes(signedReceipt)
	receiptString := string(receiptBytes)

	report := fmt.Sprintf("<EFDMS>%s<EFDMSSIGNATURE>%s</EFDMSSIGNATURE></EFDMS>", receiptString, base64SignedReceipt)
	report = fmt.Sprintf("%s%s", xml.Header, report)

	return []byte(report), nil
}

// ReceiptLink creates a link to the receipt it accepts RECEIPTCODE, GC and the RECEIPTTIME
// and base.Env to know if the receipt was created during testing or production.
func ReceiptLink(e base.Env, receiptCode string, gc int64, receiptTime string) string {
	var baseURL string

	if e == base.ProdEnv {
		baseURL = VerifyReceiptProductionURL
	} else {
		baseURL = VerifyReceiptTestingURL
	}
	return receiptLink(baseURL, receiptCode, gc, receiptTime)
}

func receiptLink(baseURL string, receiptCode string, gc int64, receiptTime string) string {
	return fmt.Sprintf(
		"%s%s%d_%s",
		baseURL,
		receiptCode,
		gc,
		strings.ReplaceAll(receiptTime, ":", ""))
}

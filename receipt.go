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

const (
	NonTaxableItemCode    = 3
	TaxableItemCode       = 1
	NonTaxableItemVatCode = "C"
	TaxableItemVatCode    = "A"
)

type (
	// ReceiptParams contains parameters icluded while sending the receipts
	//<DATE>2019-08-27</DATE>
	//		<!-- Date of issue of receipt/invoice in the format YYYY-MM-DD-->
	//		<TIME>08:36:02</TIME>
	//		<!-- Time of issue of receipt/invoice in the format HH:MM:SS-->
	//		<TIN>222222286</TIN>
	//		<!-- Tin of the taxpayer -->
	//		<REGID>TZ0100089</REGID>
	//		<!-- Registration number of the VFD -->
	//		<EFDSERIAL>10TZ1000211</EFDSERIAL>
	//		<!-- Serial number of VFD -->
	//		<CUSTIDTYPE>1</CUSTIDTYPE>
	//		<!-- Customer ID type, values range from 1 to 6 as specified in API document-->
	//		<CUSTID>111111111</CUSTID>
	//		<!-- Customer ID based on CUSTIDTYPE specified above-->
	//		<CUSTNAME></CUSTNAME>
	//		<!-- Customer name-->
	//		<MOBILENUM></MOBILENUM>
	//		<!-- Customer mobile number-->
	//		<RCTNUM>380</RCTNUM>
	//		<!-- A receipt/invoice number which is same as GC. It should compose of digits alone i.e. without letters-->
	//		<DC>1</DC>
	//		<!-- Daily counter of recipt/invoice which increments for each receipt/invoice and reset to 1 on a new day-->
	//		<GC>380</GC>
	//		<!-- Global counter of receipt/invoice which increment throughout the life of the VFD. It has the same value as RCTNUM-->
	//		<ZNUM>20190827</ZNUM>
	//		<!-- ZNUM will be a date of receipt/invoice generated as number in format of (YYYYMMDD) -->
	//		<RCTVNUM>MFT7AB380</RCTVNUM>
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

	signedPayload, err := Sign(privateKey, out)
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

func GenerateReceipt(params ReceiptParams, customer Customer, items []Item, payments []Payment) *models.RCT {
	vatAmountMap := make(map[string]float64)
	// add A and C items
	vatAmountMap["A"] = 0
	vatAmountMap["C"] = 0
	rctItems := make([]*models.ITEM, len(items))
	totals := &models.TOTALS{
		TOTALTAXEXCL: 0,
		TOTALTAXINCL: 0,
		DISCOUNT:     0,
	}
	rctVatTotals := make([]*models.VATTOTAL, 0)

	totalTax := 0.0

	for i, item := range items {
		id := item.ID
		desc := item.Description
		qty := item.Quantity
		price := item.Price
		taxCode := item.TaxCode
		discount := item.Discount
		rctItems[i] = &models.ITEM{
			ID:      id,
			DESC:    desc,
			QTY:     qty,
			TAXCODE: taxCode,
			AMT:     price,
		}

		// add discount
		totals.DISCOUNT += discount

		// if item is taxable add to total taxCode
		if taxCode == 1 {
			// add another entry to the vatAmountMap with key "A"
			amountInA := vatAmountMap["A"]
			amountPaidTaxInclusive := price * qty
			itemTax := amountPaidTaxInclusive * 0.18
			totalTax += itemTax
		}

		// add totals tax inclusive
		totals.TOTALTAXINCL += price * qty

	}

	// add totals tax exclusive
	totals.TOTALTAXEXCL = totals.TOTALTAXINCL - totalTax

	// make payments
	rctPayments := make([]*models.PAYMENT, len(payments))
	for i, payment := range payments {
		rctPayments[i] = &models.PAYMENT{
			PMTTYPE:   string(payment.Type),
			PMTAMOUNT: payment.Amount,
		}
	}

	return &models.RCT{
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
		ITEMS: models.ITEMS{
			ITEM: rctItems,
		},
		TOTALS: *totals,
		PAYMENTS: models.PAYMENTS{
			PAYMENT: rctPayments,
		},
		VATTOTALS: models.VATTOTALS{
			XMLName:  xml.Name{},
			Text:     "",
			VATTOTAL: totals,
		},
	}
}

package vfd

import (
	"context"
	"crypto/rsa"
	"github.com/vfdcloud/vfd/models"
)

const (
	StandardVatrateId                 = "A"
	StandardVatrate                   = 18
	CashPaymentType       PaymentType = "CASH"
	CreditCardPaymentType PaymentType = "CCARD"
	ChequePaymentType     PaymentType = "CHEQUE"
	InvoicePaymentType    PaymentType = "INVOICE"
	ElectronicPaymentType PaymentType = "EMONEY"
	TinCustomerID         CustomerID  = 1
	LicenceCustomerID     CustomerID  = 2
	VoterIDCustomerID     CustomerID  = 3
	PassportCustomerID    CustomerID  = 4
	NidaCustomerID        CustomerID  = 5
	NonCustomerID         CustomerID  = 6
	MeterNumberCustomerID CustomerID  = 7
)

type (
	// PaymentType represent the type of payment that is recognized by the VFD server
	// There are five types of payments: CASH, CHEQUE, CCARD, EMONEY and INVOICE
	PaymentType string

	// CustomerID is the type of ID the customer used during purchase
	// The Type of ID is to be included in the receipt.
	// Allowed values for CustomerID are 1 through 7. The number to type
	// mapping are as follows:
	// 1: Tax Identification Number (TIN), 2: Driving License, 3: Voters Number,
	// 4: Travel Passport, 5: National ID, 6: NIL (No Identity Used), 7: Meter Number
	CustomerID int

	// RequestHeaders represent collection of request headers during receipt or Z report
	// sending via VFD Service
	RequestHeaders struct {
		ContentType string
		CertSerial  string
		BearerToken string
		RoutingKey  string
	}

	Payment struct {
		Type   PaymentType
		Amount float64
	}

	VatTotal struct {
		ID     string
		Rate   float64
		Amount float64
	}

	// Response contains details returned when submitting a receipt to the VFD Service
	// or a Z report.
	// Number (int) is the receipt number in case of a receipt submission and the
	// Z report number in case of a Z report submission.
	// Date (string) is the date of the receipt or Z report submission. The format
	// is YYYY-MM-DD.
	// Time (string) is the time of the receipt or Z report submission. The format
	// is HH24:MI:SS
	// Code (int) is the response code. 0 means success.
	// Message (string) is the response message.
	Response struct {
		Number  int64  `json:"number,omitempty"`
		Date    string `json:"date,omitempty"`
		Time    string `json:"time,omitempty"`
		Code    int64  `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}

	// LoadCertFunc is a function that loads a certificate from a pfx file from the given path.
	// returns the private key and the certificate.

	Service interface {
		// Register is used to register a virtual fiscal device (VFD) with the VFD Service.
		// If successful, the VFD Service returns a registration response containing the
		// VFD details and the credentials to use when submitting receipts and Z reports.
		// Registering a VFD is a one-time operation. The subsequent calls to Register will
		// yield the same response.VFD should store the registration response to
		// avoid calling Register again.
		Register(ctx context.Context, url string, request *RegistrationRequest) (*models.RegistrationResponse, error)

		// FetchToken is used to fetch a token from the VFD Service. The token is used
		// to authenticate the VFD when submitting receipts and Z reports.
		// credentials used here are the ones returned by the Register method.
		FetchToken(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error)

		// SubmitReceipt is used to submit a receipt to the VFD Service. The receipt
		// is signed using the private key. The private key is obtained from the certificate
		// issued by the Revenue Authority during integration.
		SubmitReceipt(
			ctx context.Context, url string, headers *RequestHeaders,
			privateKey *rsa.PrivateKey, receipt *models.RCT) (*Response, error)

		// SubmitReport is used to submit a Z report to the VFD Service. The Z report
		// is signed using the private key. The private key is obtained from the certificate
		// issued by the Revenue Authority during integration.
		SubmitReport(
			ctx context.Context, url string, headers *RequestHeaders,
			privateKey *rsa.PrivateKey, report *models.Report) (*Response, error)
	}
)

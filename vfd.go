package vfd

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"github.com/vfdcloud/base"
	"github.com/vfdcloud/vfd/internal"
	"github.com/vfdcloud/vfd/models"
)

const (
	VatrateStandard         = "A"
	VatrateSpecial          = "B"
	VatrateNonVatItems      = "C"
	VatrateSpecialRelief    = "D"
	VatrateExempt           = "E"
	TaxCodeTaxableItem      = 1
	TaxCodeTaxFreeItem      = 3
	CustomerIDTaxIDNumber   = 1
	CustomerIDDriverLicense = 2
	CustomerIDVoterID       = 3
	CustomerIDPassport      = 4
	CustomerIDNationalID    = 5
	CustomerIDNone          = 6
	CustomerIDMeterNumber   = 7
	PaymentTypeCash         = "CASH"    // Cash
	PaymentTypeCheque       = "CHEQUE"  // Cheque
	PaymentTypeCCard        = "CCARD"   // Credit Card
	PaymentTypeEMoney       = "EMONEY"  // Electronic Money
	PaymentTypeInvoice      = "INVOICE" // Invoice

)

const (
	RequestTypeToken    = RequestType("TOKEN REQUEST")
	RequestTypeRegister = RequestType("REGISTRATION")
	RequestTypeReceipt  = RequestType("RECEIPT UPLOAD")
	RequestTypeReport   = RequestType("REPORT UPLOAD")
)

type (

	// API is an interface for the VFD API httpx. The interface should not hide some sort of state
	// that the implementation may need to maintain. The Ideal implementation should be stateless.
	// Hence, the interface should not hide details of the implementation.
	API interface {
		LoadCert(ctx context.Context, certPath string, certPassword string) (*rsa.PrivateKey, *x509.Certificate, error)
		VerifySignature(ctx context.Context, privateKey *rsa.PublicKey, payload []byte, signature string) error
		SignPayload(ctx context.Context, privateKey *rsa.PrivateKey, payload []byte) ([]byte, error)
		Register(ctx context.Context, request *RegistrationRequest) (*models.RegistrationResponse, error)
		Token(ctx context.Context, request *TokenRequest) (*TokenResponse, error)
		UploadReceipt(ctx context.Context, request *ReceiptRequest, receipt *models.RCT) (*ReceiptResponse, error)
	}

	RequestType string
)

func (p *Paths) Get(e base.Env, req RequestType) string {
	switch req {
	case RequestTypeToken:
		if e == base.StagingEnv {
			return internal.AppendEndpoint(p.BaseURL, "efdmsRctApi", p.TokenEndpoint)
		}

		return internal.AppendEndpoint(p.BaseURL, p.TokenEndpoint)

	case RequestTypeRegister:
		return internal.AppendEndpoint(p.BaseURL, p.APIPath, p.RegistrationEndpoint)

	case RequestTypeReceipt:
		return internal.AppendEndpoint(p.BaseURL, p.APIPath, p.ReceiptEndpoint)

	case RequestTypeReport:
		return internal.AppendEndpoint(p.BaseURL, p.APIPath, p.ReportEndpoint)

	default:
		return ""
	}
}

func GetRequestURL(e base.Env, req RequestType) string {
	p := FetchPaths(e)

	return p.Get(e, req)
}

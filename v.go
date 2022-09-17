package vfd

import (
	"context"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/vfdcloud/base/crypto"
	"github.com/vfdcloud/vfd/models"
)

type (
	// RequestHeaders represent collection of request headers during receipt or Z report
	// sending via VFD Service
	RequestHeaders struct {
		ContentType string
		CertSerial  string
		BearerToken string
		RoutingKey  string
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
	LoadCertFunc func(
		ctx context.Context, certPath string, certPassword string) (
		*rsa.PrivateKey, *x509.Certificate, error,
	)

	VerifySignatureFunc func(
		ctx context.Context, publicKey *rsa.PublicKey,
		payload []byte, signature string) error

	SignPayloadFunc func(
		ctx context.Context, privateKey *rsa.PrivateKey,
		payload []byte) ([]byte, error)

	RegisterClientFunc func(
		ctx context.Context,
		url string,
		request *RegistrationRequest,
	) (*RegistrationResponse, error)

	ReportSubmitFunc func(
		ctx context.Context,
		url string,
		headers *RequestHeaders,
		privateKey *rsa.PrivateKey,
		request *models.Report,
	) (*Response, error)

	Service interface {
		Register(
			ctx context.Context,
			url string,
			request *RegistrationRequest,
		) (*RegistrationResponse, error)
		FetchToken(
			ctx context.Context,
			url string,
			request *TokenRequest,
		) (*TokenResponse, error)
		SubmitReceipt(
			ctx context.Context,
			url string,
			headers *RequestHeaders,
			privateKey *rsa.PrivateKey,
			receipt *models.RCT,
		) (*Response, error)

		SubmitReport(
			ctx context.Context,
			url string,
			headers *RequestHeaders,
			privateKey *rsa.PrivateKey,
			report *models.Report,
		) (*Response, error)
	}

	Registrar interface {
		Register(ctx context.Context, url string, request *RegistrationRequest) (*RegistrationResponse, error)
	}

	ReportSubmitter interface {
		SubmitReport(ctx context.Context, url string, headers *RequestHeaders,
			privateKey *rsa.PrivateKey,
			report *models.Report) (*Response, error)
	}
)

func LoadCert(ctx context.Context, certPath string, certPassword string) (
	*rsa.PrivateKey, *x509.Certificate, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	return crypto.ParsePfxCertificate(certPath, certPassword)
}

func VerifySignature(ctx context.Context, publicKey *rsa.PublicKey, payload []byte, signature string) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	sg, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("could not verify signature %w", err)
	}

	hash := sha1.Sum(payload) //nolint:gosec
	err = crypto.VerifySignature(publicKey, hash[:], sg)
	if err != nil {
		return fmt.Errorf("could not verify signature %w", err)
	}

	return nil
}

func SignPayload(ctx context.Context, privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	out, err := crypto.SignPayload(privateKey, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to sign the payload: %w", err)
	}

	err = VerifySignature(ctx, &privateKey.PublicKey, payload, base64.StdEncoding.EncodeToString(out))
	if err != nil {
		return nil, fmt.Errorf("invalid signature %w", err)
	}

	return out, nil
}

func (registrar RegisterClientFunc) Register(ctx context.Context, url string,
	request *RegistrationRequest) (*RegistrationResponse, error) {
	return registrar(ctx, url, request)
}

func (fetcher TokenFetcher) FetchToken(ctx context.Context, url string,
	request *TokenRequest) (*TokenResponse, error) {
	return fetcher(ctx, url, request)
}

func (submitter ReportSubmitFunc) SubmitReport(
	ctx context.Context, url string, headers *RequestHeaders,
	privateKey *rsa.PrivateKey, report *models.Report,
) (*Response, error) {
	return submitter(ctx, url, headers, privateKey, report)
}

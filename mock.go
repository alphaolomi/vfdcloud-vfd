package vfd

import (
	"context"
	"crypto/rsa"
	"github.com/vfdcloud/vfd/models"
)

var _ Service = (*mock)(nil)

type mock struct{}

func (m *mock) Register(ctx context.Context, url string, key *rsa.PrivateKey, request *RegistrationRequest) (
	*models.RegistrationResponse, error) {
	return Register(ctx, url, key, request, nil)
}

func (m *mock) FetchToken(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error) {
	return FetchToken(ctx, url, request, nil)
}

func (m *mock) SubmitReceipt(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	receipt *ReceiptRequest) (*Response, error) {
	return SubmitReceipt(ctx, url, headers, privateKey, receipt, nil)
}

func (m *mock) SubmitReport(ctx context.Context, url string, headers *RequestHeaders, privateKey *rsa.PrivateKey,
	report *ReportRequest) (*Response, error) {
	return SubmitReport(ctx, url, headers, privateKey, report, nil)
}

package test

import (
	"context"
	"crypto/rsa"
	"github.com/vfdcloud/vfd"
	"net/http"
)

var _ vfd.Service = (*Mock)(nil)

type (

	// Mock is a mock implementation of the Service interface
	Mock struct {
		Server *http.Server
	}
)

func (m *Mock) Register(ctx context.Context, url string, privateKey *rsa.PrivateKey, request *vfd.RegistrationRequest) (*vfd.RegistrationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Mock) FetchToken(ctx context.Context, url string, request *vfd.TokenRequest) (*vfd.TokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Mock) SubmitReceipt(ctx context.Context, url string, headers *vfd.RequestHeaders, privateKey *rsa.PrivateKey, receipt *vfd.ReceiptRequest) (*vfd.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Mock) SubmitReport(ctx context.Context, url string, headers *vfd.RequestHeaders, privateKey *rsa.PrivateKey, report *vfd.ReportRequest) (*vfd.Response, error) {
	//TODO implement me
	panic("implement me")
}


package vfd

import (
	"context"
	"crypto/rsa"
	"net/http"

	xhttp "github.com/vfdcloud/vfd/internal/http"
)

type (
	Client struct {
		http *http.Client
	}

	Option func(*Client)
)

func WithHttpClient(http *http.Client) Option {
	return func(c *Client) {
		c.http = http
	}
}

func NewClient(options ...Option) *Client {
	client := &Client{
		http: xhttp.Instance(),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) Register(ctx context.Context,
	url string, privateKey *rsa.PrivateKey,
	request *RegistrationRequest,
) (*RegistrationResponse, error) {
	response, err := register(ctx, c.http, url, privateKey, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) FetchToken(ctx context.Context, url string,
	request *TokenRequest,
) (*TokenResponse, error) {
	return fetchToken(ctx, c.http, url, request)
}

func (c *Client) SubmitReceipt(
	ctx context.Context,
	url string,
	headers *RequestHeaders,
	privateKey *rsa.PrivateKey,
	receipt *ReceiptRequest,
) (*Response, error) {
	return submitReceipt(ctx, c.http, url, headers, privateKey, receipt)
}

func (c *Client) SubmitReport(
	ctx context.Context,
	url string,
	headers *RequestHeaders,
	privateKey *rsa.PrivateKey,
	report *ReportRequest,
) (*Response, error) {
	return submitReport(ctx, c.http, url, headers, privateKey, report)
}

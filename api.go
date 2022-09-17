package vfd

import (
	"context"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/vfdcloud/base/crypto"
	"net/http"
)

var (
	_                    API = (*client)(nil)
	ErrSignatureMismatch     = errors.New("signature mismatch")

	ErrReportUploadFailed = errors.New("report upload failed")
)

type (
	RequestInterceptor func(request *http.Request)

	ResponseInterceptor func(response *http.Response)
)

func (c *client) LoadCert(ctx context.Context, certPath string, certPassword string) (*rsa.PrivateKey, *x509.Certificate, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return crypto.ParsePfxCertificate(certPath, certPassword)

}

func (c *client) VerifySignature(ctx context.Context, publicKey *rsa.PublicKey, payload []byte, signature string) error {
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

func (c *client) SignPayload(ctx context.Context, privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	out, err := crypto.SignPayload(privateKey, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to sign the payload: %w", err)
	}

	err = c.VerifySignature(ctx, &privateKey.PublicKey, payload, base64.StdEncoding.EncodeToString(out))
	if err != nil {
		return nil, fmt.Errorf("invalid signature %w", err)
	}

	return out, nil
}

package vfd

import (
	"context"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/vfdcloud/base/crypto"
)

type (
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
)

func LoadCert(ctx context.Context, certPath string, certPassword string) (
	*rsa.PrivateKey, *x509.Certificate, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	return crypto.ParsePfxCertificate(certPath, certPassword)
}

func Sign(ctx context.Context, privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	signature, err := crypto.SignPayload(privateKey, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to sign the payload: %w", err)
	}

	hash := sha1.Sum(payload) //nolint:gosec
	err = crypto.VerifySignature(&privateKey.PublicKey, hash[:], signature)
	if err != nil {
		return nil, fmt.Errorf("could not verify signature %w", err)
	}

	return signature, nil
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

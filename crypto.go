package vfd

import (
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/vfdcloud/base/crypto"
)

type (
	// CertLoader loads a certificate from a file and returns the private key and the certificate
	CertLoader func(certPath string, certPassword string) (*rsa.PrivateKey, *x509.Certificate, error)

	// SignatureVerifier verifies the signature of a payload using the public key
	// of the signing certificate
	SignatureVerifier func(publicKey *rsa.PublicKey, payload []byte, signature string) error

	// PayloadSigner signs a payload using the private key of the signing certificate
	// all requests to the VFD API must be signed.
	PayloadSigner func(privateKey *rsa.PrivateKey, payload []byte) ([]byte, error)
)

func LoadCert(certPath string, certPassword string) (
	*rsa.PrivateKey, *x509.Certificate, error,
) {
	return crypto.ParsePfxCertificate(certPath, certPassword)
}

func Sign(privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
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

func VerifySignature(publicKey *rsa.PublicKey, payload []byte, signature string) error {
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

func SignPayload(privateKey *rsa.PrivateKey, payload []byte) ([]byte, error) {
	out, err := crypto.SignPayload(privateKey, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to sign the payload: %w", err)
	}

	err = VerifySignature(&privateKey.PublicKey, payload, base64.StdEncoding.EncodeToString(out))
	if err != nil {
		return nil, fmt.Errorf("invalid signature %w", err)
	}

	return out, nil
}

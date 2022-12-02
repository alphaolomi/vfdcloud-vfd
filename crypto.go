package vfd

import (
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/vfdcloud/base/crypto"
	"software.sslmate.com/src/go-pkcs12"
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

func LoadCertChain(certPath string, certPassword string) (*rsa.PrivateKey, *x509.Certificate, []*x509.Certificate, error) {
	pfxData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not read the certificate file: %w", err)
	}
	pfx, cert, caCerts, err := pkcs12.DecodeChain(pfxData, certPassword)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not decode the certificate file: %w", err)
	}

	// type check to make sure we have a private key
	privateKey, ok := pfx.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, nil, fmt.Errorf("private key is not of type *rsa.PrivateKey: %w", err)
	}

	return privateKey, cert, caCerts, nil
}

func LoadCert(path, password string) (*rsa.PrivateKey, *x509.Certificate, error) {
	pfxData, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	pfx, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		if err.Error() == "pkcs12: expected exactly two safe bags in the PFX PDU" {
			privateKey, cert, _, err := LoadCertChain(path, password)
			if err != nil {
				return nil, nil, err
			}
			return privateKey, cert, nil
		}
		return nil, nil, err
	}

	// type check to make sure we have a private key
	privateKey, ok := pfx.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("private key is not of type rsa.PrivateKey")
	}

	return privateKey, cert, nil
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

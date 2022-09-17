package vfd

import (
	"context"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec
	"fmt"
	"github.com/vfdcloud/base/crypto"
)

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

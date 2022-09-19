package vfd

import (
	"encoding/base64"
	"fmt"
	"github.com/vfdcloud/base"
	"strings"
)

// ReceiptLink creates a link to the receipt
func ReceiptLink(e base.Env, receiptVerificationNumber, receiptVerificationTime string) string {
	var baseURL string

	if e == base.ProdEnv {
		baseURL = VerifyURLProd
	} else {
		baseURL = VerifyURLTest
	}
	return receiptLink(baseURL, receiptVerificationNumber, receiptVerificationTime)
}

func receiptLink(baseURL string, receiptVerificationNumber, receiptVerificationTime string) string {
	return fmt.Sprintf(
		"%s%s_%s",
		baseURL,
		receiptVerificationNumber,
		strings.ReplaceAll(receiptVerificationTime, ":", ""))
}

// EncodeBase64Bytes calls base64.StdEncoding.EncodeToString
func EncodeBase64Bytes(val []byte) string {
	return base64.StdEncoding.EncodeToString(val)
}

// EncodeBase64String calls base64.StdEncoding.EncodeToString
func EncodeBase64String(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

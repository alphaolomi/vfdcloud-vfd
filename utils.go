package vfd

import (
	"fmt"
	"github.com/vfdcloud/base"
	"strings"
)

const (
	prodReceiptBaseURL = "https://verify.tra.go.tz/"
	devReceiptBaseURL  = "https://virtual.tra.go.tz/efdmsRctVerify/"
)

// ReceiptLink creates a link to the receipt
func ReceiptLink(e base.Env, receiptVerificationNumber, receiptVerificationTime string) string {
	var baseURL string

	if e == base.ProdEnv {
		baseURL = prodReceiptBaseURL
	} else {
		baseURL = devReceiptBaseURL
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

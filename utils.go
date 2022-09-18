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

func AppendEndpoint(url string, endpoints ...string) string {
	if len(endpoints) == 1 {
		return appendEndpoint(url, endpoints[0])
	}

	finalPath := url

	for _, endpoint := range endpoints {
		finalPath = appendEndpoint(finalPath, endpoint)
	}

	return finalPath
}

// appendEndpoint appends a path to a requestURL.
func appendEndpoint(url string, endpoint string) string {
	var (
		trimRight = strings.TrimRight
		replace   = strings.ReplaceAll
		trimLeft  = strings.TrimLeft
	)

	// remove all leading and trailing whitespaces
	url, endpoint = replace(url, " ", ""), replace(endpoint, " ", "")

	// for baseurl trim all trailing slashes and leading slashes
	// for endpoint trim all leading slashes
	url = trimRight(url, "/")
	url = trimLeft(url, "/")
	endpoint = trimLeft(endpoint, "/")

	if url == "" && endpoint == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", url, endpoint)
}

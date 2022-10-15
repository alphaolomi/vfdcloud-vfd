package vfd

import "encoding/base64"

// encodeBase64Bytes calls base64.StdEncoding.EncodeToString
func encodeBase64Bytes(val []byte) string {
	return base64.StdEncoding.EncodeToString(val)
}

// encodeBase64String calls base64.StdEncoding.EncodeToString
func encodeBase64String(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

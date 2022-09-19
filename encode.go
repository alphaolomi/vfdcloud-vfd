package vfd

import "encoding/base64"

// EncodeBase64Bytes calls base64.StdEncoding.EncodeToString
func EncodeBase64Bytes(val []byte) string {
	return base64.StdEncoding.EncodeToString(val)
}

// EncodeBase64String calls base64.StdEncoding.EncodeToString
func EncodeBase64String(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

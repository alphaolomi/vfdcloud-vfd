package internal

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

const (
	JSON EnvelopeContentType = "application/json"
	XML  EnvelopeContentType = "application/xml"
)

type (
	EnvelopeContentType string
	Payload             interface {
		Sign(*rsa.PrivateKey, EnvelopeContentType) (string, error)
	}

	// Envelope is a request body to be sent to VFD server in making a request.
	// These requests are currently Registration, Receipt Upload and Report Upload.
	// Login request does not follow the same pattern.
	// The Envelope has Body and Signature. Body contains information depends on a
	// a request and Signature is obtaining by signing the body with private key.
	// before signing the body, the body is encoded in JSON or XML.
	Envelope interface {
		setSignature(signature string) error
		setPayload(payload Payload) error
		Bytes(EnvelopeContentType) ([]byte, error)
	}
	Request struct {
		Method   string
		BaseURL  string
		Endpoint string
		Headers  map[string]string
		FormMap  map[string]string
		Payload  Payload
		Envelope Envelope
	}

	RequestInterceptor func(request *http.Request) error
)

//func NewRequestWithContext(ctx context.Context, privateKey *rsa.PrivateKey, request *Request, interceptors ...RequestInterceptor) (*http.Request, error) {
//	var (
//		method   = request.Method
//		baseURL  = request.BaseURL
//		headers  = request.Headers
//		endpoint = request.Endpoint
//		payload  = request.Payload
//		envelope = request.Envelope
//		err      error
//		req      = new(http.Request)
//		buffer   = new(bytes.Buffer)
//	)
//
//	if request.FormMap != nil {
//		formMap := request.FormMap
//		form := url.Values{}
//		for key, value := range formMap {
//			form.Set(key, value)
//		}
//
//		buffer = bytes.NewBufferString(form.Encode())
//	} else {
//		signature, err := payload.Sign(privateKey)
//		if err != nil {
//			return nil, err
//		}
//
//		err = envelope.setPayload(payload)
//		if err != nil {
//			return nil, err
//		}
//
//		err = envelope.setSignature(signature)
//		if err != nil {
//			return nil, err
//		}
//
//		out, err := envelope.Bytes()
//
//		if err != nil {
//			return nil, err
//		}
//
//		buffer = bytes.NewBuffer(out)
//	}
//
//	req, err = http.NewRequestWithContext(ctx, method, AppendEndpoint(baseURL, endpoint), buffer)
//
//	if err != nil {
//		return nil, err
//	}
//
//	for key, value := range headers {
//		req.Header.Set(key, value)
//	}
//
//	return req, nil
//
//}

func Base64(val string) string {
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

// appendEndpoint appends a path to a URL.
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

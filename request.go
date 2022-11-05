package vfd

import (
	"io"
	"net/http"
)

// CreateRequest returns a new *http.Request with the given method, URL,
// and optional body.
func CreateRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

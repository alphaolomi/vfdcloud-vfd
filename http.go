package vfd

import (
	"net/http"
	"sync"
	"time"
)

const (
	JsonContentType = "application/json"
	XMLContentType  = "application/xml"
	ClientType      = "webapi"
)

type (
	// httpx is a wrapper for the http.Client that is used internally to make
	// http requests to the VFD server.
	httpx struct {
		client *http.Client
	}
)

var (
	once     sync.Once
	instance *httpx
)

func getHttpClientInstance() *httpx {
	once.Do(func() {
		instance = defaultHTTPClient()
	})

	return instance
}

func defaultHTTPClient() *httpx {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Timeout:   70 * time.Second,
		Transport: t,
	}
	c := &httpx{
		client: httpClient,
	}

	return c
}

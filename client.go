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
	client struct {
		http *http.Client
	}
)

var (
	once     sync.Once
	instance *client
)

func getInstance() *client {

	once.Do(func() {
		instance = makeClient()
	})

	return instance
}

func makeClient() *client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: t,
	}
	c := &client{
		http: httpClient,
	}

	return c
}

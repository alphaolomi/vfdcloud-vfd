package http

import (
	"net/http"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *http.Client
)

func Instance() *http.Client {
	once.Do(func() { instance = defaultInstance() })
	return instance
}

func defaultInstance() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Timeout:   70 * time.Second,
		Transport: t,
	}
	return httpClient
}

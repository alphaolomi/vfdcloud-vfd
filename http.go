package vfd

import (
	"github.com/vfdcloud/base"
	"net/http"
	"sync"
	"time"
)

var (
	productionURLs = &requestURL{
		Registration:  RegisterProductionURL,
		FetchToken:    FetchTokenProductionURL,
		SubmitReceipt: SubmitReceiptProductionURL,
		SubmitReport:  SubmitReportProductionURL,
		VerifyReceipt: VerifyReceiptProductionURL,
	}

	stagingURLs = &requestURL{
		Registration:  RegisterTestingURL,
		FetchToken:    FetchTokenTestingURL,
		SubmitReceipt: SubmitReceiptTestingURL,
		SubmitReport:  SubmitReportTestingURL,
		VerifyReceipt: VerifyReceiptTestingURL,
	}
)

type (
	// httpx is a wrapper for the http.Client that is used internally to make
	// http requests to the VFD server.
	httpx struct{ client *http.Client }
)

var (
	once     sync.Once
	instance *httpx
)

func getHttpClientInstance() *httpx {
	once.Do(func() { instance = defaultHTTPClient() })
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

func RequestURL(e base.Env, action Action) string {
	var u *requestURL
	if e == base.ProdEnv {
		u = productionURLs
	} else {
		u = stagingURLs
	}

	switch action {
	case RegisterClientAction:
		return u.Registration
	case FetchTokenAction:
		return u.FetchToken
	case SubmitReceiptAction:
		return u.SubmitReceipt
	case SubmitReportAction:
		return u.SubmitReport
	case ReceiptVerificationAction:
		return u.VerifyReceipt
	default:
		return ""
	}
}

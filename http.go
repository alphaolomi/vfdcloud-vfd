package vfd

import (
	"github.com/vfdcloud/base"
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

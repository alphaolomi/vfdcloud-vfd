package vfd

import "github.com/vfdcloud/base"

const (
	RegisterURLProd = "https://vfd.tra.go.tz/api/vfdRegReq"
	TokenURLProd    = "https://vfd.tra.go.tz/vfdtoken"
	ReceiptURLProd  = "https://vfd.tra.go.tz/api/efdmsRctInfo"
	ReportURLProd   = "https://vfd.tra.go.tz/api/efdmszreport"
	VerifyURLProd   = "https://verify.tra.go.tz/"
	RegisterURLTest = "https://virtual.tra.go.tz/efdmsRctApi/api/vfdRegReq"
	TokenURLTest    = "https://virtual.tra.go.tz/efdmsRctApi/vfdtoken"
	ReceiptURLTest  = "https://virtual.tra.go.tz/efdmsRctApi/api/efdmsRctInfo"
	ReportURLTest   = "https://virtual.tra.go.tz/efdmsRctApi/api/efdmszreport"
	VerifyURLTest   = "https://virtual.tra.go.tz/efdmsRctVerify/"
)

const (
	RegisterClientAction      Action = "register"
	FetchTokenAction          Action = "token"
	UploadReceiptAction       Action = "receipt"
	UploadReportAction        Action = "report"
	ReceiptVerificationAction Action = "verification"
)

type (
	Action string

	requestURL struct {
		Registration  string
		FetchToken    string
		SubmitReceipt string
		SubmitReport  string
		VerifyReceipt string
	}
)

var (
	productionURLs = &requestURL{
		Registration:  RegisterURLProd,
		FetchToken:    TokenURLProd,
		SubmitReceipt: ReceiptURLProd,
		SubmitReport:  ReportURLProd,
		VerifyReceipt: VerifyURLProd,
	}

	stagingURLs = &requestURL{
		Registration:  RegisterURLTest,
		FetchToken:    TokenURLTest,
		SubmitReceipt: ReceiptURLTest,
		SubmitReport:  ReportURLTest,
		VerifyReceipt: VerifyURLTest,
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
	case UploadReceiptAction:
		return u.SubmitReceipt
	case UploadReportAction:
		return u.SubmitReport
	case ReceiptVerificationAction:
		return u.VerifyReceipt
	default:
		return ""
	}
}

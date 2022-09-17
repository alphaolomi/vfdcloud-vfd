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
	ProductionBaseURL         = "https://vfd.tra.go.tz/"
	ProductionApiPath         = "/api/"
	StagingApiPath            = "/efdmsRctApi/api/"
	StagingBaseURL            = "https://virtual.tra.go.tz/"
	TokenEndpoint             = "/vfdtoken"
	RegistrationEndpoint      = "/vfdRegReq"
	ReceiptEndpoint           = "/efdmsRctInfo"
	ReportEndpoint            = "/efdmszreport"
	StagingVerificationURL    = "https://virtual.tra.go.tz/efdmsRctVerify/"
	ProductionVerificationURL = "https://verify.tra.go.tz/"
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

	URL struct {
		RegisterURL string
		TokenURL    string
		ReceiptURL  string
		ReportURL   string
		VerifyURL   string
	}

	Paths struct {
		BaseURL              string
		APIPath              string
		RegistrationEndpoint string
		TokenEndpoint        string
		ReceiptEndpoint      string
		ReportEndpoint       string
		VerificationURL      string
	}
)

var (
	productionURLs = &URL{
		RegisterURL: RegisterURLProd,
		TokenURL:    TokenURLProd,
		ReceiptURL:  ReceiptURLProd,
		ReportURL:   ReportURLProd,
		VerifyURL:   VerifyURLProd,
	}

	stagingURLs = &URL{
		RegisterURL: RegisterURLTest,
		TokenURL:    TokenURLTest,
		ReceiptURL:  ReceiptURLTest,
		ReportURL:   ReportURLTest,
		VerifyURL:   VerifyURLTest,
	}
)

func GetURL(e base.Env, action Action) string {
	var u *URL
	if e == base.StagingEnv {
		u = stagingURLs
	} else {
		u = productionURLs
	}

	switch action {
	case RegisterClientAction:
		return u.RegisterURL
	case FetchTokenAction:
		return u.TokenURL
	case UploadReceiptAction:
		return u.ReceiptURL
	case UploadReportAction:
		return u.ReportURL
	case ReceiptVerificationAction:
		return u.VerifyURL
	default:
		return ""
	}
}

func FetchPaths(e base.Env) *Paths {
	var p *Paths
	if e == base.StagingEnv {
		p = &Paths{
			BaseURL:              StagingBaseURL,
			APIPath:              StagingApiPath,
			RegistrationEndpoint: RegistrationEndpoint,
			TokenEndpoint:        TokenEndpoint,
			ReceiptEndpoint:      ReceiptEndpoint,
			ReportEndpoint:       ReportEndpoint,
			VerificationURL:      StagingVerificationURL,
		}
	} else {
		p = &Paths{
			BaseURL:              ProductionBaseURL,
			APIPath:              ProductionApiPath,
			RegistrationEndpoint: RegistrationEndpoint,
			TokenEndpoint:        TokenEndpoint,
			ReceiptEndpoint:      ReceiptEndpoint,
			ReportEndpoint:       ReportEndpoint,
			VerificationURL:      ProductionVerificationURL,
		}
	}

	return p
}

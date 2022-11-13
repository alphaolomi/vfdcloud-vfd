package test

import (
	"context"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vfdcloud/base"
	"github.com/vfdcloud/vfd"
)

//	RegisterProductionURL                  = "https://vfd.tra.go.tz/api/vfdRegReq"
//	FetchTokenProductionURL                = "https://vfd.tra.go.tz/vfdtoken" //nolint:gosec
//	SubmitReceiptProductionURL             = "https://vfd.tra.go.tz/api/efdmsRctInfo"
//	SubmitReportProductionURL              = "https://vfd.tra.go.tz/api/efdmszreport"
//	VerifyReceiptProductionURL             = "https://verify.tra.go.tz/"
//	RegisterTestingURL                     = "https://virtual.tra.go.tz/efdmsRctApi/api/vfdRegReq"
//	FetchTokenTestingURL                   = "https://virtual.tra.go.tz/efdmsRctApi/vfdtoken" //nolint:gosec
//	SubmitReceiptTestingURL                = "https://virtual.tra.go.tz/efdmsRctApi/api/efdmsRctInfo"
//	SubmitReportTestingURL                 = "https://virtual.tra.go.tz/efdmsRctApi/api/efdmszreport"
//	VerifyReceiptTestingURL                = "https://virtual.tra.go.tz/efdmsRctVerify/"

const (
	RegisterProductionEndpoint      = "/api/vfdRegReq"
	FetchTokenProductionEndpoint    = "/vfdtoken"
	SubmitReceiptProductionEndpoint = "/api/efdmsRctInfo"
	SubmitReportProductionEndpoint  = "/api/efdmszreport"
	VerifyReceiptProductionEndpoint = "/"
	RegisterTestingEndpoint         = "/efdmsRctApi/api/vfdRegReq"
	FetchTokenTestingEndpoint       = "/efdmsRctApi/vfdtoken"
	SubmitReceiptTestingEndpoint    = "/efdmsRctApi/api/efdmsRctInfo"
	SubmitReportTestingEndpoint     = "/efdmsRctApi/api/efdmszreport"
	VerifyReceiptTestingEndpoint    = "/"
)

var _ vfd.Service = (*TestServer)(nil)

// NewTestServer returns a new test server that can be used to test the client
// against. The server is configured to respond to the given requests.
type (
	TestServer struct {
		http              *httptest.Server
		Env               base.Env
		RegistrationReq   *vfd.RegistrationRequest
		RegistrationResp  *vfd.RegistrationResponse
		GetRequestURLFunc func(env base.Env, action vfd.Action) string
	}
)

// NewTestServer returns a new test server that can be used to test the client
// against. The server is configured to respond to
func NewTestServer(t *testing.T, env base.Env, action vfd.Action) (*TestServer, error) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch action {
		case vfd.RegisterClientAction:
			if r.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", r.Method)
			}
			switch env {
			case base.ProdEnv:
				if r.URL.Path != RegisterProductionEndpoint {
					t.Errorf("expected %s request, got %s", RegisterProductionEndpoint, r.URL.Path)
				}

			case base.StagingEnv:
				if r.URL.Path != RegisterTestingEndpoint {
					t.Errorf("expected %s request, got %s", RegisterTestingEndpoint, r.URL.Path)
				}
			}

			if r.URL.Path != RegisterProductionEndpoint {
				t.Errorf("expected %s request, got %s", RegisterProductionEndpoint, r.URL.Path)
			}
		}
	}))

	s := &TestServer{
		http: httpServer,
	}

	return s, nil
}

func (t *TestServer) Register(ctx context.Context, url string, privateKey *rsa.PrivateKey, request *vfd.RegistrationRequest) (*vfd.RegistrationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (t *TestServer) FetchToken(ctx context.Context, url string, request *vfd.TokenRequest) (*vfd.TokenResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (t *TestServer) SubmitReceipt(ctx context.Context, url string, headers *vfd.RequestHeaders, privateKey *rsa.PrivateKey, receipt *vfd.ReceiptRequest) (*vfd.Response, error) {
	// TODO implement me
	panic("implement me")
}

func (t *TestServer) SubmitReport(
	ctx context.Context, url string, headers *vfd.RequestHeaders,
	privateKey *rsa.PrivateKey, report *vfd.ReportRequest,
) (*vfd.Response, error) {
	// TODO implement me
	panic("implement me")
}

func TestSubmitRawRequest(t *testing.T) {
	type testCase struct {
		name  string
		input struct {
			Env      base.Env
			Action   vfd.Action
			FilePath string
			Headers  *vfd.RequestHeaders
		}
		want struct {
			response vfd.Response
			wantErr  bool
		}
	}

	tests := []testCase{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			got, err := vfd.SubmitRawRequest(ctx, tt.input.Headers, &vfd.RawRequest{
				Env:      tt.input.Env,
				Action:   tt.input.Action,
				FilePath: tt.input.FilePath,
			})

			if (err != nil) != tt.want.wantErr {
				t.Errorf("SubmitRawRequest() error = %v, wantErr %v", err, tt.want.wantErr)
				return
			}

			if got != nil {
				same := reflect.DeepEqual(got, tt.want.response)
				if !same {
					t.Errorf("SubmitRawRequest() = %v, want %v", got, tt.want.response)
				}
			}
		})
	}
}

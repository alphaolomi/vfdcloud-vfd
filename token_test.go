package vfd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type (
	FetchTokenTestParams struct {
		Username   string
		Password   string
		GrantType  string
		Panic      bool
		ContextErr bool
		Timeout    time.Duration
	}
)

func CreateTestFetchTokenFunc(t *testing.T, params *FetchTokenTestParams) FetchTokenFunc {
	t.Helper()

	return func(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error) {
		_, cancel := context.WithTimeout(context.TODO(), params.Timeout)

		t.Logf("\nrequest: %+v, params %+v\n", request, params)
		defer cancel()
		// check if we should panic
		if params.Panic {
			return nil, fmt.Errorf("panic")
		}
		if params.ContextErr {
			return nil, context.Canceled
		}

		if request.Username != params.Username {
			return &TokenResponse{
				Code:        "240",
				Message:     "invalid username",
				AccessToken: "",
				TokenType:   "",
				ExpiresIn:   0,
				Error:       "invalid username",
			}, nil
		}

		if request.Password != params.Password {
			return &TokenResponse{
				Code:        "241",
				Message:     "invalid password",
				AccessToken: "",
				TokenType:   "",
				ExpiresIn:   0,
				Error:       "invalid password",
			}, nil
		}

		if request.GrantType != params.GrantType {
			return &TokenResponse{
				Code:        "242",
				Message:     "invalid grant type",
				AccessToken: "",
				TokenType:   "",
				ExpiresIn:   0,
				Error:       "invalid grant type",
			}, nil
		}

		return &TokenResponse{
			Code:        "0",
			Message:     "success",
			AccessToken: "access_token_for_test",
			TokenType:   "bearer",
			ExpiresIn:   3600,
			Error:       "",
		}, nil
	}
}

func TestFetchToken(t *testing.T) {
	type args struct {
		ctx     context.Context
		url     string
		request *TokenRequest
		params  *FetchTokenTestParams
	}

	type test struct {
		name           string
		args           args
		want           *TokenResponse
		wantErr        bool
		isNetworkError bool
	}

	tests := []test{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				url: "/register",
				request: &TokenRequest{
					Username:  "admin",
					Password:  "admin",
					GrantType: "credentials",
				},
				params: &FetchTokenTestParams{
					Username:   "admin",
					Password:   "admin",
					GrantType:  "credentials",
					Panic:      false,
					ContextErr: false,
					Timeout:    1 * time.Second,
				},
			},
			want: &TokenResponse{
				Code:        "0",
				Message:     "success",
				AccessToken: "access_token_for_test",
				TokenType:   "bearer",
				ExpiresIn:   3600,
				Error:       "",
			},
			wantErr:        false,
			isNetworkError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// launch a test server in a goroutine
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// get request body then call fn
				// compare response with want
				request := &TokenRequest{}
				err := json.NewDecoder(r.Body).Decode(request)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
				}
				fn := CreateTestFetchTokenFunc(t, tt.args.params)
				response, err := fn(tt.args.ctx, tt.args.url, request)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
				}
				err = json.NewEncoder(w).Encode(response)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))
			defer srv.Close()
			reqURL, err := url.JoinPath(srv.URL, tt.args.url)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("\nreqURL: %s, request %+v\n", reqURL, tt.args.request)
			got, err := FetchToken(tt.args.ctx, reqURL, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.isNetworkError {
				if got == nil {
					t.Errorf("FetchToken() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

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
		defer cancel()
		// check if we should panic
		if params.Panic {
			return nil, fmt.Errorf("panic")
		}

		// check if we should return context error
		if params.ContextErr {
			return nil, context.Canceled
		}

		return &TokenResponse{
			Code:        "",
			Message:     "",
			AccessToken: "",
			TokenType:   "",
			ExpiresIn:   0,
			Error:       "",
		}, nil
	}
}

func TestFetchToken(t *testing.T) {
	t.Parallel()
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
					GrantType: "client_credentials",
				},
				params: &FetchTokenTestParams{
					Username:   "",
					Password:   "",
					GrantType:  "",
					Panic:      true,
					ContextErr: false,
					Timeout:    0,
				},
			},
			want: &TokenResponse{
				Code:        "",
				Message:     "",
				AccessToken: "",
				TokenType:   "",
				ExpiresIn:   0,
				Error:       "",
			},
			wantErr:        false,
			isNetworkError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

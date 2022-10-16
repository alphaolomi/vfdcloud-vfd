package vfd

import (
	"bytes"
	"context"
	"github.com/vfdcloud/base"
	"io"
	"net/http"
	"reflect"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// equals checks if two string values are equal.
func equals(t *testing.T, a, b string) {
	t.Helper()
	if a != b {
		t.Errorf("expected %q to equal %q", a, b)
	}
}

func TestFetchToken(t *testing.T) {
	type args struct {
		Env     base.Env
		Request *TokenRequest
		URL     string
		Ctx     context.Context
	}
	type wanted struct {
		Username  string
		Password  string
		GrantType string
	}
	type expected struct {
		StatusCode   int
		Error        bool
		Response     *TokenResponse
		ErrorMessage string
	}

	type test struct {
		name     string
		args     args
		expected expected
		wanted   wanted
	}

	tests := []test{
		{
			name: "TestFetchTokenStagingSuccess",
			args: args{
				Env: base.StagingEnv,
				Request: &TokenRequest{
					Username:  "username",
					Password:  "password",
					GrantType: "password",
				},
				URL: FetchTokenTestingURL,
				Ctx: context.Background(),
			},
			expected: expected{
				StatusCode: 200,
				Error:      false,
				Response: &TokenResponse{
					Code:        "",
					Message:     "",
					AccessToken: "access_token",
					TokenType:   "",
					ExpiresIn:   0,
					Error:       "",
				},
			},
			wanted: wanted{
				Username:  "username1",
				Password:  "password",
				GrantType: "password",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := NewTestClient(func(req *http.Request) *http.Response {
				tokenURL := RequestURL(base.StagingEnv, FetchTokenAction)
				equals(t, req.URL.String(), tokenURL)
				// extract form values
				err := req.ParseForm()
				if err != nil {
					t.Errorf("error parsing form: %v", err)
				}
				password := req.FormValue("password")
				equals(t, password, tt.wanted.Password)
				username := req.FormValue("username")
				equals(t, username, tt.wanted.Username)
				grantType := req.FormValue("grant_type")
				equals(t, grantType, tt.wanted.GrantType)

				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: io.NopCloser(bytes.NewBufferString(`OK`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			got, err := fetchToken(tt.args.Ctx, client, tt.args.URL, tt.args.Request)
			if (err != nil) != tt.expected.Error {
				t.Errorf("fetchToken() error = %v, wantErr %v", err, tt.expected.Error)
				return
			}
			if !reflect.DeepEqual(got, tt.expected.Response) {
				t.Errorf("fetchToken() got = %v, want %v", got, tt.expected.Response)
			}
		})
	}
}

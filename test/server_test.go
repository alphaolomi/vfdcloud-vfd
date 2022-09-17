// Package test is for vfd package integration testing
package test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

type (
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}
)

// requestOK is a test func that randomly generates a number
// between 1 and 100 and returns true if the number is less than
// or equal to 50. This is used to mock request failures.
func requestOK(t *testing.T) bool {
	t.Helper()
	n := rand.Intn(100)
	if n <= 50 {
		return true
	}
	return false
}

func HandlerFuncTest(t *testing.T) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		switch path {
		case "/token":
			if method != http.MethodPost {
				// return not allowed error
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}

			// mock request failure
			if !requestOK(t) {
				http.Error(w, "request failed", http.StatusInternalServerError)
				return
			}

			// return success
			w.WriteHeader(http.StatusOK)
			return

		case "/register":
			if method != http.MethodPost {
				// return not allowed error
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}

			// mock request failure
			if !requestOK(t) {
				http.Error(w, "request failed", http.StatusInternalServerError)
				return
			}

			// return success
			w.WriteHeader(http.StatusOK)
			return
		}

	}
}

func MakeTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(HandlerFuncTest(t))
}

func Test_Token_Request(t *testing.T) {
	t.Helper()
	t.Parallel()
	server := MakeTestServer(t)
	defer server.Close()
	serverURL := server.URL
	t.Logf("serverURL: %s", serverURL)
}

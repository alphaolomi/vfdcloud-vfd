package vfd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ErrFetchToken is the error returned when the token request fails.
// It is a wrapper for the underlying error.
var ErrFetchToken = errors.New("fetch token failed")

type (
	// TokenRequest contains the request parameters needed to get a token.
	// GrantType - The type of the grant_type.
	// Username - The username of the user.
	// Password - The password of the user.
	TokenRequest struct {
		Username  string
		Password  string
		GrantType string
	}

	// TokenResponse contains the response parameters returned by the token endpoint.
	TokenResponse struct {
		Code        string `json:"code,omitempty"`
		Message     string `json:"message,omitempty"`
		AccessToken string `json:"access_token,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
		ExpiresIn   int64  `json:"expires_in,omitempty"`
		Error       string `json:"error,omitempty"`
	}

	// TokenFetcher is an interface that fetches a token from the VFD Service
	// using the given url and the token request.
	// If the response status code is not 200, an error is returned.
	// The error message will contain the TokenResponse.Code and TokenResponse.Message
	// fields.
	TokenFetcher interface {
		FetchToken(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error)
	}

	TokenFetcherMiddleware func(fetcher TokenFetchFunc) TokenFetchFunc
)

func FetchTokenMiddleware(mw ...TokenFetcherMiddleware) TokenFetchFunc {
	fetcher := TokenFetchFunc(FetchToken)
	// Loop backwards through the middleware invoking each one. Replace the
	// fetcher with the new wrapped fetcher. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			fetcher = h(fetcher)
		}
	}

	return fetcher

}

// FetchToken retrieves a token from the VFD server. If the status code is not 200, an error is returned.
// Error Message will contain TokenResponse.Code and TokenResponse.Message
// fetchToken uses an internal client httpx with a timeout of 70 seconds.
// It is advised to call this only when the previous token has expired. It will still work if called before
// the token expires.
func FetchToken(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error) {
	httpClient := httpClientInstance().client

	return fetchToken(ctx, httpClient, url, request)
}

func (c *httpx) Token(ctx context.Context, requestUrl string, request *TokenRequest) (*TokenResponse, error) {
	httpClient := c.client

	return fetchToken(ctx, httpClient, requestUrl, request)
}

// fetchToken retrieves a token from the VFD server. If the status code is not 200, an error is returned.
// It is a context-aware function with a timeout of 1 minute
func fetchToken(ctx context.Context, client *http.Client, path string, request *TokenRequest) (*TokenResponse, error) {
	var (
		username  = request.Username
		password  = request.Password
		grantType = request.GrantType
	)

	// this request should have a max of 1 Minute timeout
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	var form = url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	form.Set("grant_type", grantType)
	buffer := bytes.NewBufferString(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, buffer)

	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrFetchToken, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client call error: %w: %v", ErrFetchToken, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "token error: could not close response body %v", err)
		}
	}(resp.Body)

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrFetchToken, err)
	}

	response := new(TokenResponse)

	if err := json.NewDecoder(bytes.NewBuffer(out)).Decode(response); err != nil {
		return nil, fmt.Errorf("response decode error: %w", err)
	}

	response.Code = resp.Header.Get("ACKCODE")
	response.Message = resp.Header.Get("ACKMSG")

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: error code=[%s],message=[%s], error=[%s]",
			ErrFetchToken, response.Code, response.Message, response.Error)
	}

	return response, nil
}

func (tr *TokenResponse) String() string {
	return fmt.Sprintf("FetchToken Response: [Code=%s,Message=%s,AccessToken=%s,TokenType=%s,ExpiresIn=%d seconds,Error=%s]",
		tr.Code, tr.Message, tr.AccessToken, tr.TokenType, tr.ExpiresIn, tr.Error)
}

//func TokenSaverMiddleware() TokenFetcherMiddleware {
//	// This is the actual middleware function to be executed.
//	m := func(handler TokenFetchFunc) TokenFetchFunc{
//		f := func(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error) {
//			response, err := handler(ctx, url, request)
//			if err != nil {
//				return nil, err
//			}
//
//			if response.AccessToken != "" {
//				// save the token
//				_ = SaveToken(response.AccessToken)
//			}
//
//			return response, nil
//		}
//
//		return f
//	}
//
//	return m
//}

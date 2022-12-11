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
	"time"

	xhttp "github.com/vfdcloud/vfd/internal/http"
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

	// OnTokenResponse is a callback function that is called when a token is received.
	OnTokenResponse func(context.Context, *TokenResponse) error

	// TokenResponseMiddleware is a middleware function that is called when a token is received.
	TokenResponseMiddleware func(next OnTokenResponse) OnTokenResponse
)

// WrapTokenResponseMiddleware wraps a TokenResponseMiddleware with a OnTokenResponse.
func WrapTokenResponseMiddleware(next OnTokenResponse, middlewares ...TokenResponseMiddleware) OnTokenResponse {
	// loop backwards through the middlewares
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

// FetchTokenWithMw retrieves a token from the VFD server then passes it to the callback function
// This is beacuse the response might have a code and message that needs to be handled.
func FetchTokenWithMw(ctx context.Context, url string, request *TokenRequest, callback OnTokenResponse) (*TokenResponse, error) {
	httpClient := xhttp.Instance()

	response, err := fetchToken(ctx, httpClient, url, request)
	if err != nil {
		return nil, err
	}

	err = callback(ctx, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// FetchToken retrieves a token from the VFD server. If the status code is not 200, an error is
// returned. Error Message will contain TokenResponse.Code and TokenResponse.Message
// FetchToken wraps internally a *http.Client responsible for making http calls. It has a timeout
// of 70 seconds. It is advised to call this only when the previous token has expired. It will still
// work if called before the token expires.
// It is a context-aware function with a timeout of 1 minute
func FetchToken(ctx context.Context, url string, request *TokenRequest) (*TokenResponse, error) {
	httpClient := xhttp.Instance()
	return fetchToken(ctx, httpClient, url, request)
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
	form := url.Values{}
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
		return nil, fmt.Errorf("http call error: %w: %v", ErrFetchToken, err)
	}
	defer resp.Body.Close()

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
	return fmt.Sprintf(
		"FetchToken Response: [Code=%s,Message=%s,AccessToken=%s,TokenType=%s,ExpiresIn=%d seconds,Error=%s]",
		tr.Code, tr.Message, tr.AccessToken, tr.TokenType, tr.ExpiresIn, tr.Error)
}

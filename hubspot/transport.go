package hubspot

import (
	"fmt"
	"net/http"
)

// authTransport injects the Authorization header using a TokenSource.
type authTransport struct {
	tokenSource TokenSource
	base        http.RoundTripper
}

// RoundTrip injects the Authorization header using a TokenSource.
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.tokenSource.Token(req.Context())
	if err != nil {
		return nil, fmt.Errorf("hubspot: auth: %w", err)
	}

	// Clone to avoid mutating the caller's request.
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("hubspot: execute request: %w", err)
	}

	return resp, nil
}

// headerTransport adds standard SDK headers to every request.
type headerTransport struct {
	userAgent string
	base      http.RoundTripper
}

// RoundTrip adds standard SDK headers to every request.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("hubspot: execute request: %w", err)
	}

	return resp, nil
}

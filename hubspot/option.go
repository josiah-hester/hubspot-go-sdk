package hubspot

import "net/http"

// Option configures a [Client].
type Option func(*Client)

// WithBaseURL overrides the default HubSpot API base URL.
// This is useful for testing against a mock server.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient overrides the default HTTP client. The provided client's
// Transport will be wrapped with auth and header transports — do not set
// Authorization headers on it directly.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent overrides the default User-Agent header sent with every
// request. The default is "hubspot-go-sdk/{version}".
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

package hubspot

import (
	"context"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.hubapi.com"

// Client is the top-level HubSpot API client. Create one with [NewClient]
// and pass it to service constructors like [crm.NewService].
//
//	client := hubspot.NewClient(hubspot.PrivateAppToken("pat-na1-xxxxx"))
//	crmService := crm.NewService(client)
//	contact, err := crmService.Contacts().Get(ctx, "123", nil)
//
// Client implements [Requester] and is safe for concurrent use.
type Client struct {
	baseURL   string
	userAgent string

	httpClient *http.Client

	tokenSource TokenSource

	req *requester
}

// NewClient creates a new HubSpot API client with the given [TokenSource]
// and options.
func NewClient(tokenSource TokenSource, opts ...Option) *Client {
	c := &Client{
		baseURL:   defaultBaseURL,
		userAgent: "hubspot-go-sdk/" + Version,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenSource: tokenSource,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Build the transport chain. Order matters — outermost runs first.
	// Final order: caller → headerTransport → authTransport → base transport.
	transport := c.httpClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	transport = &authTransport{
		tokenSource: c.tokenSource,
		base:        transport,
	}
	transport = &headerTransport{
		userAgent: c.userAgent,
		base:      transport,
	}

	// Replace the transport on a shallow copy so we don't mutate a
	// user-provided http.Client.
	hc := *c.httpClient
	hc.Transport = transport

	c.req = newRequester(&hc, c.baseURL)

	return c
}

// Do executes an API request. This implements [Requester], allowing Client
// to be passed directly to service constructors.
func (c *Client) Do(ctx context.Context, req *Request, result any) error {
	return c.req.Do(ctx, req, result)
}

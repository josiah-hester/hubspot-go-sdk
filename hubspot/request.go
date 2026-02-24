package hubspot

import (
	"context"
	"io"
	"net/url"
)

// Requester executes API requests. It is the internal interface passed to
// service packages so they do not depend on [Client] directly.
type Requester interface {
	// Do executes an API request and decodes the response into result.
	// If result is nil, the response body is discarded.
	Do(ctx context.Context, req *Request, result any) error
}

// Request describes an API request to be executed by a [Requester].
type Request struct {
	// Method is the HTTP method (GET, POST, PATCH, DELETE, etc.).
	Method string

	// Path is the URL path relative to the base URL (e.g., "/crm/v3/objects/contacts").
	Path string

	// Query contains URL query parameters.
	Query url.Values

	// Body is the request body, JSON-encoded by the requester.
	// Ignored for GET and DELETE requests. For non-JSON bodies, use RawBody instead.
	Body any

	// ContentType overrides the default "application/json" content type.
	ContentType string

	// RawBody provides a pre-encoded request body (e.g., for multipart uploads).
	// When set, Body is ignored.
	RawBody io.Reader
}

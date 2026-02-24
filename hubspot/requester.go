package hubspot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// requester is the internal [Requester] implementation that executes HTTP
// requests against the HubSpot API.
type requester struct {
	httpClient *http.Client
	baseURL    string
}

func newRequester(httpClient *http.Client, baseURL string) *requester {
	return &requester{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

// Do executes an API request and decodes the JSON response into result.
func (r *requester) Do(ctx context.Context, req *Request, result any) error {
	httpReq, err := r.buildHTTPRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("hubspot: build request: %w", err)
	}

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("hubspot: execute request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("hubspot: close response: %v", err)
		}
	}()

	// Read the full body so the connection can be reused.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("hubspot: read response: %w", err)
	}

	// Error responses.
	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, body, resp.Header)
	}

	// Success with no body expected.
	if result == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// Decode JSON response.
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("hubspot: decode response: %w", err)
	}

	return nil
}

func (r *requester) buildHTTPRequest(ctx context.Context, req *Request) (*http.Request, error) {
	url := r.baseURL + req.Path
	if len(req.Query) > 0 {
		url += "?" + req.Query.Encode()
	}

	var bodyReader io.Reader
	contentType := "application/json"

	if req.RawBody != nil {
		bodyReader = req.RawBody
		if req.ContentType != "" {
			contentType = req.ContentType
		}
	} else if req.Body != nil {
		data, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("encode body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("hubspot: new request: %w", err)
	}

	if bodyReader != nil {
		httpReq.Header.Set("Content-Type", contentType)
	}

	return httpReq, nil
}

func parseAPIError(statusCode int, body []byte, _ http.Header) error {
	apiErr := &APIError{
		StatusCode: statusCode,
	}

	// Try to decode HubSpot's JSON error body.
	if len(body) > 0 {
		if err := json.Unmarshal(body, apiErr); err != nil {
			// If we can't parse the body, still return a usable error.
			apiErr.Message = string(body)
		}
	}

	// Ensure StatusCode is always set (JSON doesn't include it).
	apiErr.StatusCode = statusCode

	return apiErr
}

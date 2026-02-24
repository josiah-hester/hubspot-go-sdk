package hubspot

import (
	"errors"
	"fmt"
)

// APIError represents an error response from the HubSpot API.
type APIError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int `json:"-"`

	// Status is the error status string (typically "error").
	Status string `json:"status"`

	// Message is a human-readable error description.
	Message string `json:"message"`

	// CorrelationID is a unique identifier for the request, useful when
	// contacting HubSpot support.
	CorrelationID string `json:"correlationId"`

	// Category is the error category (e.g., "RATE_LIMITS", "VALIDATION_ERROR").
	Category string `json:"category"`

	// SubCategory provides additional error classification.
	SubCategory string `json:"subCategory,omitempty"`

	// ErrorType is the specific error type (e.g., "RATE_LIMIT", "NOT_FOUND").
	ErrorType string `json:"errorType,omitempty"`

	// PolicyName indicates which rate limit was hit ("DAILY" or "SECONDLY").
	// Only present on 429 responses.
	PolicyName string `json:"policyName,omitempty"`

	// RequestID is the HubSpot request identifier.
	RequestID string `json:"requestId,omitempty"`

	// Context contains additional error context from the API.
	Context map[string][]string `json:"context,omitempty"`

	// Links contains related resource URLs.
	Links map[string]string `json:"links,omitempty"`
}

// Error returns a human-readable error string.
func (e *APIError) Error() string {
	if e.ErrorType != "" {
		return fmt.Sprintf("hubspot: %d %s (correlation_id=%s): %s",
			e.StatusCode, e.ErrorType, e.CorrelationID, e.Message)
	}
	return fmt.Sprintf("hubspot: %d (correlation_id=%s): %s",
		e.StatusCode, e.CorrelationID, e.Message)
}

// IsNotFound reports whether the error is a 404 Not Found response.
func IsNotFound(err error) bool {
	return hasStatusCode(err, 404)
}

// IsRateLimited reports whether the error is a 429 Too Many Requests response.
func IsRateLimited(err error) bool {
	return hasStatusCode(err, 429)
}

// IsUnauthorized reports whether the error is a 401 Unauthorized response.
func IsUnauthorized(err error) bool {
	return hasStatusCode(err, 401)
}

// IsForbidden reports whether the error is a 403 Forbidden response.
func IsForbidden(err error) bool {
	return hasStatusCode(err, 403)
}

// IsConflict reports whether the error is a 409 Conflict response.
func IsConflict(err error) bool {
	return hasStatusCode(err, 409)
}

// IsServerError reports whether the error is a 5xx server error.
func IsServerError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode >= 500
}

// IsDailyRateLimit reports whether the error is a 429 caused by hitting
// the daily API call quota.
func IsDailyRateLimit(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.PolicyName == "DAILY"
}

// IsBurstRateLimit reports whether the error is a 429 caused by hitting
// the per-second burst limit.
func IsBurstRateLimit(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.PolicyName == "SECONDLY"
}

func hasStatusCode(err error, code int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == code
}

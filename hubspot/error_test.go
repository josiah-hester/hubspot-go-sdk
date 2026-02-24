package hubspot_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  hubspot.APIError
		want string
	}{
		{
			name: "with error type",
			err: hubspot.APIError{
				StatusCode:    404,
				ErrorType:     "NOT_FOUND",
				CorrelationID: "abc-123",
				Message:       "Object not found",
			},
			want: "hubspot: 404 NOT_FOUND (correlation_id=abc-123): Object not found",
		},
		{
			name: "without error type",
			err: hubspot.APIError{
				StatusCode:    400,
				CorrelationID: "def-456",
				Message:       "Invalid input",
			},
			want: "hubspot: 400 (correlation_id=def-456): Invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError_Implements_error(t *testing.T) {
	var err error = &hubspot.APIError{StatusCode: 500, Message: "fail"}
	if err.Error() == "" {
		t.Error("expected non-empty error string")
	}
}

func TestSentinelHelpers(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		check  func(error) bool
		expect bool
	}{
		{"IsNotFound true", &hubspot.APIError{StatusCode: 404}, hubspot.IsNotFound, true},
		{"IsNotFound false", &hubspot.APIError{StatusCode: 400}, hubspot.IsNotFound, false},
		{"IsRateLimited true", &hubspot.APIError{StatusCode: 429}, hubspot.IsRateLimited, true},
		{"IsRateLimited false", &hubspot.APIError{StatusCode: 200}, hubspot.IsRateLimited, false},
		{"IsUnauthorized true", &hubspot.APIError{StatusCode: 401}, hubspot.IsUnauthorized, true},
		{"IsForbidden true", &hubspot.APIError{StatusCode: 403}, hubspot.IsForbidden, true},
		{"IsConflict true", &hubspot.APIError{StatusCode: 409}, hubspot.IsConflict, true},
		{"IsServerError 500", &hubspot.APIError{StatusCode: 500}, hubspot.IsServerError, true},
		{"IsServerError 503", &hubspot.APIError{StatusCode: 503}, hubspot.IsServerError, true},
		{"IsServerError 400", &hubspot.APIError{StatusCode: 400}, hubspot.IsServerError, false},
		{
			"IsDailyRateLimit true",
			&hubspot.APIError{StatusCode: 429, PolicyName: "DAILY"},
			hubspot.IsDailyRateLimit,
			true,
		},
		{
			"IsDailyRateLimit wrong policy",
			&hubspot.APIError{StatusCode: 429, PolicyName: "SECONDLY"},
			hubspot.IsDailyRateLimit,
			false,
		},
		{
			"IsBurstRateLimit true",
			&hubspot.APIError{StatusCode: 429, PolicyName: "SECONDLY"},
			hubspot.IsBurstRateLimit,
			true,
		},
		{
			"non-APIError returns false",
			fmt.Errorf("some other error"),
			hubspot.IsNotFound,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.check(tt.err)
			if got != tt.expect {
				t.Errorf("got %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestSentinelHelpers_WrappedError(t *testing.T) {
	// Sentinels should work through fmt.Errorf wrapping.
	inner := &hubspot.APIError{StatusCode: 404, Message: "not found"}
	wrapped := fmt.Errorf("outer: %w", inner)

	if !hubspot.IsNotFound(wrapped) {
		t.Error("IsNotFound should detect wrapped APIError")
	}

	var apiErr *hubspot.APIError
	if !errors.As(wrapped, &apiErr) {
		t.Error("errors.As should unwrap to *APIError")
	}
}

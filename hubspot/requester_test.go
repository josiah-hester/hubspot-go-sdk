package hubspot_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

func TestRequester_SuccessfulGET(t *testing.T) {
	type response struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/123" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/123", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response{ID: "123", Name: "Alice"})
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	var got response
	err := client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/crm/v3/objects/contacts/123",
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "123" || got.Name != "Alice" {
		t.Errorf("got %+v, want {ID:123 Name:Alice}", got)
	}
}

func TestRequester_QueryParameters(t *testing.T) {
	var gotQuery url.Values

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	err := client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
		Query: url.Values{
			"properties": {"email,firstname"},
			"archived":   {"true"},
		},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery.Get("properties") != "email,firstname" {
		t.Errorf("properties = %q, want %q", gotQuery.Get("properties"), "email,firstname")
	}
	if gotQuery.Get("archived") != "true" {
		t.Errorf("archived = %q, want %q", gotQuery.Get("archived"), "true")
	}
}

func TestRequester_POSTWithBody(t *testing.T) {
	type input struct {
		Properties map[string]string `json:"properties"`
	}

	var gotBody input

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}

		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"456"}`))
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	type result struct {
		ID string `json:"id"`
	}

	var got result
	err := client.Do(context.Background(), &hubspot.Request{
		Method: "POST",
		Path:   "/crm/v3/objects/contacts",
		Body:   input{Properties: map[string]string{"email": "a@b.com"}},
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "456" {
		t.Errorf("ID = %q, want 456", got.ID)
	}
	if gotBody.Properties["email"] != "a@b.com" {
		t.Errorf("body email = %q, want a@b.com", gotBody.Properties["email"])
	}
}

func TestRequester_NoContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	err := client.Do(context.Background(), &hubspot.Request{
		Method: "DELETE",
		Path:   "/crm/v3/objects/contacts/123",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequester_APIError_JSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"status":        "error",
			"message":       "Object not found. objectId=999 objectType=contacts",
			"correlationId": "abc-123-def",
			"category":      "OBJECT_NOT_FOUND",
		})
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	err := client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/crm/v3/objects/contacts/999",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !hubspot.IsNotFound(err) {
		t.Errorf("IsNotFound = false, want true; err = %v", err)
	}

	var apiErr *hubspot.APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *APIError")
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if apiErr.CorrelationID != "abc-123-def" {
		t.Errorf("CorrelationID = %q, want abc-123-def", apiErr.CorrelationID)
	}
	if apiErr.Category != "OBJECT_NOT_FOUND" {
		t.Errorf("Category = %q, want OBJECT_NOT_FOUND", apiErr.Category)
	}
}

func TestRequester_APIError_RateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "10")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{
			"status":     "error",
			"message":    "You have reached your secondly limit.",
			"errorType":  "RATE_LIMIT",
			"policyName": "SECONDLY",
		})
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	err := client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if !hubspot.IsRateLimited(err) {
		t.Errorf("IsRateLimited = false, want true; err = %v", err)
	}
	if !hubspot.IsBurstRateLimit(err) {
		t.Errorf("IsBurstRateLimit = false, want true")
	}
	if hubspot.IsDailyRateLimit(err) {
		t.Error("IsDailyRateLimit = true, want false")
	}
}

func TestRequester_APIError_NonJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Bad Gateway"))
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	err := client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *hubspot.APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *APIError")
	}
	if apiErr.StatusCode != 502 {
		t.Errorf("StatusCode = %d, want 502", apiErr.StatusCode)
	}
	if apiErr.Message != "Bad Gateway" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "Bad Gateway")
	}
}

func TestRequester_ContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	err := client.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}

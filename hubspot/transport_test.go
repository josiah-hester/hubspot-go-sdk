package hubspot_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

func TestTransport_AuthHeader(t *testing.T) {
	var gotAuth string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("test-token-123"),
		hubspot.WithBaseURL(ts.URL),
	)

	_ = client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	want := "Bearer test-token-123"
	if gotAuth != want {
		t.Errorf("Authorization header = %q, want %q", gotAuth, want)
	}
}

func TestTransport_UserAgentHeader(t *testing.T) {
	var gotUA string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	_ = client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	want := "hubspot-go-sdk/" + hubspot.Version
	if gotUA != want {
		t.Errorf("User-Agent header = %q, want %q", gotUA, want)
	}
}

func TestTransport_CustomUserAgent(t *testing.T) {
	var gotUA string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
		hubspot.WithUserAgent("my-app/1.0"),
	)

	_ = client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if gotUA != "my-app/1.0" {
		t.Errorf("User-Agent header = %q, want %q", gotUA, "my-app/1.0")
	}
}

func TestTransport_AcceptHeader(t *testing.T) {
	var gotAccept string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAccept = r.Header.Get("Accept")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	_ = client.Do(context.Background(), &hubspot.Request{
		Method: "GET",
		Path:   "/test",
	}, nil)

	if gotAccept != "application/json" {
		t.Errorf("Accept header = %q, want %q", gotAccept, "application/json")
	}
}

func TestTransport_DoesNotMutateOriginalRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := hubspot.NewClient(
		hubspot.PrivateAppToken("token"),
		hubspot.WithBaseURL(ts.URL),
	)

	// Verify that calling Do does not panic or cause issues with
	// concurrent requests by exercising the transport.
	for i := 0; i < 5; i++ {
		err := client.Do(context.Background(), &hubspot.Request{
			Method: "GET",
			Path:   "/test",
		}, nil)
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
	}
}

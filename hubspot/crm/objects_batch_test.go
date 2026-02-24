package crm_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot/crm"
)

func batchResponseJSON(ids ...string) map[string]any {
	results := make([]any, len(ids))
	for i, id := range ids {
		results[i] = contactJSON(id, map[string]string{})
	}
	return map[string]any{
		"status":      "COMPLETE",
		"results":     results,
		"startedAt":   "2024-01-01T00:00:00.000Z",
		"completedAt": "2024-01-01T00:00:01.000Z",
	}
}

func TestBatchService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/batch/create" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/batch/create", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.BatchCreateInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)

		if len(payload.Inputs) != 2 {
			t.Errorf("len(inputs) = %d, want 2", len(payload.Inputs))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(batchResponseJSON("1", "2"))
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Contacts().Batch()
	resp, err := batch.Create(context.Background(), []crm.BatchCreateInput{
		{Properties: map[string]string{"email": "a@b.com"}},
		{Properties: map[string]string{"email": "c@d.com"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Errorf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Status != "COMPLETE" {
		t.Errorf("Status = %q, want COMPLETE", resp.Status)
	}
}

func TestBatchService_Read(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/crm/v3/objects/contacts/batch/read" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/batch/read", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.BatchReadInput
		json.Unmarshal(body, &input)

		if len(input.Inputs) != 2 {
			t.Errorf("len(inputs) = %d, want 2", len(input.Inputs))
		}
		if len(input.Properties) != 1 || input.Properties[0] != "email" {
			t.Errorf("properties = %v, want [email]", input.Properties)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(batchResponseJSON("10", "20"))
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Contacts().Batch()
	resp, err := batch.Read(context.Background(), &crm.BatchReadInput{
		Properties: []string{"email"},
		Inputs:     []crm.BatchReadID{{ID: "10"}, {ID: "20"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Errorf("len(Results) = %d, want 2", len(resp.Results))
	}
}

func TestBatchService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/crm/v3/objects/contacts/batch/update" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/batch/update", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.BatchUpdateInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)

		if len(payload.Inputs) != 1 {
			t.Errorf("len(inputs) = %d, want 1", len(payload.Inputs))
		}
		if payload.Inputs[0].ID != "5" {
			t.Errorf("ID = %q, want 5", payload.Inputs[0].ID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(batchResponseJSON("5"))
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Contacts().Batch()
	resp, err := batch.Update(context.Background(), []crm.BatchUpdateInput{
		{ID: "5", Properties: map[string]string{"lastname": "New"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want 1", len(resp.Results))
	}
}

func TestBatchService_Upsert(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/crm/v3/objects/contacts/batch/upsert" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/batch/upsert", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.BatchUpsertInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)

		if len(payload.Inputs) != 1 {
			t.Errorf("len(inputs) = %d, want 1", len(payload.Inputs))
		}
		if payload.Inputs[0].IDProperty != "email" {
			t.Errorf("IDProperty = %q, want email", payload.Inputs[0].IDProperty)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(batchResponseJSON("100"))
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Contacts().Batch()
	resp, err := batch.Upsert(context.Background(), []crm.BatchUpsertInput{
		{
			IDProperty: "email",
			ID:         "upsert@example.com",
			Properties: map[string]string{"firstname": "Upserted"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want 1", len(resp.Results))
	}
}

func TestBatchService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/batch/archive" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/batch/archive", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.BatchArchiveInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)

		if len(payload.Inputs) != 3 {
			t.Errorf("len(inputs) = %d, want 3", len(payload.Inputs))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Contacts().Batch()
	err := batch.Archive(context.Background(), []crm.BatchArchiveInput{
		{ID: "1"}, {ID: "2"}, {ID: "3"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBatchService_DifferentObjectType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/crm/v3/objects/deals/batch/create" {
			t.Errorf("path = %s, want /crm/v3/objects/deals/batch/create", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(batchResponseJSON("1"))
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Deals().Batch()
	_, err := batch.Create(context.Background(), []crm.BatchCreateInput{
		{Properties: map[string]string{"dealname": "Big Deal"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

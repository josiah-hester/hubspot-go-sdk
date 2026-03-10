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

func TestPropertyBatchService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/properties/contacts/batch/create"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.PropertyCreateInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)
		if len(payload.Inputs) != 2 {
			t.Errorf("len(inputs) = %d, want 2", len(payload.Inputs))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "COMPLETE",
			"results": []map[string]any{
				{"name": "prop_a", "label": "Prop A", "type": "string", "fieldType": "text", "groupName": "contactinformation"},
				{"name": "prop_b", "label": "Prop B", "type": "number", "fieldType": "number", "groupName": "contactinformation"},
			},
			"startedAt":   "2024-01-01T00:00:00.000Z",
			"completedAt": "2024-01-01T00:00:01.000Z",
		})
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Properties("contacts").Batch()
	resp, err := batch.Create(context.Background(), []crm.PropertyCreateInput{
		{Name: "prop_a", Label: "Prop A", Type: "string", FieldType: "text", GroupName: "contactinformation"},
		{Name: "prop_b", Label: "Prop B", Type: "number", FieldType: "number", GroupName: "contactinformation"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "COMPLETE" {
		t.Errorf("Status = %q, want COMPLETE", resp.Status)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Results[0].Name != "prop_a" {
		t.Errorf("Results[0].Name = %q, want prop_a", resp.Results[0].Name)
	}
}

func TestPropertyBatchService_Read(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/properties/contacts/batch/read"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs   []crm.PropertyBatchReadInput `json:"inputs"`
			Archived bool                         `json:"archived"`
		}
		json.Unmarshal(body, &payload)
		if len(payload.Inputs) != 2 {
			t.Errorf("len(inputs) = %d, want 2", len(payload.Inputs))
		}
		if payload.Archived {
			t.Error("archived = true, want false")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "COMPLETE",
			"results": []map[string]any{
				{"name": "email", "label": "Email", "type": "string", "fieldType": "text", "groupName": "contactinformation"},
				{"name": "firstname", "label": "First Name", "type": "string", "fieldType": "text", "groupName": "contactinformation"},
			},
			"startedAt":   "2024-01-01T00:00:00.000Z",
			"completedAt": "2024-01-01T00:00:01.000Z",
		})
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Properties("contacts").Batch()
	resp, err := batch.Read(context.Background(), &crm.PropertyBatchReadRequest{
		Inputs: []crm.PropertyBatchReadInput{
			{Name: "email"},
			{Name: "firstname"},
		},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "COMPLETE" {
		t.Errorf("Status = %q, want COMPLETE", resp.Status)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
}

func TestPropertyBatchService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/properties/contacts/batch/archive"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.PropertyBatchReadInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)
		if len(payload.Inputs) != 1 {
			t.Errorf("len(inputs) = %d, want 1", len(payload.Inputs))
		}
		if payload.Inputs[0].Name != "old_prop" {
			t.Errorf("inputs[0].name = %q, want old_prop", payload.Inputs[0].Name)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Properties("contacts").Batch()
	err := batch.Archive(context.Background(), []crm.PropertyBatchReadInput{
		{Name: "old_prop"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

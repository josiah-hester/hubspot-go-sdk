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

func TestAssociationsBatchService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/batch/create"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.AssociationBatchCreateInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)
		if len(payload.Inputs) != 2 {
			t.Errorf("len(inputs) = %d, want 2", len(payload.Inputs))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "COMPLETE",
			"results": []map[string]any{
				{
					"fromObjectTypeId": "0-1",
					"fromObjectId":     "c1",
					"toObjectTypeId":   "0-2",
					"toObjectId":       "co1",
					"labels":           []map[string]any{{"category": "HUBSPOT_DEFINED", "typeId": 1}},
				},
			},
			"startedAt":   "2024-01-01T00:00:00.000Z",
			"completedAt": "2024-01-01T00:00:01.000Z",
		})
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Batch()
	resp, err := batch.Create(context.Background(), []crm.AssociationBatchCreateInput{
		{
			From:  crm.AssociationObjectID{ID: "c1"},
			To:    crm.AssociationObjectID{ID: "co1"},
			Types: []crm.AssociationTypeInput{{AssociationCategory: crm.AssociationCategoryHubSpotDefined, AssociationTypeID: 1}},
		},
		{
			From:  crm.AssociationObjectID{ID: "c2"},
			To:    crm.AssociationObjectID{ID: "co2"},
			Types: []crm.AssociationTypeInput{{AssociationCategory: crm.AssociationCategoryHubSpotDefined, AssociationTypeID: 1}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "COMPLETE" {
		t.Errorf("Status = %q, want COMPLETE", resp.Status)
	}
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want 1", len(resp.Results))
	}
}

func TestAssociationsBatchService_Read(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/batch/read"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.AssociationBatchReadInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)
		if len(payload.Inputs) != 2 {
			t.Errorf("len(inputs) = %d, want 2", len(payload.Inputs))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "COMPLETE",
			"results": []map[string]any{
				{
					"from": map[string]string{"id": "c1"},
					"to": []map[string]any{
						{
							"toObjectId":       "co1",
							"associationTypes": []map[string]any{{"category": "HUBSPOT_DEFINED", "typeId": 1}},
						},
					},
				},
			},
			"startedAt":   "2024-01-01T00:00:00.000Z",
			"completedAt": "2024-01-01T00:00:01.000Z",
		})
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Batch()
	resp, err := batch.Read(context.Background(), []crm.AssociationBatchReadInput{
		{ID: "c1"},
		{ID: "c2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "COMPLETE" {
		t.Errorf("Status = %q, want COMPLETE", resp.Status)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(resp.Results))
	}
	if resp.Results[0].From.ID != "c1" {
		t.Errorf("From.ID = %q, want c1", resp.Results[0].From.ID)
	}
	if len(resp.Results[0].To) != 1 {
		t.Fatalf("len(To) = %d, want 1", len(resp.Results[0].To))
	}
	if resp.Results[0].To[0].ToObjectID != "co1" {
		t.Errorf("ToObjectID = %q, want co1", resp.Results[0].To[0].ToObjectID)
	}
}

func TestAssociationsBatchService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/batch/archive"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Inputs []crm.AssociationBatchArchiveInput `json:"inputs"`
		}
		json.Unmarshal(body, &payload)
		if len(payload.Inputs) != 1 {
			t.Errorf("len(inputs) = %d, want 1", len(payload.Inputs))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	batch := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Batch()
	err := batch.Archive(context.Background(), []crm.AssociationBatchArchiveInput{
		{
			From: crm.AssociationObjectID{ID: "c1"},
			To:   crm.AssociationObjectID{ID: "co1"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

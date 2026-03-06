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

func TestAssociationsService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if want := "/crm/v4/objects/contacts/contact-1/associations/companies/company-1"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var types []crm.AssociationTypeInput
		json.Unmarshal(body, &types)
		if len(types) != 1 {
			t.Fatalf("len(types) = %d, want 1", len(types))
		}
		if types[0].AssociationCategory != crm.AssociationCategoryHubSpotDefined {
			t.Errorf("category = %q, want %q", types[0].AssociationCategory, crm.AssociationCategoryHubSpotDefined)
		}
		if types[0].AssociationTypeID != 1 {
			t.Errorf("typeID = %d, want 1", types[0].AssociationTypeID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"fromObjectTypeId": "0-1",
			"fromObjectId":     "contact-1",
			"toObjectTypeId":   "0-2",
			"toObjectId":       "company-1",
			"labels": []map[string]any{
				{"category": "HUBSPOT_DEFINED", "typeId": 1, "label": "Primary"},
			},
		})
	}))
	defer ts.Close()

	assoc := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies")
	result, err := assoc.Create(context.Background(), "contact-1", "company-1", []crm.AssociationTypeInput{
		{AssociationCategory: crm.AssociationCategoryHubSpotDefined, AssociationTypeID: 1},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.FromObjectID != "contact-1" {
		t.Errorf("FromObjectID = %q, want contact-1", result.FromObjectID)
	}
	if result.ToObjectID != "company-1" {
		t.Errorf("ToObjectID = %q, want company-1", result.ToObjectID)
	}
	if len(result.Labels) != 1 {
		t.Fatalf("len(Labels) = %d, want 1", len(result.Labels))
	}
	if result.Labels[0].Label != "Primary" {
		t.Errorf("Label = %q, want Primary", result.Labels[0].Label)
	}
}

func TestAssociationsService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v4/objects/contacts/contact-1/associations/companies"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"toObjectId": "company-1",
					"associationTypes": []map[string]any{
						{"category": "HUBSPOT_DEFINED", "typeId": 1, "label": "Primary"},
					},
				},
				{
					"toObjectId": "company-2",
					"associationTypes": []map[string]any{
						{"category": "HUBSPOT_DEFINED", "typeId": 1},
					},
				},
			},
			"paging": map[string]any{
				"next": map[string]string{"after": "cursor-123"},
			},
		})
	}))
	defer ts.Close()

	assoc := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies")
	resp, err := assoc.List(context.Background(), "contact-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Results[0].ToObjectID != "company-1" {
		t.Errorf("ToObjectID = %q, want company-1", resp.Results[0].ToObjectID)
	}
	if len(resp.Results[0].AssociationTypes) != 1 {
		t.Fatalf("len(AssociationTypes) = %d, want 1", len(resp.Results[0].AssociationTypes))
	}
	if resp.Paging == nil || resp.Paging.Next == nil {
		t.Fatal("expected paging cursor")
	}
	if resp.Paging.Next.After != "cursor-123" {
		t.Errorf("After = %q, want cursor-123", resp.Paging.Next.After)
	}
}

func TestAssociationsService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if want := "/crm/v4/objects/contacts/contact-1/associations/companies/company-1"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	assoc := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies")
	err := assoc.Archive(context.Background(), "contact-1", "company-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssociationsService_CreateDefault(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if want := "/crm/v4/objects/contacts/contact-1/associations/default/companies/company-1"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			if len(body) > 0 {
				t.Errorf("expected no request body, got %s", body)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "COMPLETE",
			"results": []map[string]any{
				{
					"from":            map[string]string{"id": "contact-1"},
					"to":              map[string]string{"id": "company-1"},
					"associationSpec": map[string]any{"associationCategory": "HUBSPOT_DEFINED", "associationTypeId": 1},
				},
			},
			"startedAt":   "2024-01-01T00:00:00.000Z",
			"completedAt": "2024-01-01T00:00:01.000Z",
		})
	}))
	defer ts.Close()

	assoc := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies")
	resp, err := assoc.CreateDefault(context.Background(), "contact-1", "company-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "COMPLETE" {
		t.Errorf("Status = %q, want COMPLETE", resp.Status)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(resp.Results))
	}
	if resp.Results[0].From.ID != "contact-1" {
		t.Errorf("From.ID = %q, want contact-1", resp.Results[0].From.ID)
	}
	if resp.Results[0].To.ID != "company-1" {
		t.Errorf("To.ID = %q, want company-1", resp.Results[0].To.ID)
	}
	if resp.Results[0].AssociationSpec.AssociationCategory != crm.AssociationCategoryHubSpotDefined {
		t.Errorf("AssociationCategory = %q, want %q", resp.Results[0].AssociationSpec.AssociationCategory, crm.AssociationCategoryHubSpotDefined)
	}
	if resp.Results[0].AssociationSpec.AssociationTypeID != 1 {
		t.Errorf("AssociationTypeID = %d, want 1", resp.Results[0].AssociationSpec.AssociationTypeID)
	}
}

func TestAssociationsService_CustomObjectTypes(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
	defer ts.Close()

	assoc := crm.NewService(newTestClient(t, ts)).Associations("2-12345", "deals")
	_, err := assoc.List(context.Background(), "obj-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "/crm/v4/objects/2-12345/obj-1/associations/deals"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
}

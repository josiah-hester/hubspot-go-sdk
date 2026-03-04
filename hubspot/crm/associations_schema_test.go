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

func TestAssociationsSchemaService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/labels"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"category": "HUBSPOT_DEFINED", "typeId": 1, "label": "Primary"},
				{"category": "USER_DEFINED", "typeId": 42, "label": "Partner"},
			},
		})
	}))
	defer ts.Close()

	schema := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Schema()
	resp, err := schema.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Results[0].Label != "Primary" {
		t.Errorf("Label = %q, want Primary", resp.Results[0].Label)
	}
	if resp.Results[1].Category != crm.AssociationCategoryUserDefined {
		t.Errorf("Category = %q, want %q", resp.Results[1].Category, crm.AssociationCategoryUserDefined)
	}
}

func TestAssociationsSchemaService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/labels"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.CreateAssociationLabelInput
		json.Unmarshal(body, &input)
		if input.Label != "Partner" {
			t.Errorf("Label = %q, want Partner", input.Label)
		}
		if input.Name != "partner" {
			t.Errorf("Name = %q, want partner", input.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"category": "USER_DEFINED", "typeId": 42, "label": "Partner"},
			},
		})
	}))
	defer ts.Close()

	schema := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Schema()
	resp, err := schema.Create(context.Background(), &crm.CreateAssociationLabelInput{
		Label: "Partner",
		Name:  "partner",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(resp.Results))
	}
	if resp.Results[0].TypeID != 42 {
		t.Errorf("TypeID = %d, want 42", resp.Results[0].TypeID)
	}
}

func TestAssociationsSchemaService_Create_WithInverseLabel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var input crm.CreateAssociationLabelInput
		json.Unmarshal(body, &input)
		if input.InverseLabel != "Client" {
			t.Errorf("InverseLabel = %q, want Client", input.InverseLabel)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"category": "USER_DEFINED", "typeId": 43, "label": "Vendor"},
				{"category": "USER_DEFINED", "typeId": 44, "label": "Client"},
			},
		})
	}))
	defer ts.Close()

	schema := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Schema()
	resp, err := schema.Create(context.Background(), &crm.CreateAssociationLabelInput{
		Label:        "Vendor",
		InverseLabel: "Client",
		Name:         "vendor_client",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
}

func TestAssociationsSchemaService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/labels"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.UpdateAssociationLabelInput
		json.Unmarshal(body, &input)
		if input.AssociationTypeID != 42 {
			t.Errorf("AssociationTypeID = %d, want 42", input.AssociationTypeID)
		}
		if input.Label != "New Label" {
			t.Errorf("Label = %q, want New Label", input.Label)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	schema := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Schema()
	err := schema.Update(context.Background(), &crm.UpdateAssociationLabelInput{
		AssociationTypeID: 42,
		Label:             "New Label",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssociationsSchemaService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if want := "/crm/v4/associations/contacts/companies/labels/42"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	schema := crm.NewService(newTestClient(t, ts)).Associations("contacts", "companies").Schema()
	err := schema.Delete(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

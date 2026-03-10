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

func TestPropertyGroupsService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v3/properties/contacts/groups"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"name": "contactinformation", "label": "Contact Information", "displayOrder": 0, "archived": false},
				{"name": "sales_properties", "label": "Sales Properties", "displayOrder": 1, "archived": false},
			},
		})
	}))
	defer ts.Close()

	groups := crm.NewService(newTestClient(t, ts)).Properties("contacts").Groups()
	resp, err := groups.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Results[0].Name != "contactinformation" {
		t.Errorf("Name = %q, want contactinformation", resp.Results[0].Name)
	}
}

func TestPropertyGroupsService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v3/properties/contacts/groups/contactinformation"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"name": "contactinformation", "label": "Contact Information", "displayOrder": 0, "archived": false,
		})
	}))
	defer ts.Close()

	groups := crm.NewService(newTestClient(t, ts)).Properties("contacts").Groups()
	group, err := groups.Get(context.Background(), "contactinformation", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "contactinformation" {
		t.Errorf("Name = %q, want contactinformation", group.Name)
	}
	if group.Label != "Contact Information" {
		t.Errorf("Label = %q, want Contact Information", group.Label)
	}
}

func TestPropertyGroupsService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/properties/contacts/groups"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		json.Unmarshal(body, &input)
		if input["name"] != "custom_group" {
			t.Errorf("name = %v, want custom_group", input["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"name": "custom_group", "label": "Custom Group", "displayOrder": 5, "archived": false,
		})
	}))
	defer ts.Close()

	groups := crm.NewService(newTestClient(t, ts)).Properties("contacts").Groups()
	group, err := groups.Create(context.Background(), &crm.PropertyGroupCreateInput{
		Name:  "custom_group",
		Label: "Custom Group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "custom_group" {
		t.Errorf("Name = %q, want custom_group", group.Name)
	}
}

func TestPropertyGroupsService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		if want := "/crm/v3/properties/contacts/groups/custom_group"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"name": "custom_group", "label": "Updated Group", "displayOrder": 10, "archived": false,
		})
	}))
	defer ts.Close()

	groups := crm.NewService(newTestClient(t, ts)).Properties("contacts").Groups()
	group, err := groups.Update(context.Background(), "custom_group", &crm.PropertyGroupUpdateInput{
		Label: "Updated Group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Label != "Updated Group" {
		t.Errorf("Label = %q, want Updated Group", group.Label)
	}
}

func TestPropertyGroupsService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if want := "/crm/v3/properties/contacts/groups/custom_group"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	groups := crm.NewService(newTestClient(t, ts)).Properties("contacts").Groups()
	err := groups.Archive(context.Background(), "custom_group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

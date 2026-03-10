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

func TestPropertiesService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v3/properties/contacts"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"name":      "email",
					"label":     "Email",
					"type":      "string",
					"fieldType": "text",
					"groupName": "contactinformation",
				},
				{
					"name":      "firstname",
					"label":     "First Name",
					"type":      "string",
					"fieldType": "text",
					"groupName": "contactinformation",
				},
			},
		})
	}))
	defer ts.Close()

	props := crm.NewService(newTestClient(t, ts)).Properties("contacts")
	resp, err := props.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Results[0].Name != "email" {
		t.Errorf("Name = %q, want email", resp.Results[0].Name)
	}
	if resp.Results[0].GroupName != "contactinformation" {
		t.Errorf("GroupName = %q, want contactinformation", resp.Results[0].GroupName)
	}
}

func TestPropertiesService_List_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("archived") != "true" {
			t.Errorf("archived = %q, want true", r.URL.Query().Get("archived"))
		}
		if r.URL.Query().Get("dataSensitivity") != "sensitive" {
			t.Errorf("dataSensitivity = %q, want sensitive", r.URL.Query().Get("dataSensitivity"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
	defer ts.Close()

	props := crm.NewService(newTestClient(t, ts)).Properties("contacts")
	_, err := props.List(context.Background(), &crm.PropertyListOptions{Archived: true, DataSensitivity: "sensitive"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPropertiesService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v3/properties/contacts/email"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"name":           "email",
			"label":          "Email",
			"type":           "string",
			"fieldType":      "text",
			"groupName":      "contactinformation",
			"hasUniqueValue": true,
			"hubspotDefined": true,
		})
	}))
	defer ts.Close()

	props := crm.NewService(newTestClient(t, ts)).Properties("contacts")
	prop, err := props.Get(context.Background(), "email", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prop.Name != "email" {
		t.Errorf("Name = %q, want email", prop.Name)
	}
	if !prop.HasUniqueValue {
		t.Error("HasUniqueValue = false, want true")
	}
	if !prop.HubSpotDefined {
		t.Error("HubSpotDefined = false, want true")
	}
}

func TestPropertiesService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/properties/contacts"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		json.Unmarshal(body, &input)
		if input["name"] != "favorite_food" {
			t.Errorf("name = %v, want favorite_food", input["name"])
		}
		if input["type"] != "string" {
			t.Errorf("type = %v, want string", input["type"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"name":      "favorite_food",
			"label":     "Favorite Food",
			"type":      "string",
			"fieldType": "text",
			"groupName": "contactinformation",
		})
	}))
	defer ts.Close()

	props := crm.NewService(newTestClient(t, ts)).Properties("contacts")
	prop, err := props.Create(context.Background(), &crm.PropertyCreateInput{
		Name:      "favorite_food",
		Label:     "Favorite Food",
		Type:      "string",
		FieldType: "text",
		GroupName: "contactinformation",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prop.Name != "favorite_food" {
		t.Errorf("Name = %q, want favorite_food", prop.Name)
	}
}

func TestPropertiesService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		if want := "/crm/v3/properties/contacts/favorite_food"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		json.Unmarshal(body, &input)
		if input["label"] != "Preferred Food" {
			t.Errorf("label = %v, want Preferred Food", input["label"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"name":      "favorite_food",
			"label":     "Preferred Food",
			"type":      "string",
			"fieldType": "text",
			"groupName": "contactinformation",
		})
	}))
	defer ts.Close()

	props := crm.NewService(newTestClient(t, ts)).Properties("contacts")
	prop, err := props.Update(context.Background(), "favorite_food", &crm.PropertyUpdateInput{
		Label: "Preferred Food",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prop.Label != "Preferred Food" {
		t.Errorf("Label = %q, want Preferred Food", prop.Label)
	}
}

func TestPropertiesService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if want := "/crm/v3/properties/contacts/favorite_food"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	props := crm.NewService(newTestClient(t, ts)).Properties("contacts")
	err := props.Archive(context.Background(), "favorite_food")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

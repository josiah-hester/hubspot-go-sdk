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

func TestSchemasService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v3/schemas"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{
					"id":                 "2-123456",
					"name":               "cars",
					"labels":             map[string]string{"singular": "Car", "plural": "Cars"},
					"requiredProperties": []string{"vin"},
					"properties":         []any{},
					"associations":       []any{},
				},
			},
		})
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	resp, err := schemas.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(resp.Results))
	}
	if resp.Results[0].ID != "2-123456" {
		t.Errorf("ID = %q, want 2-123456", resp.Results[0].ID)
	}
	if resp.Results[0].Name != "cars" {
		t.Errorf("Name = %q, want cars", resp.Results[0].Name)
	}
}

func TestSchemasService_List_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("archived") != "true" {
			t.Errorf("archived = %q, want true", r.URL.Query().Get("archived"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	_, err := schemas.List(context.Background(), &crm.SchemaListOptions{Archived: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchemasService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if want := "/crm/v3/schemas/my_object"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":                     "2-123456",
			"name":                   "my_object",
			"labels":                 map[string]string{"singular": "My Object", "plural": "My Objects"},
			"primaryDisplayProperty": "name",
			"requiredProperties":     []string{"name"},
			"properties": []map[string]any{
				{
					"name":      "name",
					"label":     "Name",
					"type":      "string",
					"fieldType": "text",
				},
			},
			"associations": []map[string]any{
				{
					"id":               "100",
					"fromObjectTypeId": "2-123456",
					"toObjectTypeId":   "0-1",
				},
			},
		})
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	schema, err := schemas.Get(context.Background(), "my_object", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if schema.ID != "2-123456" {
		t.Errorf("ID = %q, want 2-123456", schema.ID)
	}
	if schema.Name != "my_object" {
		t.Errorf("Name = %q, want my_object", schema.Name)
	}
	if len(schema.Properties) != 1 {
		t.Fatalf("len(Properties) = %d, want 1", len(schema.Properties))
	}
	if schema.Properties[0].Name != "name" {
		t.Errorf("Properties[0].Name = %q, want name", schema.Properties[0].Name)
	}
	if len(schema.Associations) != 1 {
		t.Fatalf("len(Associations) = %d, want 1", len(schema.Associations))
	}
	if schema.Associations[0].ID != "100" {
		t.Errorf("Associations[0].ID = %q, want 100", schema.Associations[0].ID)
	}
}

func TestSchemasService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/schemas"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		json.Unmarshal(body, &input)
		if input["name"] != "cars" {
			t.Errorf("name = %v, want cars", input["name"])
		}
		labels, _ := input["labels"].(map[string]any)
		if labels["singular"] != "Car" {
			t.Errorf("labels.singular = %v, want Car", labels["singular"])
		}
		props, _ := input["properties"].([]any)
		if len(props) != 1 {
			t.Fatalf("len(properties) = %d, want 1", len(props))
		}
		assocObjs, _ := input["associatedObjects"].([]any)
		if len(assocObjs) != 1 || assocObjs[0] != "CONTACT" {
			t.Errorf("associatedObjects = %v, want [CONTACT]", input["associatedObjects"])
		}
		reqProps, _ := input["requiredProperties"].([]any)
		if len(reqProps) != 1 || reqProps[0] != "vin" {
			t.Errorf("requiredProperties = %v, want [vin]", input["requiredProperties"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":                     "2-999999",
			"name":                   "cars",
			"labels":                 map[string]string{"singular": "Car", "plural": "Cars"},
			"primaryDisplayProperty": "vin",
			"requiredProperties":     []string{"vin"},
			"properties": []map[string]any{
				{"name": "vin", "label": "VIN", "type": "string", "fieldType": "text"},
			},
			"associations": []any{},
		})
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	schema, err := schemas.Create(context.Background(), &crm.SchemaCreateInput{
		Name:   "cars",
		Labels: crm.SchemaLabels{Singular: "Car", Plural: "Cars"},
		Properties: []crm.SchemaPropertyCreate{
			{Name: "vin", Label: "VIN", Type: "string", FieldType: "text", HasUniqueValue: true},
		},
		AssociatedObjects:      []string{"CONTACT"},
		RequiredProperties:     []string{"vin"},
		PrimaryDisplayProperty: "vin",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if schema.ID != "2-999999" {
		t.Errorf("ID = %q, want 2-999999", schema.ID)
	}
	if schema.Name != "cars" {
		t.Errorf("Name = %q, want cars", schema.Name)
	}
}

func TestSchemasService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		if want := "/crm/v3/schemas/2-123456"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		json.Unmarshal(body, &input)
		secDisplay, _ := input["secondaryDisplayProperties"].([]any)
		if len(secDisplay) != 2 {
			t.Errorf("secondaryDisplayProperties length = %d, want 2", len(secDisplay))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":                         "2-123456",
			"name":                       "cars",
			"labels":                     map[string]string{"singular": "Car", "plural": "Cars"},
			"requiredProperties":         []string{"vin"},
			"secondaryDisplayProperties": []string{"make", "model"},
		})
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	def, err := schemas.Update(context.Background(), "2-123456", &crm.SchemaUpdateInput{
		SecondaryDisplayProperties: []string{"make", "model"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def.ID != "2-123456" {
		t.Errorf("ID = %q, want 2-123456", def.ID)
	}
	if len(def.SecondaryDisplayProperties) != 2 {
		t.Fatalf("len(SecondaryDisplayProperties) = %d, want 2", len(def.SecondaryDisplayProperties))
	}
}

func TestSchemasService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if want := "/crm/v3/schemas/2-123456"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		if r.URL.Query().Get("archived") != "" {
			t.Errorf("archived query param should not be set, got %q", r.URL.Query().Get("archived"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	err := schemas.Delete(context.Background(), "2-123456", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchemasService_Delete_HardDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Query().Get("archived") != "true" {
			t.Errorf("archived = %q, want true", r.URL.Query().Get("archived"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	err := schemas.Delete(context.Background(), "2-123456", &crm.SchemaDeleteOptions{Archived: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchemasService_CreateAssociation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if want := "/crm/v3/schemas/2-123456/associations"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}

		body, _ := io.ReadAll(r.Body)
		var input map[string]any
		json.Unmarshal(body, &input)
		if input["fromObjectTypeId"] != "2-123456" {
			t.Errorf("fromObjectTypeId = %v, want 2-123456", input["fromObjectTypeId"])
		}
		if input["toObjectTypeId"] != "ticket" {
			t.Errorf("toObjectTypeId = %v, want ticket", input["toObjectTypeId"])
		}
		if input["name"] != "car_to_ticket" {
			t.Errorf("name = %v, want car_to_ticket", input["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":               "105",
			"fromObjectTypeId": "2-123456",
			"toObjectTypeId":   "ticket",
			"name":             "car_to_ticket",
		})
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	assoc, err := schemas.CreateAssociation(context.Background(), "2-123456", &crm.SchemaAssociationCreateInput{
		FromObjectTypeID: "2-123456",
		ToObjectTypeID:   "ticket",
		Name:             "car_to_ticket",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assoc.ID != "105" {
		t.Errorf("ID = %q, want 105", assoc.ID)
	}
	if assoc.Name != "car_to_ticket" {
		t.Errorf("Name = %q, want car_to_ticket", assoc.Name)
	}
}

func TestSchemasService_DeleteAssociation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if want := "/crm/v3/schemas/2-123456/associations/105"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	schemas := crm.NewService(newTestClient(t, ts)).Schemas()
	err := schemas.DeleteAssociation(context.Background(), "2-123456", "105")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

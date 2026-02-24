package crm_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
	"github.com/josiah-hester/hubspot-go-sdk/hubspot/crm"
)

// newTestClient creates a hubspot.Client pointed at the given test server.
func newTestClient(t *testing.T, ts *httptest.Server) *hubspot.Client {
	t.Helper()
	return hubspot.NewClient(
		hubspot.PrivateAppToken("test-token"),
		hubspot.WithBaseURL(ts.URL),
	)
}

// contactJSON returns a minimal valid CRM object JSON for use in test servers.
func contactJSON(id string, props map[string]string) map[string]any {
	return map[string]any{
		"id":         id,
		"properties": props,
		"createdAt":  "2024-01-15T10:30:00.000Z",
		"updatedAt":  "2024-06-20T14:00:00.000Z",
		"archived":   false,
	}
}

// ---- Get ----

func TestObjectsService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/123" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/123", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contactJSON("123", map[string]string{
			"email": "alice@example.com", "firstname": "Alice", "lastname": "Smith",
		}))
	}))
	defer ts.Close()

	client := newTestClient(t, ts)
	contact, err := crm.NewService(client).Contacts().Get(context.Background(), "123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contact.ID != "123" {
		t.Errorf("ID = %q, want 123", contact.ID)
	}
	if contact.Properties["email"] != "alice@example.com" {
		t.Errorf("email = %q, want alice@example.com", contact.Properties["email"])
	}
}

func TestObjectsService_Get_WithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("properties"); got != "email,firstname" {
			t.Errorf("properties = %q, want email,firstname", got)
		}
		if got := q.Get("associations"); got != "companies" {
			t.Errorf("associations = %q, want companies", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contactJSON("123", map[string]string{"email": "test@test.com"}))
	}))
	defer ts.Close()

	client := newTestClient(t, ts)
	_, err := crm.NewService(client).Contacts().Get(context.Background(), "123", &crm.GetOptions{
		Properties:   []string{"email", "firstname"},
		Associations: []string{"companies"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestObjectsService_Get_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"status": "error", "message": "Object not found",
			"correlationId": "test-corr-id", "category": "OBJECT_NOT_FOUND",
		})
	}))
	defer ts.Close()

	_, err := crm.NewService(newTestClient(t, ts)).Contacts().Get(context.Background(), "999", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !hubspot.IsNotFound(err) {
		t.Errorf("IsNotFound = false, want true; err = %v", err)
	}
}

// ---- List ----

func TestObjectsService_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Errorf("limit = %q, want 2", got)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"results": []any{
				contactJSON("1", map[string]string{"email": "a@b.com"}),
				contactJSON("2", map[string]string{"email": "c@d.com"}),
			},
			"paging": map[string]any{
				"next": map[string]string{"after": "cursor-abc"},
			},
		})
	}))
	defer ts.Close()

	resp, err := crm.NewService(newTestClient(t, ts)).Contacts().List(context.Background(), &crm.ListOptions{
		Limit: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(resp.Results))
	}
	if resp.Paging == nil || resp.Paging.Next == nil {
		t.Fatal("expected paging cursor")
	}
	if resp.Paging.Next.After != "cursor-abc" {
		t.Errorf("After = %q, want cursor-abc", resp.Paging.Next.After)
	}
}

func TestObjectsService_List_WithAfter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("after"); got != "cursor-xyz" {
			t.Errorf("after = %q, want cursor-xyz", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
	defer ts.Close()

	_, err := crm.NewService(newTestClient(t, ts)).Contacts().List(context.Background(), &crm.ListOptions{
		After: "cursor-xyz",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---- Create ----

func TestObjectsService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.CreateInput
		json.Unmarshal(body, &input)

		if input.Properties["email"] != "new@example.com" {
			t.Errorf("email = %q, want new@example.com", input.Properties["email"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contactJSON("789", input.Properties))
	}))
	defer ts.Close()

	obj, err := crm.NewService(newTestClient(t, ts)).Contacts().Create(context.Background(), &crm.CreateInput{
		Properties: map[string]string{"email": "new@example.com", "firstname": "New"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.ID != "789" {
		t.Errorf("ID = %q, want 789", obj.ID)
	}
}

func TestObjectsService_Create_WithAssociations(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var input crm.CreateInput
		json.Unmarshal(body, &input)

		if len(input.Associations) != 1 {
			t.Fatalf("len(Associations) = %d, want 1", len(input.Associations))
		}
		if input.Associations[0].To.ID != "company-1" {
			t.Errorf("association target = %q, want company-1", input.Associations[0].To.ID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contactJSON("789", input.Properties))
	}))
	defer ts.Close()

	_, err := crm.NewService(newTestClient(t, ts)).Contacts().Create(context.Background(), &crm.CreateInput{
		Properties: map[string]string{"email": "a@b.com"},
		Associations: []crm.CreateAssociation{{
			To: crm.CreateAssociationTarget{ID: "company-1"},
			Types: []crm.AssociationType{{
				AssociationCategory: "HUBSPOT_DEFINED",
				AssociationTypeID:   1,
			}},
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---- Update ----

func TestObjectsService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/123" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/123", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.UpdateInput
		json.Unmarshal(body, &input)

		if input.Properties["lastname"] != "Updated" {
			t.Errorf("lastname = %q, want Updated", input.Properties["lastname"])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contactJSON("123", map[string]string{
			"email": "a@b.com", "lastname": "Updated",
		}))
	}))
	defer ts.Close()

	obj, err := crm.NewService(newTestClient(t, ts)).Contacts().Update(context.Background(), "123", &crm.UpdateInput{
		Properties: map[string]string{"lastname": "Updated"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Properties["lastname"] != "Updated" {
		t.Errorf("lastname = %q, want Updated", obj.Properties["lastname"])
	}
}

// ---- Archive ----

func TestObjectsService_Archive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/123" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/123", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	err := crm.NewService(newTestClient(t, ts)).Contacts().Archive(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---- Merge ----

func TestObjectsService_Merge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/merge" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/merge", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.MergeInput
		json.Unmarshal(body, &input)
		if input.PrimaryObjectID != "1" || input.ObjectIDToMerge != "2" {
			t.Errorf("merge input = %+v, want primary=1 merge=2", input)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contactJSON("1", map[string]string{"email": "merged@test.com"}))
	}))
	defer ts.Close()

	obj, err := crm.NewService(newTestClient(t, ts)).Contacts().Merge(context.Background(), &crm.MergeInput{
		PrimaryObjectID: "1",
		ObjectIDToMerge: "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.ID != "1" {
		t.Errorf("ID = %q, want 1", obj.ID)
	}
}

// ---- GDPRDelete ----

func TestObjectsService_GDPRDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/gdpr-delete" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/gdpr-delete", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var input crm.GDPRDeleteInput
		json.Unmarshal(body, &input)
		if input.ObjectID != "123" {
			t.Errorf("ObjectID = %q, want 123", input.ObjectID)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	err := crm.NewService(newTestClient(t, ts)).Contacts().GDPRDelete(context.Background(), &crm.GDPRDeleteInput{
		ObjectID: "123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---- UnmarshalProperties ----

func TestObjectsService_UnmarshalProperties(t *testing.T) {
	obj := &crm.Object{
		ID: "123",
		Properties: map[string]string{
			"email": "bob@example.com", "firstname": "Bob", "lastname": "Jones",
		},
	}

	type Contact struct {
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	var c Contact
	if err := obj.UnmarshalProperties(&c); err != nil {
		t.Fatalf("UnmarshalProperties failed: %v", err)
	}
	if c.Email != "bob@example.com" {
		t.Errorf("Email = %q, want bob@example.com", c.Email)
	}
	if c.FirstName != "Bob" {
		t.Errorf("FirstName = %q, want Bob", c.FirstName)
	}
}

// ---- Object type routing ----

func TestObjectsService_DifferentObjectTypes(t *testing.T) {
	tests := []struct {
		name     string
		service  func(*crm.Service) *crm.ObjectsService
		wantPath string
	}{
		{"contacts", (*crm.Service).Contacts, "/crm/v3/objects/contacts/1"},
		{"companies", (*crm.Service).Companies, "/crm/v3/objects/companies/1"},
		{"deals", (*crm.Service).Deals, "/crm/v3/objects/deals/1"},
		{"tickets", (*crm.Service).Tickets, "/crm/v3/objects/tickets/1"},
		{"products", (*crm.Service).Products, "/crm/v3/objects/products/1"},
		{"line_items", (*crm.Service).LineItems, "/crm/v3/objects/line_items/1"},
		{"quotes", (*crm.Service).Quotes, "/crm/v3/objects/quotes/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(contactJSON("1", map[string]string{}))
			}))
			defer ts.Close()

			svc := tt.service(crm.NewService(newTestClient(t, ts)))
			_, err := svc.Get(context.Background(), "1", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotPath != tt.wantPath {
				t.Errorf("path = %q, want %q", gotPath, tt.wantPath)
			}
		})
	}
}

func TestObjectsService_CustomObject(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(contactJSON("1", map[string]string{}))
	}))
	defer ts.Close()

	_, err := crm.NewService(newTestClient(t, ts)).Object("2-12345").Get(context.Background(), "1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "/crm/v3/objects/2-12345/1"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
}

func TestObjectsService_ObjectType(t *testing.T) {
	client := hubspot.NewClient(hubspot.PrivateAppToken("token"))
	svc := crm.NewService(client)

	if got := svc.Contacts().ObjectType(); got != "contacts" {
		t.Errorf("ObjectType = %q, want contacts", got)
	}
	if got := svc.Object("pets").ObjectType(); got != "pets" {
		t.Errorf("ObjectType = %q, want pets", got)
	}
}

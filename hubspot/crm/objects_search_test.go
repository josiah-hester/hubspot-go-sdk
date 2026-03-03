package crm_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot/crm"
)

func TestObjectsService_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/crm/v3/objects/contacts/search" {
			t.Errorf("path = %s, want /crm/v3/objects/contacts/search", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req crm.SearchRequest
		json.Unmarshal(body, &req)

		if len(req.FilterGroups) != 1 {
			t.Fatalf("len(FilterGroups) = %d, want 1", len(req.FilterGroups))
		}
		if len(req.FilterGroups[0].Filters) != 1 {
			t.Fatalf("len(Filters) = %d, want 1", len(req.FilterGroups[0].Filters))
		}

		f := req.FilterGroups[0].Filters[0]
		if f.PropertyName != "email" || f.Operator != crm.OpContainsToken || f.Value != "example.com" {
			t.Errorf("filter = %+v, want email CONTAINS_TOKEN example.com", f)
		}
		if req.Limit != 10 {
			t.Errorf("limit = %d, want 10", req.Limit)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total": 42,
			"results": []any{
				contactJSON("1", map[string]string{"email": "a@example.com"}),
			},
		})
	}))
	defer ts.Close()

	resp, err := crm.NewService(newTestClient(t, ts)).Contacts().Search(context.Background(), &crm.SearchRequest{
		FilterGroups: []crm.FilterGroup{{
			Filters: []crm.Filter{{
				PropertyName: "email",
				Operator:     crm.OpContainsToken,
				Value:        "example.com",
			}},
		}},
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 42 {
		t.Errorf("Total = %d, want 42", resp.Total)
	}
	if len(resp.Results) != 1 {
		t.Errorf("len(Results) = %d, want 1", len(resp.Results))
	}
}

func TestObjectsService_Search_WithSort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req crm.SearchRequest
		json.Unmarshal(body, &req)

		if len(req.Sorts) != 1 {
			t.Fatalf("len(Sorts) = %d, want 1", len(req.Sorts))
		}
		if req.Sorts[0] != "createdate" {
			t.Errorf("sort = %+v, want createdate DESCENDING", req.Sorts[0])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"total": 0, "results": []any{}})
	}))
	defer ts.Close()

	_, err := crm.NewService(newTestClient(t, ts)).Contacts().Search(context.Background(), &crm.SearchRequest{
		Sorts: []string{"createdate"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---- SearchBuilder ----

func TestSearchBuilder_SingleGroup(t *testing.T) {
	req := crm.NewSearch().
		Where("email", crm.OpContainsToken, "example.com").
		Where("firstname", crm.OpEQ, "Alice").
		Build()

	if len(req.FilterGroups) != 1 {
		t.Fatalf("len(FilterGroups) = %d, want 1", len(req.FilterGroups))
	}
	if len(req.FilterGroups[0].Filters) != 2 {
		t.Fatalf("len(Filters) = %d, want 2", len(req.FilterGroups[0].Filters))
	}
	if req.FilterGroups[0].Filters[0].Operator != crm.OpContainsToken {
		t.Errorf("filter[0].Operator = %q, want CONTAINS_TOKEN", req.FilterGroups[0].Filters[0].Operator)
	}
	if req.FilterGroups[0].Filters[1].PropertyName != "firstname" {
		t.Errorf("filter[1].PropertyName = %q, want firstname", req.FilterGroups[0].Filters[1].PropertyName)
	}
}

func TestSearchBuilder_MultipleGroups(t *testing.T) {
	req := crm.NewSearch().
		Where("email", crm.OpContainsToken, "example.com").
		Or().
		Where("email", crm.OpContainsToken, "test.com").
		Build()

	if len(req.FilterGroups) != 2 {
		t.Fatalf("len(FilterGroups) = %d, want 2", len(req.FilterGroups))
	}
	if req.FilterGroups[0].Filters[0].Value != "example.com" {
		t.Errorf("group[0] value = %q, want example.com", req.FilterGroups[0].Filters[0].Value)
	}
	if req.FilterGroups[1].Filters[0].Value != "test.com" {
		t.Errorf("group[1] value = %q, want test.com", req.FilterGroups[1].Filters[0].Value)
	}
}

func TestSearchBuilder_WhereIn(t *testing.T) {
	req := crm.NewSearch().
		WhereIn("hs_pipeline", "pipeline-1", "pipeline-2").
		Build()

	f := req.FilterGroups[0].Filters[0]
	if f.Operator != crm.OpIN {
		t.Errorf("Operator = %q, want IN", f.Operator)
	}
	if !reflect.DeepEqual(f.Values, []string{"pipeline-1", "pipeline-2"}) {
		t.Errorf("Values = %v, want [pipeline-1 pipeline-2]", f.Values)
	}
}

func TestSearchBuilder_WhereNotIn(t *testing.T) {
	req := crm.NewSearch().
		WhereNotIn("lifecyclestage", "subscriber").
		Build()

	f := req.FilterGroups[0].Filters[0]
	if f.Operator != crm.OpNotIN {
		t.Errorf("Operator = %q, want NOT_IN", f.Operator)
	}
}

func TestSearchBuilder_WhereBetween(t *testing.T) {
	req := crm.NewSearch().
		WhereBetween("amount", "1000", "5000").
		Build()

	f := req.FilterGroups[0].Filters[0]
	if f.Operator != crm.OpBetween {
		t.Errorf("Operator = %q, want BETWEEN", f.Operator)
	}
	if f.Value != "1000" {
		t.Errorf("Value = %q, want 1000", f.Value)
	}
	if f.HighValue != "5000" {
		t.Errorf("HighValue = %q, want 5000", f.HighValue)
	}
}

func TestSearchBuilder_WhereHasProperty(t *testing.T) {
	req := crm.NewSearch().
		WhereHasProperty("email").
		Build()

	f := req.FilterGroups[0].Filters[0]
	if f.Operator != crm.OpHasProperty {
		t.Errorf("Operator = %q, want HAS_PROPERTY", f.Operator)
	}
	if f.PropertyName != "email" {
		t.Errorf("PropertyName = %q, want email", f.PropertyName)
	}
}

func TestSearchBuilder_WhereNotHasProperty(t *testing.T) {
	req := crm.NewSearch().
		WhereNotHasProperty("phone").
		Build()

	f := req.FilterGroups[0].Filters[0]
	if f.Operator != crm.OpNotHasProperty {
		t.Errorf("Operator = %q, want NOT_HAS_PROPERTY", f.Operator)
	}
}

func TestSearchBuilder_SortSelectLimitAfter(t *testing.T) {
	req := crm.NewSearch().
		Where("email", crm.OpHasProperty, "").
		SortBy("createdate").
		SortBy("email").
		Select("email", "firstname").
		Limit(50).
		After("cursor-123").
		Build()

	if len(req.Sorts) != 2 {
		t.Fatalf("len(Sorts) = %d, want 2", len(req.Sorts))
	}
	if req.Sorts[0] != "createdate" {
		t.Errorf("sort[0] = %+v, want createdate DESCENDING", req.Sorts[0])
	}
	if !reflect.DeepEqual(req.Properties, []string{"email", "firstname"}) {
		t.Errorf("Properties = %v, want [email firstname]", req.Properties)
	}
	if req.Limit != 50 {
		t.Errorf("Limit = %d, want 50", req.Limit)
	}
	if req.After != "cursor-123" {
		t.Errorf("After = %q, want cursor-123", req.After)
	}
}

func TestSearchBuilder_EmptyBuild(t *testing.T) {
	req := crm.NewSearch().Build()

	if len(req.FilterGroups) != 0 {
		t.Errorf("len(FilterGroups) = %d, want 0", len(req.FilterGroups))
	}
	if req.Limit != 0 {
		t.Errorf("Limit = %d, want 0", req.Limit)
	}
}

func TestSearchBuilder_OrWithoutPriorFilters(t *testing.T) {
	// Or() with no prior filters should be a no-op.
	req := crm.NewSearch().
		Or().
		Where("email", crm.OpEQ, "a@b.com").
		Build()

	if len(req.FilterGroups) != 1 {
		t.Fatalf("len(FilterGroups) = %d, want 1", len(req.FilterGroups))
	}
}

func TestSearchBuilder_ComplexQuery(t *testing.T) {
	// (email contains example.com AND firstname=Alice) OR (email contains test.com)
	req := crm.NewSearch().
		Where("email", crm.OpContainsToken, "example.com").
		Where("firstname", crm.OpEQ, "Alice").
		Or().
		Where("email", crm.OpContainsToken, "test.com").
		SortBy("createdate").
		Select("email", "firstname", "lastname").
		Limit(20).
		Build()

	if len(req.FilterGroups) != 2 {
		t.Fatalf("len(FilterGroups) = %d, want 2", len(req.FilterGroups))
	}
	if len(req.FilterGroups[0].Filters) != 2 {
		t.Errorf("group[0] filters = %d, want 2", len(req.FilterGroups[0].Filters))
	}
	if len(req.FilterGroups[1].Filters) != 1 {
		t.Errorf("group[1] filters = %d, want 1", len(req.FilterGroups[1].Filters))
	}
	if len(req.Properties) != 3 {
		t.Errorf("len(Properties) = %d, want 3", len(req.Properties))
	}
}

// ---- Operator constants ----

func TestOperatorConstants(t *testing.T) {
	// Verify constants match HubSpot's expected string values.
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"EQ", crm.OpEQ, "EQ"},
		{"NEQ", crm.OpNEQ, "NEQ"},
		{"LT", crm.OpLT, "LT"},
		{"LTE", crm.OpLTE, "LTE"},
		{"GT", crm.OpGT, "GT"},
		{"GTE", crm.OpGTE, "GTE"},
		{"BETWEEN", crm.OpBetween, "BETWEEN"},
		{"IN", crm.OpIN, "IN"},
		{"NOT_IN", crm.OpNotIN, "NOT_IN"},
		{"HAS_PROPERTY", crm.OpHasProperty, "HAS_PROPERTY"},
		{"NOT_HAS_PROPERTY", crm.OpNotHasProperty, "NOT_HAS_PROPERTY"},
		{"CONTAINS_TOKEN", crm.OpContainsToken, "CONTAINS_TOKEN"},
		{"NOT_CONTAINS_TOKEN", crm.OpNotContainsToken, "NOT_CONTAINS_TOKEN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

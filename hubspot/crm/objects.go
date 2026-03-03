package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// ObjectsService provides CRUD, batch, and search operations for a single
// CRM object type. All standard and custom object types share the same
// API surface.
//
// Obtain an ObjectsService from a [Service]:
//
//	crm := crm.NewService(client)
//	contacts := crm.Contacts()
//	companies := crm.Companies()
//	custom := crm.Object("2-12345")
type ObjectsService struct {
	r          hubspot.Requester
	objectType string
}

func newObjectsService(r hubspot.Requester, objectType string) *ObjectsService {
	return &ObjectsService{
		r:          r,
		objectType: objectType,
	}
}

// ObjectType returns the CRM object type this service operates on.
func (s *ObjectsService) ObjectType() string {
	return s.objectType
}

// Get retrieves a single CRM object by its ID.
//
//	contact, err := contacts.Get(ctx, "123", &crm.GetOptions{
//	    Properties: []string{"email", "firstname", "lastname"},
//	})
func (s *ObjectsService) Get(ctx context.Context, id string, opts *GetOptions) (*Object, error) {
	obj := &Object{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), id),
		Query:  opts.toQuery(),
	}, obj); err != nil {
		return nil, fmt.Errorf("hubspot: get object: %w", err)
	}

	return obj, nil
}

// List returns a page of CRM objects. Use [ListOptions.After] with the
// cursor from [ListResponse.Paging] to paginate through all results.
//
//	resp, err := contacts.List(ctx, &crm.ListOptions{
//	    Limit:      10,
//	    Properties: []string{"email", "firstname"},
//	})
func (s *ObjectsService) List(ctx context.Context, opts *ListOptions) (*ListResponse, error) {
	resp := &ListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   s.basePath(),
		Query:  opts.toQuery(),
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: list objects: %w", err)
	}

	return resp, nil
}

// Create creates a new CRM object with the given properties.
//
//	contact, err := contacts.Create(ctx, &crm.CreateInput{
//	    Properties: map[string]string{
//	        "email":     "alice@example.com",
//	        "firstname": "Alice",
//	    },
//	})
func (s *ObjectsService) Create(ctx context.Context, input *CreateInput) (*Object, error) {
	obj := &Object{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   s.basePath(),
		Body:   input,
	}, obj); err != nil {
		return nil, fmt.Errorf("hubspot: create object: %w", err)
	}

	return obj, nil
}

// Update modifies properties on an existing CRM object. Only properties
// included in the input are changed; all others remain untouched.
//
//	contact, err := contacts.Update(ctx, "123", &crm.UpdateInput{
//	    Properties: map[string]string{"lastname": "Smith"},
//	})
func (s *ObjectsService) Update(ctx context.Context, id string, idProperty string, input *UpdateInput) (*Object, error) {
	obj := &Object{}
	req := &hubspot.Request{
		Method: "PATCH",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), id),
		Body:   input,
	}
	if idProperty != "" {
		req.Query.Set("idProperty", idProperty)
	}
	if err := s.r.Do(ctx, req, obj); err != nil {
		return nil, fmt.Errorf("hubspot: update object: %w", err)
	}

	return obj, nil
}

// Archive soft-deletes a CRM object. The object can be restored from the
// HubSpot recycling bin within 90 days.
//
//	err := contacts.Archive(ctx, "123")
func (s *ObjectsService) Archive(ctx context.Context, id string) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), id),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: archive object: %w", err)
	}

	return nil
}

// Merge combines two CRM objects into one. The primary object is kept
// and the secondary object is archived. Property values from the secondary
// object fill in any blank properties on the primary.
//
//	err := contacts.Merge(ctx, &crm.MergeInput{
//	    PrimaryObjectID: "123",
//	    ObjectIDToMerge: "456",
//	})
func (s *ObjectsService) Merge(ctx context.Context, input *MergeInput) (*Object, error) {
	obj := &Object{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/merge", s.basePath()),
		Body:   input,
	}, obj); err != nil {
		return nil, fmt.Errorf("hubspot: merge objects: %w", err)
	}

	return obj, nil
}

// GDPRDelete permanently deletes a CRM object in compliance with GDPR.
// Unlike [Archive], this cannot be undone.
//
//	err := contacts.GDPRDelete(ctx, &crm.GDPRDeleteInput{
//	    ObjectID: "123",
//	})
func (s *ObjectsService) GDPRDelete(ctx context.Context, input *GDPRDeleteInput) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/crm/v3/objects/%s/gdpr-delete", s.objectType),
		Body:   input,
	}, nil); err != nil {
		return fmt.Errorf("hubspot: gdpr delete object: %w", err)
	}

	return nil
}

// Batch returns a [BatchService] for performing bulk operations on this
// object type.
func (s *ObjectsService) Batch() *BatchService {
	return newBatchService(s.r, s.objectType)
}

// Search returns objects matching the given search request. See
// [SearchRequest] for filter and sort options, or use [NewSearch] for
// a fluent builder.
//
//	resp, err := contacts.Search(ctx, &crm.SearchRequest{
//	    FilterGroups: []crm.FilterGroup{{
//	        Filters: []crm.Filter{{
//	            PropertyName: "email",
//	            Operator:     crm.OpContainsToken,
//	            Value:        "example.com",
//	        }},
//	    }},
//	    Limit: 10,
//	})
func (s *ObjectsService) Search(ctx context.Context, req *SearchRequest) (*ListResponse, error) {
	resp := &ListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/search", s.basePath()),
		Body:   req,
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: search objects: %w", err)
	}

	return resp, nil
}

func (s *ObjectsService) basePath() string {
	return fmt.Sprintf("/crm/v3/objects/%s", s.objectType)
}

package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// AssociationsService provides CRUD operations for associations between two
// CRM object types using the v4 Associations API.
//
// Obtain an AssociationsService from a [Service]:
//
//	assoc := crm.NewService(client).Associations("contacts", "companies")
//	result, err := assoc.Create(ctx, "contact-1", "company-1", types)
type AssociationsService struct {
	r              hubspot.Requester
	fromObjectType string
	toObjectType   string
}

func newAssociationsService(r hubspot.Requester, fromObjectType, toObjectType string) *AssociationsService {
	return &AssociationsService{
		r:              r,
		fromObjectType: fromObjectType,
		toObjectType:   toObjectType,
	}
}

// Create creates an association between two objects. The body is a bare
// array of association types (not wrapped in {"inputs": ...}).
//
//	result, err := assoc.Create(ctx, "contact-1", "company-1", []crm.AssociationTypeInput{
//	    {AssociationCategory: crm.AssociationCategoryHubSpotDefined, AssociationTypeID: 1},
//	})
func (s *AssociationsService) Create(ctx context.Context, fromObjectID, toObjectID string, types []AssociationTypeInput) (*AssociationResult, error) {
	result := &AssociationResult{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "PUT",
		Path:   fmt.Sprintf("%s/%s/associations/%s/%s", s.basePath(), fromObjectID, s.toObjectType, toObjectID),
		Body:   types,
	}, result); err != nil {
		return nil, fmt.Errorf("hubspot: create association: %w", err)
	}
	return result, nil
}

// List returns all associations from the given object to the configured
// target object type.
//
//	resp, err := assoc.List(ctx, "contact-1")
func (s *AssociationsService) List(ctx context.Context, fromObjectID string) (*AssociationListResponse, error) {
	resp := &AssociationListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s/associations/%s", s.basePath(), fromObjectID, s.toObjectType),
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: list associations: %w", err)
	}
	return resp, nil
}

// Archive removes all associations between two objects.
//
//	err := assoc.Archive(ctx, "contact-1", "company-1")
func (s *AssociationsService) Archive(ctx context.Context, fromObjectID, toObjectID string) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s/associations/%s/%s", s.basePath(), fromObjectID, s.toObjectType, toObjectID),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: archive association: %w", err)
	}
	return nil
}

// Batch returns an [AssociationsBatchService] for bulk association operations.
func (s *AssociationsService) Batch() *AssociationsBatchService {
	return newAssociationsBatchService(s.r, s.fromObjectType, s.toObjectType)
}

// Schema returns an [AssociationsSchemaService] for managing custom
// association labels and definitions.
func (s *AssociationsService) Schema() *AssociationsSchemaService {
	return newAssociationsSchemaService(s.r, s.fromObjectType, s.toObjectType)
}

func (s *AssociationsService) basePath() string {
	return fmt.Sprintf("/crm/v4/objects/%s", s.fromObjectType)
}

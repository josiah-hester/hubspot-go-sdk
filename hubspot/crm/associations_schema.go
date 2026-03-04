package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// AssociationsSchemaService manages custom association labels and definitions
// using the v4 Associations API.
//
// Obtain an AssociationsSchemaService from an [AssociationsService]:
//
//	schema := crm.NewService(client).Associations("contacts", "companies").Schema()
type AssociationsSchemaService struct {
	r              hubspot.Requester
	fromObjectType string
	toObjectType   string
}

func newAssociationsSchemaService(r hubspot.Requester, fromObjectType, toObjectType string) *AssociationsSchemaService {
	return &AssociationsSchemaService{
		r:              r,
		fromObjectType: fromObjectType,
		toObjectType:   toObjectType,
	}
}

// List returns all association label definitions between the configured
// object types.
//
//	resp, err := schema.List(ctx)
func (s *AssociationsSchemaService) List(ctx context.Context) (*AssociationLabelListResponse, error) {
	resp := &AssociationLabelListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   s.basePath(),
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: list association labels: %w", err)
	}
	return resp, nil
}

// Create creates a new custom association label.
//
//	resp, err := schema.Create(ctx, &crm.CreateAssociationLabelInput{
//	    Label: "Partner",
//	    Name:  "partner",
//	})
func (s *AssociationsSchemaService) Create(ctx context.Context, input *CreateAssociationLabelInput) (*CreateAssociationLabelResponse, error) {
	resp := &CreateAssociationLabelResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   s.basePath(),
		Body:   input,
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: create association label: %w", err)
	}
	return resp, nil
}

// Update updates an existing custom association label.
//
//	err := schema.Update(ctx, &crm.UpdateAssociationLabelInput{
//	    AssociationTypeID: 42,
//	    Label:             "New Label",
//	})
func (s *AssociationsSchemaService) Update(ctx context.Context, input *UpdateAssociationLabelInput) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "PUT",
		Path:   s.basePath(),
		Body:   input,
	}, nil); err != nil {
		return fmt.Errorf("hubspot: update association label: %w", err)
	}
	return nil
}

// Delete removes a custom association label by its type ID.
//
//	err := schema.Delete(ctx, 42)
func (s *AssociationsSchemaService) Delete(ctx context.Context, associationTypeID int) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%d", s.basePath(), associationTypeID),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: delete association label: %w", err)
	}
	return nil
}

func (s *AssociationsSchemaService) basePath() string {
	return fmt.Sprintf("/crm/v4/associations/%s/%s/labels", s.fromObjectType, s.toObjectType)
}

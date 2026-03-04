package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

const schemasBasePath = "/crm/v3/schemas"

// SchemasService provides CRUD operations for CRM custom object schemas.
//
// Obtain a SchemasService from a [Service]:
//
//	schemas := crm.NewService(client).Schemas()
//	resp, err := schemas.List(ctx, nil)
type SchemasService struct {
	r hubspot.Requester
}

func newSchemasService(r hubspot.Requester) *SchemasService {
	return &SchemasService{r: r}
}

// List returns all custom object schemas.
func (s *SchemasService) List(ctx context.Context, opts *SchemaListOptions) (*SchemaListResponse, error) {
	resp := &SchemaListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   schemasBasePath,
		Query:  opts.toQuery(),
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: list schemas: %w", err)
	}
	return resp, nil
}

// Get retrieves a single schema by object type name or ID.
func (s *SchemasService) Get(ctx context.Context, objectType string, opts *SchemaGetOptions) (*Schema, error) {
	schema := &Schema{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", schemasBasePath, objectType),
		Query:  opts.toQuery(),
	}, schema); err != nil {
		return nil, fmt.Errorf("hubspot: get schema: %w", err)
	}
	return schema, nil
}

// Create creates a new custom object schema.
func (s *SchemasService) Create(ctx context.Context, input *SchemaCreateInput) (*Schema, error) {
	schema := &Schema{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   schemasBasePath,
		Body:   input,
	}, schema); err != nil {
		return nil, fmt.Errorf("hubspot: create schema: %w", err)
	}
	return schema, nil
}

// Update modifies an existing schema's configuration. Only fields included
// in the input are changed.
func (s *SchemasService) Update(ctx context.Context, objectType string, input *SchemaUpdateInput) (*SchemaDefinition, error) {
	def := &SchemaDefinition{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "PATCH",
		Path:   fmt.Sprintf("%s/%s", schemasBasePath, objectType),
		Body:   input,
	}, def); err != nil {
		return nil, fmt.Errorf("hubspot: update schema: %w", err)
	}
	return def, nil
}

// Delete removes a schema. Pass nil opts for a soft delete, or set
// Archived to true for a permanent hard delete.
func (s *SchemasService) Delete(ctx context.Context, objectType string, opts *SchemaDeleteOptions) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", schemasBasePath, objectType),
		Query:  opts.toQuery(),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: delete schema: %w", err)
	}
	return nil
}

// CreateAssociation adds an association definition to a schema.
func (s *SchemasService) CreateAssociation(ctx context.Context, objectType string, input *SchemaAssociationCreateInput) (*SchemaAssociationDef, error) {
	assoc := &SchemaAssociationDef{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/%s/associations", schemasBasePath, objectType),
		Body:   input,
	}, assoc); err != nil {
		return nil, fmt.Errorf("hubspot: create schema association: %w", err)
	}
	return assoc, nil
}

// DeleteAssociation removes an association definition from a schema.
func (s *SchemasService) DeleteAssociation(ctx context.Context, objectType string, associationIdentifier string) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s/associations/%s", schemasBasePath, objectType, associationIdentifier),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: delete schema association: %w", err)
	}
	return nil
}

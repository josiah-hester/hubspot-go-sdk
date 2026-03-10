package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

const propertiesBasePath = "/crm/v3/properties"

// PropertiesService provides CRUD operations for CRM property definitions
// using the v3 Properties API.
//
// Obtain a PropertiesService from a [Service]:
//
//	props := crm.NewService(client).Properties("contacts")
//	resp, err := props.List(ctx, nil)
type PropertiesService struct {
	r          hubspot.Requester
	objectType string
}

func newPropertiesService(r hubspot.Requester, objectType string) *PropertiesService {
	return &PropertiesService{r: r, objectType: objectType}
}

// List returns all property definitions for the object type.
func (s *PropertiesService) List(ctx context.Context, opts *PropertyListOptions) (*PropertyListResponse, error) {
	resp := &PropertyListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   s.basePath(),
		Query:  opts.toQuery(),
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: list properties: %w", err)
	}
	return resp, nil
}

// Get retrieves a single property definition by name.
func (s *PropertiesService) Get(ctx context.Context, propertyName string, opts *PropertyGetOptions) (*Property, error) {
	prop := &Property{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), propertyName),
		Query:  opts.toQuery(),
	}, prop); err != nil {
		return nil, fmt.Errorf("hubspot: get property: %w", err)
	}
	return prop, nil
}

// Create creates a new property definition.
func (s *PropertiesService) Create(ctx context.Context, input *PropertyCreateInput) (*Property, error) {
	prop := &Property{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   s.basePath(),
		Body:   input,
	}, prop); err != nil {
		return nil, fmt.Errorf("hubspot: create property: %w", err)
	}
	return prop, nil
}

// Update modifies an existing property definition.
func (s *PropertiesService) Update(ctx context.Context, propertyName string, input *PropertyUpdateInput) (*Property, error) {
	prop := &Property{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "PATCH",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), propertyName),
		Body:   input,
	}, prop); err != nil {
		return nil, fmt.Errorf("hubspot: update property: %w", err)
	}
	return prop, nil
}

// Archive removes a property definition.
func (s *PropertiesService) Archive(ctx context.Context, propertyName string) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), propertyName),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: archive property: %w", err)
	}
	return nil
}

// Groups returns a [PropertyGroupsService] for managing property groups.
func (s *PropertiesService) Groups() *PropertyGroupsService {
	return newPropertyGroupsService(s.r, s.objectType)
}

// Batch returns a [PropertyBatchService] for bulk property operations.
func (s *PropertiesService) Batch() *PropertyBatchService {
	return newPropertyBatchService(s.r, s.objectType)
}

func (s *PropertiesService) basePath() string {
	return fmt.Sprintf("%s/%s", propertiesBasePath, s.objectType)
}

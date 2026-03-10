package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// PropertyGroupsService provides CRUD operations for CRM property groups.
//
// Obtain a PropertyGroupsService from a [PropertiesService]:
//
//	groups := crm.NewService(client).Properties("contacts").Groups()
//	resp, err := groups.List(ctx, nil)
type PropertyGroupsService struct {
	r          hubspot.Requester
	objectType string
}

func newPropertyGroupsService(r hubspot.Requester, objectType string) *PropertyGroupsService {
	return &PropertyGroupsService{r: r, objectType: objectType}
}

// List returns all property groups for the object type.
func (s *PropertyGroupsService) List(ctx context.Context, opts *PropertyGroupListOptions) (*PropertyGroupListResponse, error) {
	resp := &PropertyGroupListResponse{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   s.basePath(),
		Query:  opts.toQuery(),
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: list property groups: %w", err)
	}
	return resp, nil
}

// Get retrieves a single property group by name.
func (s *PropertyGroupsService) Get(ctx context.Context, groupName string, opts *PropertyGroupGetOptions) (*PropertyGroup, error) {
	group := &PropertyGroup{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), groupName),
		Query:  opts.toQuery(),
	}, group); err != nil {
		return nil, fmt.Errorf("hubspot: get property group: %w", err)
	}
	return group, nil
}

// Create creates a new property group.
func (s *PropertyGroupsService) Create(ctx context.Context, input *PropertyGroupCreateInput) (*PropertyGroup, error) {
	group := &PropertyGroup{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   s.basePath(),
		Body:   input,
	}, group); err != nil {
		return nil, fmt.Errorf("hubspot: create property group: %w", err)
	}
	return group, nil
}

// Update modifies an existing property group.
func (s *PropertyGroupsService) Update(ctx context.Context, groupName string, input *PropertyGroupUpdateInput) (*PropertyGroup, error) {
	group := &PropertyGroup{}
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "PATCH",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), groupName),
		Body:   input,
	}, group); err != nil {
		return nil, fmt.Errorf("hubspot: update property group: %w", err)
	}
	return group, nil
}

// Archive removes a property group.
func (s *PropertyGroupsService) Archive(ctx context.Context, groupName string) error {
	if err := s.r.Do(ctx, &hubspot.Request{
		Method: "DELETE",
		Path:   fmt.Sprintf("%s/%s", s.basePath(), groupName),
	}, nil); err != nil {
		return fmt.Errorf("hubspot: archive property group: %w", err)
	}
	return nil
}

func (s *PropertyGroupsService) basePath() string {
	return fmt.Sprintf("%s/%s/groups", propertiesBasePath, s.objectType)
}

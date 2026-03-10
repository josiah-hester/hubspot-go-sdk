package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// PropertyBatchService provides bulk property operations using the
// v3 Properties API.
//
// Obtain a PropertyBatchService from a [PropertiesService]:
//
//	batch := crm.NewService(client).Properties("contacts").Batch()
//	resp, err := batch.Read(ctx, req, nil)
type PropertyBatchService struct {
	r          hubspot.Requester
	objectType string
}

func newPropertyBatchService(r hubspot.Requester, objectType string) *PropertyBatchService {
	return &PropertyBatchService{r: r, objectType: objectType}
}

// Create creates properties in bulk.
func (b *PropertyBatchService) Create(ctx context.Context, inputs []PropertyCreateInput) (*PropertyBatchResponse, error) {
	resp := &PropertyBatchResponse{}
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/create", b.basePath()),
		Body:   map[string]any{"inputs": inputs},
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: batch create properties: %w", err)
	}
	return resp, nil
}

// Read reads properties in bulk by name.
func (b *PropertyBatchService) Read(ctx context.Context, req *PropertyBatchReadRequest, opts *PropertyBatchReadOptions) (*PropertyBatchResponse, error) {
	resp := &PropertyBatchResponse{}
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/read", b.basePath()),
		Query:  opts.toQuery(),
		Body:   req,
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: batch read properties: %w", err)
	}
	return resp, nil
}

// Archive removes properties in bulk by name.
func (b *PropertyBatchService) Archive(ctx context.Context, inputs []PropertyBatchReadInput) error {
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/archive", b.basePath()),
		Body:   map[string]any{"inputs": inputs},
	}, nil); err != nil {
		return fmt.Errorf("hubspot: batch archive properties: %w", err)
	}
	return nil
}

func (b *PropertyBatchService) basePath() string {
	return fmt.Sprintf("%s/%s", propertiesBasePath, b.objectType)
}

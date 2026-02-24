package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// BatchService provides bulk create, read, update, upsert, and archive
// operations for a single CRM object type. Batch endpoints accept up to
// 100 items per request.
//
// Obtain a BatchService from an [ObjectsService]:
//
//	batch := crm.Contacts().Batch()
//	resp, err := batch.Create(ctx, []crm.BatchCreateInput{...})
type BatchService struct {
	r          hubspot.Requester
	objectType string
}

func newBatchService(r hubspot.Requester, objectType string) *BatchService {
	return &BatchService{
		r:          r,
		objectType: objectType,
	}
}

// Create creates up to 100 CRM objects in a single request.
//
//	resp, err := batch.Create(ctx, []crm.BatchCreateInput{
//	    {Properties: map[string]string{"email": "a@b.com"}},
//	    {Properties: map[string]string{"email": "c@d.com"}},
//	})
func (b *BatchService) Create(ctx context.Context, inputs []BatchCreateInput) (*BatchResponse, error) {
	resp := &BatchResponse{}
	err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/crm/v3/objects/%s/batch/create", b.objectType),
		Body:   map[string]any{"inputs": inputs},
	}, resp)
	if err != nil {
		return nil, fmt.Errorf("hubspot: create objects: %w", err)
	}

	return resp, nil
}

// Read retrieves up to 100 CRM objects by ID in a single request.
//
//	resp, err := batch.Read(ctx, &crm.BatchReadInput{
//	    Properties: []string{"email", "firstname"},
//	    Inputs:     []crm.BatchReadID{{ID: "1"}, {ID: "2"}},
//	})
func (b *BatchService) Read(ctx context.Context, input *BatchReadInput) (*BatchResponse, error) {
	resp := &BatchResponse{}
	err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/crm/v3/objects/%s/batch/read", b.objectType),
		Body:   input,
	}, resp)
	if err != nil {
		return nil, fmt.Errorf("hubspot: read objects: %w", err)
	}

	return resp, nil
}

// Update modifies up to 100 CRM objects in a single request.
//
//	resp, err := batch.Update(ctx, []crm.BatchUpdateInput{
//	    {ID: "1", Properties: map[string]string{"lastname": "Smith"}},
//	    {ID: "2", Properties: map[string]string{"lastname": "Jones"}},
//	})
func (b *BatchService) Update(ctx context.Context, inputs []BatchUpdateInput) (*BatchResponse, error) {
	resp := &BatchResponse{}
	err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/crm/v3/objects/%s/batch/update", b.objectType),
		Body:   map[string]any{"inputs": inputs},
	}, resp)
	if err != nil {
		return nil, fmt.Errorf("hubspot: update objects: %w", err)
	}

	return resp, nil
}

// Upsert creates or updates up to 100 CRM objects in a single request.
// Objects are matched by the IDProperty field on each input. If a match
// is found the object is updated; otherwise a new object is created.
//
//	resp, err := batch.Upsert(ctx, []crm.BatchUpsertInput{
//	    {
//	        IDProperty: "email",
//	        ID:         "a@b.com",
//	        Properties: map[string]string{"firstname": "Alice"},
//	    },
//	})
func (b *BatchService) Upsert(ctx context.Context, inputs []BatchUpsertInput) (*BatchResponse, error) {
	resp := &BatchResponse{}
	err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/crm/v3/objects/%s/batch/upsert", b.objectType),
		Body:   map[string]any{"inputs": inputs},
	}, resp)
	if err != nil {
		return nil, fmt.Errorf("hubspot: upsert objects: %w", err)
	}

	return resp, nil
}

// Archive soft-deletes up to 100 CRM objects in a single request.
//
//	err := batch.Archive(ctx, []crm.BatchArchiveInput{
//	    {ID: "1"}, {ID: "2"},
//	})
func (b *BatchService) Archive(ctx context.Context, inputs []BatchArchiveInput) error {
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/crm/v3/objects/%s/batch/archive", b.objectType),
		Body:   map[string]any{"inputs": inputs},
	}, nil); err != nil {
		return fmt.Errorf("hubspot: archive objects: %w", err)
	}
	return nil
}

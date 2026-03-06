package crm

import (
	"context"
	"fmt"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
)

// AssociationsBatchService provides bulk association operations using the
// v4 Associations API.
//
// Obtain an AssociationsBatchService from an [AssociationsService]:
//
//	batch := crm.NewService(client).Associations("contacts", "companies").Batch()
type AssociationsBatchService struct {
	r              hubspot.Requester
	fromObjectType string
	toObjectType   string
}

func newAssociationsBatchService(r hubspot.Requester, fromObjectType, toObjectType string) *AssociationsBatchService {
	return &AssociationsBatchService{
		r:              r,
		fromObjectType: fromObjectType,
		toObjectType:   toObjectType,
	}
}

// Create creates associations in bulk.
//
//	resp, err := batch.Create(ctx, []crm.AssociationBatchCreateInput{
//	    {
//	        From:  crm.AssociationObjectID{ID: "contact-1"},
//	        To:    crm.AssociationObjectID{ID: "company-1"},
//	        Types: []crm.AssociationTypeInput{{AssociationCategory: "HUBSPOT_DEFINED", AssociationTypeID: 1}},
//	    },
//	})
func (b *AssociationsBatchService) Create(ctx context.Context, inputs []AssociationBatchCreateInput) (*AssociationBatchCreateResponse, error) {
	resp := &AssociationBatchCreateResponse{}
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/create", b.basePath()),
		Body:   map[string]any{"inputs": inputs},
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: batch create associations: %w", err)
	}
	return resp, nil
}

// Read reads associations in bulk for a list of object IDs.
//
//	resp, err := batch.Read(ctx, []crm.AssociationBatchReadInput{
//	    {ID: "contact-1"},
//	    {ID: "contact-2"},
//	})
func (b *AssociationsBatchService) Read(ctx context.Context, inputs []AssociationBatchReadInput) (*AssociationBatchReadResponse, error) {
	resp := &AssociationBatchReadResponse{}
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/read", b.basePath()),
		Body:   map[string]any{"inputs": inputs},
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: batch read associations: %w", err)
	}
	return resp, nil
}

// Archive removes associations in bulk.
//
//	err := batch.Archive(ctx, []crm.AssociationBatchArchiveInput{
//	    {
//	        From: crm.AssociationObjectID{ID: "contact-1"},
//	        To:   crm.AssociationObjectID{ID: "company-1"},
//	    },
//	})
func (b *AssociationsBatchService) Archive(ctx context.Context, inputs []AssociationBatchArchiveInput) error {
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/archive", b.basePath()),
		Body:   map[string]any{"inputs": inputs},
	}, nil); err != nil {
		return fmt.Errorf("hubspot: batch archive associations: %w", err)
	}
	return nil
}

// CreateDefault creates default (unlabeled) associations in bulk.
//
//	resp, err := batch.CreateDefault(ctx, []crm.AssociationBatchArchiveInput{
//	    {
//	        From: crm.AssociationObjectID{ID: "contact-1"},
//	        To:   crm.AssociationObjectID{ID: "company-1"},
//	    },
//	})
func (b *AssociationsBatchService) CreateDefault(ctx context.Context, inputs []AssociationBatchArchiveInput) (*DefaultAssociationResponse, error) {
	resp := &DefaultAssociationResponse{}
	if err := b.r.Do(ctx, &hubspot.Request{
		Method: "POST",
		Path:   fmt.Sprintf("%s/batch/associate/default", b.basePath()),
		Body:   map[string]any{"inputs": inputs},
	}, resp); err != nil {
		return nil, fmt.Errorf("hubspot: batch create default associations: %w", err)
	}
	return resp, nil
}

func (b *AssociationsBatchService) basePath() string {
	return fmt.Sprintf("/crm/v4/associations/%s/%s", b.fromObjectType, b.toObjectType)
}

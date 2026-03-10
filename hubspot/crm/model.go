package crm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Object represents any CRM object (contact, company, deal, ticket, etc.).
// Properties are returned as string key-value pairs because HubSpot objects
// support dynamic custom properties.
type Object struct {
	// ID is the unique HubSpot identifier for this object.
	ID string `json:"id"`

	// Properties contains the object's property values keyed by internal
	// property name (e.g., "email", "firstname", "hs_object_id").
	Properties map[string]string `json:"properties"`

	// PropertiesWithHistory contains property values with their full
	// change history, when requested via [GetOptions.PropertiesWithHistory].
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory,omitempty"`

	// Associations contains associated object IDs, when requested via
	// [GetOptions.Associations] or [ListOptions.Associations].
	Associations map[string]AssociationList `json:"associations,omitempty"`

	// CreatedAt is when the object was created.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is when the object was last modified.
	UpdatedAt time.Time `json:"updatedAt"`

	// Archived indicates whether the object has been soft-deleted.
	Archived bool `json:"archived"`

	// ArchivedAt is when the object was archived, if applicable.
	ArchivedAt *time.Time `json:"archivedAt,omitempty"`
}

// UnmarshalProperties decodes the object's [Properties] map into a typed
// struct. The struct should use json tags matching HubSpot's internal
// property names.
//
//	type Contact struct {
//	    Email     string `json:"email"`
//	    FirstName string `json:"firstname"`
//	    LastName  string `json:"lastname"`
//	}
//
//	var c Contact
//	err := object.UnmarshalProperties(&c)
func (o *Object) UnmarshalProperties(v any) error {
	data, err := json.Marshal(o.Properties)
	if err != nil {
		return fmt.Errorf("marshal properties: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal properties: %w", err)
	}

	return nil
}

// PropertyHistory records a single historical value for a property.
type PropertyHistory struct {
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	SourceID  string    `json:"sourceId"`
	Source    string    `json:"sourceType"`
}

// AssociationList contains associated object references.
type AssociationList struct {
	Results []Association `json:"results"`
}

// Association represents a link between two CRM objects.
type Association struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Paging contains pagination cursors for list responses.
type Paging struct {
	Next *PagingCursor `json:"next,omitempty"`
}

// PagingCursor holds the cursor value for the next page.
type PagingCursor struct {
	After string `json:"after"`
	Link  string `json:"link,omitempty"`
}

// GetOptions configures a [ObjectsService.Get] request.
type GetOptions struct {
	// Properties is the list of property names to return.
	// If empty, HubSpot returns default properties.
	Properties []string

	// PropertiesWithHistory is the list of property names to return
	// with their full change history.
	PropertiesWithHistory []string

	// Associations is the list of object types to return associated IDs for.
	Associations []string

	// Archived includes archived (soft-deleted) objects when true.
	Archived bool

	// IdProperty is the name of a property whose value identifies the object
	// instead of the standard hs_object_id. For example, set IdProperty to
	// "email" to look up a contact by email address.
	IdProperty string
}

func (o *GetOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}

	q := url.Values{}

	if len(o.Properties) > 0 {
		q.Set("properties", strings.Join(o.Properties, ","))
	}
	if len(o.PropertiesWithHistory) > 0 {
		q.Set("propertiesWithHistory", strings.Join(o.PropertiesWithHistory, ","))
	}
	if len(o.Associations) > 0 {
		q.Set("associations", strings.Join(o.Associations, ","))
	}
	if o.Archived {
		q.Set("archived", "true")
	}
	if o.IdProperty != "" {
		q.Set("idProperty", o.IdProperty)
	}

	return q
}

// ListOptions configures a [ObjectsService.List] request.
type ListOptions struct {
	// Limit is the maximum number of results per page (max 100).
	Limit int

	// After is the paging cursor from a previous response.
	After string

	// Properties is the list of property names to return.
	Properties []string

	// PropertiesWithHistory is the list of property names to return
	// with their full change history.
	PropertiesWithHistory []string

	// Associations is the list of object types to return associated IDs for.
	Associations []string

	// Archived includes archived (soft-deleted) objects when true.
	Archived bool
}

func (o *ListOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}

	q := url.Values{}

	if o.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", o.Limit))
	}
	if o.After != "" {
		q.Set("after", o.After)
	}
	if len(o.Properties) > 0 {
		q.Set("properties", strings.Join(o.Properties, ","))
	}
	if len(o.PropertiesWithHistory) > 0 {
		q.Set("propertiesWithHistory", strings.Join(o.PropertiesWithHistory, ","))
	}
	if len(o.Associations) > 0 {
		q.Set("associations", strings.Join(o.Associations, ","))
	}
	if o.Archived {
		q.Set("archived", "true")
	}

	return q
}

// ListResponse is the response from a list or search operation.
type ListResponse struct {
	// Results contains the objects returned.
	Results []Object `json:"results"`

	// Paging contains the cursor for fetching the next page.
	Paging *Paging `json:"paging,omitempty"`

	// Total is the total number of matching results (only present in
	// search responses).
	Total int `json:"total,omitempty"`
}

// CreateInput is the request body for creating a CRM object.
type CreateInput struct {
	// Properties contains the property values to set on the new object.
	Properties map[string]string `json:"properties"`

	// Associations optionally creates associations during object creation.
	Associations []CreateAssociation `json:"associations,omitempty"`
}

// CreateAssociation defines an association to create alongside a new object.
type CreateAssociation struct {
	To    CreateAssociationTarget `json:"to"`
	Types []AssociationType       `json:"types"`
}

// CreateAssociationTarget identifies the object to associate with.
type CreateAssociationTarget struct {
	ID string `json:"id"`
}

// AssociationType identifies the type of association.
type AssociationType struct {
	AssociationCategory string `json:"associationCategory"`
	AssociationTypeID   int    `json:"associationTypeId"`
}

// UpdateInput is the request body for updating a CRM object.
type UpdateInput struct {
	// Properties contains the property values to update. Only specified
	// properties are modified; others are left unchanged.
	Properties map[string]string `json:"properties"`
}

// MergeInput is the request body for merging two CRM objects.
type MergeInput struct {
	// PrimaryObjectID is the ID of the object that will remain after the merge.
	PrimaryObjectID string `json:"primaryObjectId"`

	// ObjectIDToMerge is the ID of the object that will be merged into the primary.
	ObjectIDToMerge string `json:"objectIdToMerge"`
}

// GDPRDeleteInput is the request body for GDPR-compliant permanent deletion.
type GDPRDeleteInput struct {
	// ObjectID is the ID of the object to permanently delete.
	ObjectID string `json:"objectId"`

	// IDProperty is the property name used to identify the object
	// (e.g., "email" for contacts). If empty, "objectId" is used.
	IDProperty string `json:"idProperty,omitempty"`
}

// --- Batch types ---

// BatchCreateInput is a single item in a batch create request.
type BatchCreateInput struct {
	Properties   map[string]string   `json:"properties"`
	Associations []CreateAssociation `json:"associations,omitempty"`
}

// BatchReadInput identifies objects to read in a batch request.
type BatchReadInput struct {
	// Properties is the list of property names to return for each object.
	Properties []string `json:"properties"`

	// PropertiesWithHistory is the list of property names to return with history.
	PropertiesWithHistory []string `json:"propertiesWithHistory,omitempty"`

	// IDProperty is the property name used to identify objects.
	// Defaults to "hs_object_id" if empty.
	IDProperty string `json:"idProperty,omitempty"`

	// Inputs contains the IDs of objects to read.
	Inputs []BatchReadID `json:"inputs"`
}

// BatchReadID identifies a single object to read in a batch.
type BatchReadID struct {
	ID string `json:"id"`
}

// BatchUpdateInput is a single item in a batch update request.
type BatchUpdateInput struct {
	ID         string            `json:"id"`
	Properties map[string]string `json:"properties"`
}

// BatchUpsertInput is a single item in a batch upsert request.
type BatchUpsertInput struct {
	// IDProperty is the property used to match existing objects for upsert.
	IDProperty string `json:"idProperty,omitempty"`

	// ID is the value of the IDProperty to match.
	ID string `json:"id"`

	// Properties contains the property values to set.
	Properties map[string]string `json:"properties"`
}

// BatchArchiveInput is a single item in a batch archive request.
type BatchArchiveInput struct {
	ID string `json:"id"`
}

// BatchResponse is the response from a batch create, update, or upsert operation.
type BatchResponse struct {
	Status      string   `json:"status"`
	Results     []Object `json:"results"`
	RequestedAt string   `json:"requestedAt,omitempty"`
	StartedAt   string   `json:"startedAt"`
	CompletedAt string   `json:"completedAt"`
}

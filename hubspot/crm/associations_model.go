package crm

// Association category constants for the v4 Associations API.
const (
	AssociationCategoryHubSpotDefined    = "HUBSPOT_DEFINED"
	AssociationCategoryUserDefined       = "USER_DEFINED"
	AssociationCategoryIntegratorDefined = "INTEGRATOR_DEFINED"
)

// --- CRUD types ---

// AssociationTypeInput specifies an association type when creating associations.
type AssociationTypeInput struct {
	AssociationCategory string `json:"associationCategory"`
	AssociationTypeID   int    `json:"associationTypeId"`
}

// AssociationLabel describes a single association type on a result.
type AssociationLabel struct {
	Category string `json:"category"`
	TypeID   int    `json:"typeId"`
	Label    string `json:"label,omitempty"`
}

// AssociationListResult is a single associated object with its association types.
type AssociationListResult struct {
	ToObjectID       string             `json:"toObjectId"`
	AssociationTypes []AssociationLabel `json:"associationTypes"`
}

// AssociationListResponse is the response from listing associations.
type AssociationListResponse struct {
	Results []AssociationListResult `json:"results"`
	Paging  *Paging                 `json:"paging,omitempty"`
}

// AssociationResult is the response from creating an association.
type AssociationResult struct {
	FromObjectTypeID string             `json:"fromObjectTypeId"`
	FromObjectID     string             `json:"fromObjectId"`
	ToObjectTypeID   string             `json:"toObjectTypeId"`
	ToObjectID       string             `json:"toObjectId"`
	Labels           []AssociationLabel `json:"labels"`
}

// --- Batch types ---

// AssociationObjectID wraps an object ID for batch request bodies.
type AssociationObjectID struct {
	ID string `json:"id"`
}

// AssociationBatchCreateInput is a single item in a batch association create request.
type AssociationBatchCreateInput struct {
	From  AssociationObjectID    `json:"from"`
	To    AssociationObjectID    `json:"to"`
	Types []AssociationTypeInput `json:"types"`
}

// AssociationBatchReadInput identifies a single object to read associations for.
type AssociationBatchReadInput struct {
	ID string `json:"id"`
}

// AssociationBatchArchiveInput identifies a pair of objects whose association should be removed.
type AssociationBatchArchiveInput struct {
	From AssociationObjectID `json:"from"`
	To   AssociationObjectID `json:"to"`
}

// AssociationBatchCreateResponse is the response from a batch association create.
type AssociationBatchCreateResponse struct {
	Status      string              `json:"status"`
	Results     []AssociationResult `json:"results"`
	StartedAt   string              `json:"startedAt"`
	CompletedAt string              `json:"completedAt"`
}

// AssociationBatchReadResult is a single object's associations in a batch read response.
type AssociationBatchReadResult struct {
	From   AssociationObjectID     `json:"from"`
	To     []AssociationListResult `json:"to"`
	Paging *Paging                 `json:"paging,omitempty"`
}

// AssociationBatchReadResponse is the response from a batch association read.
type AssociationBatchReadResponse struct {
	Status      string                       `json:"status"`
	Results     []AssociationBatchReadResult `json:"results"`
	StartedAt   string                       `json:"startedAt"`
	CompletedAt string                       `json:"completedAt"`
}

// --- Schema types ---

// AssociationLabelDefinition describes a custom association label/definition.
type AssociationLabelDefinition struct {
	Category string `json:"category"`
	TypeID   int    `json:"typeId"`
	Label    string `json:"label"`
}

// AssociationLabelListResponse is the response from listing association label definitions.
type AssociationLabelListResponse struct {
	Results []AssociationLabelDefinition `json:"results"`
}

// CreateAssociationLabelInput is the request body for creating a custom association label.
type CreateAssociationLabelInput struct {
	Label        string `json:"label"`
	InverseLabel string `json:"inverseLabel,omitempty"`
	Name         string `json:"name"`
}

// CreateAssociationLabelResponse is the response from creating a custom association label.
type CreateAssociationLabelResponse struct {
	Results []AssociationLabelDefinition `json:"results"`
}

// UpdateAssociationLabelInput is the request body for updating a custom association label.
type UpdateAssociationLabelInput struct {
	AssociationTypeID int    `json:"associationTypeId"`
	Label             string `json:"label"`
}

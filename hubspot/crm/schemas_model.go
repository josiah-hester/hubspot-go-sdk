package crm

import (
	"net/url"
	"time"
)

// --- Response types ---

// Schema is a full custom object schema as returned by Get, Create, and List.
type Schema struct {
	ID                         string                 `json:"id"`
	Name                       string                 `json:"name"`
	Labels                     SchemaLabels           `json:"labels"`
	Description                string                 `json:"description,omitempty"`
	ObjectTypeID               string                 `json:"objectTypeId,omitempty"`
	FullyQualifiedName         string                 `json:"fullyQualifiedName,omitempty"`
	PrimaryDisplayProperty     string                 `json:"primaryDisplayProperty,omitempty"`
	SecondaryDisplayProperties []string               `json:"secondaryDisplayProperties,omitempty"`
	RequiredProperties         []string               `json:"requiredProperties"`
	SearchableProperties       []string               `json:"searchableProperties,omitempty"`
	Properties                 []SchemaProperty       `json:"properties"`
	Associations               []SchemaAssociationDef `json:"associations"`
	AllowsSensitiveProperties  bool                   `json:"allowsSensitiveProperties,omitempty"`
	Archived                   bool                   `json:"archived,omitempty"`
	CreatedAt                  time.Time              `json:"createdAt,omitempty"`
	UpdatedAt                  time.Time              `json:"updatedAt,omitempty"`
	CreatedByUserID            int                    `json:"createdByUserId,omitempty"`
	UpdatedByUserID            int                    `json:"updatedByUserId,omitempty"`
}

// SchemaDefinition is the schema metadata returned by Update (without
// properties or associations arrays).
type SchemaDefinition struct {
	ID                         string       `json:"id"`
	Name                       string       `json:"name"`
	Labels                     SchemaLabels `json:"labels"`
	Description                string       `json:"description,omitempty"`
	ObjectTypeID               string       `json:"objectTypeId,omitempty"`
	FullyQualifiedName         string       `json:"fullyQualifiedName,omitempty"`
	PrimaryDisplayProperty     string       `json:"primaryDisplayProperty,omitempty"`
	SecondaryDisplayProperties []string     `json:"secondaryDisplayProperties,omitempty"`
	RequiredProperties         []string     `json:"requiredProperties"`
	SearchableProperties       []string     `json:"searchableProperties,omitempty"`
	PortalID                   int          `json:"portalId,omitempty"`
	AllowsSensitiveProperties  bool         `json:"allowsSensitiveProperties,omitempty"`
	Archived                   bool         `json:"archived,omitempty"`
	CreatedAt                  time.Time    `json:"createdAt,omitempty"`
	UpdatedAt                  time.Time    `json:"updatedAt,omitempty"`
}

// SchemaListResponse is the response from listing all schemas.
type SchemaListResponse struct {
	Results []Schema `json:"results"`
}

// --- Nested types ---

// SchemaLabels holds the singular and plural display names for a schema.
type SchemaLabels struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
}

// SchemaProperty is a full property definition as returned by the API.
type SchemaProperty struct {
	Name                 string                 `json:"name"`
	Label                string                 `json:"label"`
	Type                 string                 `json:"type"`
	FieldType            string                 `json:"fieldType"`
	Description          string                 `json:"description"`
	GroupName            string                 `json:"groupName"`
	Options              []SchemaPropertyOption `json:"options"`
	DisplayOrder         int                    `json:"displayOrder,omitempty"`
	HasUniqueValue       bool                   `json:"hasUniqueValue,omitempty"`
	Hidden               bool                   `json:"hidden,omitempty"`
	FormField            bool                   `json:"formField,omitempty"`
	Archived             bool                   `json:"archived,omitempty"`
	Calculated           bool                   `json:"calculated,omitempty"`
	ExternalOptions      bool                   `json:"externalOptions,omitempty"`
	HubSpotDefined       bool                   `json:"hubspotDefined,omitempty"`
	ShowCurrencySymbol   bool                   `json:"showCurrencySymbol,omitempty"`
	ReferencedObjectType string                 `json:"referencedObjectType,omitempty"`
	CalculationFormula   string                 `json:"calculationFormula,omitempty"`
	DataSensitivity      string                 `json:"dataSensitivity,omitempty"`
	ModificationMetadata *SchemaPropertyModMeta `json:"modificationMetadata,omitempty"`
	CreatedAt            time.Time              `json:"createdAt,omitempty"`
	UpdatedAt            time.Time              `json:"updatedAt,omitempty"`
	ArchivedAt           *time.Time             `json:"archivedAt,omitempty"`
	CreatedUserId        string                 `json:"createdUserId,omitempty"`
	UpdatedUserId        string                 `json:"updatedUserId,omitempty"`
}

// SchemaPropertyOption is a single option for an enumeration property.
type SchemaPropertyOption struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
	Hidden       bool   `json:"hidden"`
}

// SchemaPropertyModMeta holds modification metadata for a property.
type SchemaPropertyModMeta struct {
	Archivable         bool `json:"archivable"`
	ReadOnlyDefinition bool `json:"readOnlyDefinition"`
	ReadOnlyOptions    bool `json:"readOnlyOptions,omitempty"`
	ReadOnlyValue      bool `json:"readOnlyValue"`
}

// SchemaAssociationDef is an association definition on a schema.
type SchemaAssociationDef struct {
	ID               string    `json:"id"`
	FromObjectTypeID string    `json:"fromObjectTypeId"`
	ToObjectTypeID   string    `json:"toObjectTypeId"`
	Name             string    `json:"name,omitempty"`
	CreatedAt        time.Time `json:"createdAt,omitempty"`
	UpdatedAt        time.Time `json:"updatedAt,omitempty"`
}

// --- Input types ---

// SchemaCreateInput is the request body for creating a new schema.
type SchemaCreateInput struct {
	Name                       string                 `json:"name"`
	Labels                     SchemaLabels           `json:"labels"`
	Description                string                 `json:"description,omitempty"`
	Properties                 []SchemaPropertyCreate `json:"properties"`
	AssociatedObjects          []string               `json:"associatedObjects"`
	RequiredProperties         []string               `json:"requiredProperties"`
	SearchableProperties       []string               `json:"searchableProperties,omitempty"`
	PrimaryDisplayProperty     string                 `json:"primaryDisplayProperty,omitempty"`
	SecondaryDisplayProperties []string               `json:"secondaryDisplayProperties,omitempty"`
	AllowsSensitiveProperties  bool                   `json:"allowsSensitiveProperties,omitempty"`
}

// SchemaPropertyCreate is a property definition within a create schema request.
type SchemaPropertyCreate struct {
	Name           string                 `json:"name"`
	Label          string                 `json:"label"`
	Type           string                 `json:"type"`
	FieldType      string                 `json:"fieldType"`
	Description    string                 `json:"description,omitempty"`
	GroupName      string                 `json:"groupName,omitempty"`
	Options        []SchemaPropertyOption `json:"options,omitempty"`
	DisplayOrder   int                    `json:"displayOrder,omitempty"`
	HasUniqueValue bool                   `json:"hasUniqueValue,omitempty"`
	Hidden         bool                   `json:"hidden,omitempty"`
	FormField      bool                   `json:"formField,omitempty"`
}

// SchemaUpdateInput is the request body for updating an existing schema.
type SchemaUpdateInput struct {
	Labels                     *SchemaLabels `json:"labels,omitempty"`
	Description                string        `json:"description,omitempty"`
	ClearDescription           bool          `json:"clearDescription,omitempty"`
	PrimaryDisplayProperty     string        `json:"primaryDisplayProperty,omitempty"`
	RequiredProperties         []string      `json:"requiredProperties,omitempty"`
	SearchableProperties       []string      `json:"searchableProperties,omitempty"`
	SecondaryDisplayProperties []string      `json:"secondaryDisplayProperties,omitempty"`
	AllowsSensitiveProperties  *bool         `json:"allowsSensitiveProperties,omitempty"`
	Restorable                 *bool         `json:"restorable,omitempty"`
}

// SchemaAssociationCreateInput is the request body for adding an association
// definition to a schema.
type SchemaAssociationCreateInput struct {
	FromObjectTypeID string `json:"fromObjectTypeId"`
	ToObjectTypeID   string `json:"toObjectTypeId"`
	Name             string `json:"name"`
}

// --- Options types ---

// SchemaListOptions are query parameters for listing schemas.
type SchemaListOptions struct {
	Archived                      bool
	IncludeAssociationDefinitions bool // default true on API
	IncludeAuditMetadata          bool // default true on API
	IncludePropertyDefinitions    bool // default true on API
}

func (o *SchemaListOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	q := url.Values{}
	if o.Archived {
		q.Set("archived", "true")
	}
	if !o.IncludeAssociationDefinitions {
		q.Set("includeAssociationDefinitions", "false")
	}
	if !o.IncludeAuditMetadata {
		q.Set("includeAuditMetadata", "false")
	}
	if !o.IncludePropertyDefinitions {
		q.Set("includePropertyDefinitions", "false")
	}
	if len(q) == 0 {
		return nil
	}
	return q
}

// SchemaGetOptions are query parameters for getting a single schema.
type SchemaGetOptions struct {
	IncludeAssociationDefinitions bool // default true on API
	IncludeAuditMetadata          bool // default true on API
	IncludePropertyDefinitions    bool // default true on API
}

func (o *SchemaGetOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	q := url.Values{}
	if !o.IncludeAssociationDefinitions {
		q.Set("includeAssociationDefinitions", "false")
	}
	if !o.IncludeAuditMetadata {
		q.Set("includeAuditMetadata", "false")
	}
	if !o.IncludePropertyDefinitions {
		q.Set("includePropertyDefinitions", "false")
	}
	if len(q) == 0 {
		return nil
	}
	return q
}

// SchemaDeleteOptions are query parameters for deleting a schema.
type SchemaDeleteOptions struct {
	Archived bool // true = hard delete (permanent)
}

func (o *SchemaDeleteOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	if o.Archived {
		return url.Values{"archived": {"true"}}
	}
	return nil
}

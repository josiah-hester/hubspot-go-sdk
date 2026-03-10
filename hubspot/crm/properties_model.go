package crm

import (
	"net/url"
	"time"
)

// --- Response types ---

// Property is a CRM property definition as returned by the Properties API.
type Property struct {
	Name                    string           `json:"name"`
	Label                   string           `json:"label"`
	Type                    string           `json:"type"`
	FieldType               string           `json:"fieldType"`
	Description             string           `json:"description"`
	GroupName               string           `json:"groupName"`
	Options                 []PropertyOption `json:"options"`
	DisplayOrder            int              `json:"displayOrder,omitempty"`
	HasUniqueValue          bool             `json:"hasUniqueValue,omitempty"`
	Hidden                  bool             `json:"hidden,omitempty"`
	FormField               bool             `json:"formField,omitempty"`
	Archived                bool             `json:"archived,omitempty"`
	Calculated              bool             `json:"calculated,omitempty"`
	ExternalOptions         bool             `json:"externalOptions,omitempty"`
	HubSpotDefined          bool             `json:"hubspotDefined,omitempty"`
	ShowCurrencySymbol      bool             `json:"showCurrencySymbol,omitempty"`
	ReferencedObjectType    string           `json:"referencedObjectType,omitempty"`
	CalculationFormula      string           `json:"calculationFormula,omitempty"`
	DataSensitivity         string           `json:"dataSensitivity,omitempty"`
	DateDisplayHint         string           `json:"dateDisplayHint,omitempty"`
	SensitiveDataCategories []string         `json:"sensitiveDataCategories,omitempty"`
	ModificationMetadata    *PropertyModMeta `json:"modificationMetadata,omitempty"`
	CreatedAt               time.Time        `json:"createdAt,omitempty"`
	UpdatedAt               time.Time        `json:"updatedAt,omitempty"`
	ArchivedAt              *time.Time       `json:"archivedAt,omitempty"`
	CreatedUserId           string           `json:"createdUserId,omitempty"`
	UpdatedUserId           string           `json:"updatedUserId,omitempty"`
}

// PropertyOption is a single option for an enumeration property.
type PropertyOption struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
	Hidden       bool   `json:"hidden"`
}

// PropertyModMeta holds modification metadata for a property.
type PropertyModMeta struct {
	Archivable         bool `json:"archivable"`
	ReadOnlyDefinition bool `json:"readOnlyDefinition"`
	ReadOnlyOptions    bool `json:"readOnlyOptions,omitempty"`
	ReadOnlyValue      bool `json:"readOnlyValue"`
}

// PropertyListResponse is the response from listing all properties for an object type.
type PropertyListResponse struct {
	Results []Property `json:"results"`
}

// PropertyGroup is a property group definition.
type PropertyGroup struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	DisplayOrder int    `json:"displayOrder"`
	Archived     bool   `json:"archived"`
}

// PropertyGroupListResponse is the response from listing all property groups.
type PropertyGroupListResponse struct {
	Results []PropertyGroup `json:"results"`
}

// PropertyBatchResponse is the response from batch property operations.
type PropertyBatchResponse struct {
	Status      string     `json:"status"`
	Results     []Property `json:"results"`
	StartedAt   string     `json:"startedAt"`
	CompletedAt string     `json:"completedAt"`
	NumErrors   int        `json:"numErrors,omitempty"`
}

// --- Input types ---

// PropertyCreateInput is the request body for creating a property.
type PropertyCreateInput struct {
	Name                 string           `json:"name"`
	Label                string           `json:"label"`
	Type                 string           `json:"type"`
	FieldType            string           `json:"fieldType"`
	GroupName            string           `json:"groupName"`
	Description          string           `json:"description,omitempty"`
	DisplayOrder         int              `json:"displayOrder,omitempty"`
	Options              []PropertyOption `json:"options,omitempty"`
	Hidden               bool             `json:"hidden,omitempty"`
	FormField            bool             `json:"formField,omitempty"`
	HasUniqueValue       bool             `json:"hasUniqueValue,omitempty"`
	ExternalOptions      bool             `json:"externalOptions,omitempty"`
	ReferencedObjectType string           `json:"referencedObjectType,omitempty"`
	CalculationFormula   string           `json:"calculationFormula,omitempty"`
	DataSensitivity      string           `json:"dataSensitivity,omitempty"`
}

// PropertyUpdateInput is the request body for updating a property.
type PropertyUpdateInput struct {
	Label              string           `json:"label,omitempty"`
	Type               string           `json:"type,omitempty"`
	FieldType          string           `json:"fieldType,omitempty"`
	GroupName          string           `json:"groupName,omitempty"`
	Description        string           `json:"description,omitempty"`
	DisplayOrder       int              `json:"displayOrder,omitempty"`
	Options            []PropertyOption `json:"options,omitempty"`
	Hidden             *bool            `json:"hidden,omitempty"`
	FormField          *bool            `json:"formField,omitempty"`
	CalculationFormula string           `json:"calculationFormula,omitempty"`
}

// PropertyGroupCreateInput is the request body for creating a property group.
type PropertyGroupCreateInput struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// PropertyGroupUpdateInput is the request body for updating a property group.
type PropertyGroupUpdateInput struct {
	Label        string `json:"label,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// PropertyBatchReadInput identifies a property by name for batch read.
type PropertyBatchReadInput struct {
	Name string `json:"name"`
}

// PropertyBatchReadRequest is the request body for batch reading properties.
type PropertyBatchReadRequest struct {
	Inputs          []PropertyBatchReadInput `json:"inputs"`
	Archived        bool                     `json:"archived"`
	DataSensitivity string                   `json:"dataSensitivity,omitempty"`
}

// --- Options types ---

// PropertyListOptions are query parameters for listing properties.
type PropertyListOptions struct {
	Archived        bool
	DataSensitivity string
	Properties      string
}

func (o *PropertyListOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	q := url.Values{}
	if o.Archived {
		q.Set("archived", "true")
	}
	if o.DataSensitivity != "" {
		q.Set("dataSensitivity", o.DataSensitivity)
	}
	if o.Properties != "" {
		q.Set("properties", o.Properties)
	}
	if len(q) == 0 {
		return nil
	}
	return q
}

// PropertyGetOptions are query parameters for getting a single property.
type PropertyGetOptions struct {
	Archived        bool
	DataSensitivity string
	Properties      string
}

func (o *PropertyGetOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	q := url.Values{}
	if o.Archived {
		q.Set("archived", "true")
	}
	if o.DataSensitivity != "" {
		q.Set("dataSensitivity", o.DataSensitivity)
	}
	if o.Properties != "" {
		q.Set("properties", o.Properties)
	}
	if len(q) == 0 {
		return nil
	}
	return q
}

// PropertyGroupListOptions are query parameters for listing property groups.
type PropertyGroupListOptions struct {
	Locale string
}

func (o *PropertyGroupListOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	if o.Locale != "" {
		return url.Values{"locale": {o.Locale}}
	}
	return nil
}

// PropertyGroupGetOptions are query parameters for getting a single property group.
type PropertyGroupGetOptions struct {
	Locale string
}

func (o *PropertyGroupGetOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	if o.Locale != "" {
		return url.Values{"locale": {o.Locale}}
	}
	return nil
}

// PropertyBatchReadOptions are query parameters for batch reading properties.
type PropertyBatchReadOptions struct {
	Locale string
}

func (o *PropertyBatchReadOptions) toQuery() url.Values {
	if o == nil {
		return nil
	}
	if o.Locale != "" {
		return url.Values{"locale": {o.Locale}}
	}
	return nil
}

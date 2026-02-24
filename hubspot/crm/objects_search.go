package crm

// Operator constants for CRM search filters.
const (
	// OpEQ matches objects where the property equals the value.
	OpEQ = "EQ"

	// OpNEQ matches objects where the property does not equal the value.
	OpNEQ = "NEQ"

	// OpLT matches objects where the property is less than the value.
	OpLT = "LT"

	// OpLTE matches objects where the property is less than or equal to the value.
	OpLTE = "LTE"

	// OpGT matches objects where the property is greater than the value.
	OpGT = "GT"

	// OpGTE matches objects where the property is greater than or equal to the value.
	OpGTE = "GTE"

	// OpBetween matches objects where the property is between highValue and value.
	OpBetween = "BETWEEN"

	// OpIN matches objects where the property is one of the specified values.
	OpIN = "IN"

	// OpNotIN matches objects where the property is not one of the specified values.
	OpNotIN = "NOT_IN"

	// OpHasProperty matches objects that have a value for the property.
	OpHasProperty = "HAS_PROPERTY"

	// OpNotHasProperty matches objects that do not have a value for the property.
	OpNotHasProperty = "NOT_HAS_PROPERTY"

	// OpContainsToken matches objects where the property contains the token.
	OpContainsToken = "CONTAINS_TOKEN"

	// OpNotContainsToken matches objects where the property does not contain the token.
	OpNotContainsToken = "NOT_CONTAINS_TOKEN"
)

// SearchRequest defines a CRM object search query.
type SearchRequest struct {
	// FilterGroups are ORed together. Within each group, filters are ANDed.
	FilterGroups []FilterGroup `json:"filterGroups,omitempty"`

	// Sorts defines the sort order of results.
	Sorts []Sort `json:"sorts,omitempty"`

	// Properties is the list of property names to return.
	Properties []string `json:"properties,omitempty"`

	// Limit is the maximum number of results to return (max 100).
	Limit int `json:"limit,omitempty"`

	// After is the paging cursor for fetching the next page of results.
	After string `json:"after,omitempty"`
}

// FilterGroup contains filters that are ANDed together. Multiple filter
// groups in a [SearchRequest] are ORed.
type FilterGroup struct {
	Filters []Filter `json:"filters"`
}

// Filter is a single property filter in a search query.
type Filter struct {
	// PropertyName is the internal name of the property to filter on.
	PropertyName string `json:"propertyName"`

	// Operator is the comparison operator (use Op* constants).
	Operator string `json:"operator"`

	// Value is the filter value. Use for single-value operators.
	Value string `json:"value,omitempty"`

	// HighValue is the upper bound for [OpBetween] queries.
	HighValue string `json:"highValue,omitempty"`

	// Values is used for [OpIN] and [OpNotIN] operators.
	Values []string `json:"values,omitempty"`
}

// Sort defines a sort criterion for search results.
type Sort struct {
	// PropertyName is the property to sort by.
	PropertyName string `json:"propertyName"`

	// Direction is the sort direction: "ASCENDING" or "DESCENDING".
	Direction string `json:"direction"`
}

// SortAscending is the ascending sort direction.
const SortAscending = "ASCENDING"

// SortDescending is the descending sort direction.
const SortDescending = "DESCENDING"

// --- Fluent search builder ---

// SearchBuilder provides a fluent interface for constructing CRM search
// queries.
//
//	req := crm.NewSearch().
//	    Where("email", crm.OpContainsToken, "example.com").
//	    Where("firstname", crm.OpEQ, "Alice").
//	    Or().
//	    Where("email", crm.OpContainsToken, "test.com").
//	    SortBy("createdate", crm.SortDescending).
//	    Select("email", "firstname", "lastname").
//	    Limit(20).
//	    Build()
type SearchBuilder struct {
	groups     []FilterGroup
	current    []Filter
	sorts      []Sort
	properties []string
	limit      int
	after      string
}

// NewSearch creates a new [SearchBuilder]. Filters added with [Where]
// are ANDed within a group. Call [Or] to start a new group (OR condition).
func NewSearch() *SearchBuilder {
	return &SearchBuilder{}
}

// Where adds a filter to the current filter group. Multiple calls to Where
// without an intervening [Or] are ANDed together.
func (b *SearchBuilder) Where(property, operator, value string) *SearchBuilder {
	b.current = append(b.current, Filter{
		PropertyName: property,
		Operator:     operator,
		Value:        value,
	})
	return b
}

// WhereIn adds an IN filter for matching against multiple values.
func (b *SearchBuilder) WhereIn(property string, values ...string) *SearchBuilder {
	b.current = append(b.current, Filter{
		PropertyName: property,
		Operator:     OpIN,
		Values:       values,
	})
	return b
}

// WhereNotIn adds a NOT_IN filter for excluding multiple values.
func (b *SearchBuilder) WhereNotIn(property string, values ...string) *SearchBuilder {
	b.current = append(b.current, Filter{
		PropertyName: property,
		Operator:     OpNotIN,
		Values:       values,
	})
	return b
}

// WhereBetween adds a BETWEEN filter with a low and high value.
func (b *SearchBuilder) WhereBetween(property, low, high string) *SearchBuilder {
	b.current = append(b.current, Filter{
		PropertyName: property,
		Operator:     OpBetween,
		Value:        low,
		HighValue:    high,
	})
	return b
}

// WhereHasProperty adds a HAS_PROPERTY filter.
func (b *SearchBuilder) WhereHasProperty(property string) *SearchBuilder {
	b.current = append(b.current, Filter{
		PropertyName: property,
		Operator:     OpHasProperty,
	})
	return b
}

// WhereNotHasProperty adds a NOT_HAS_PROPERTY filter.
func (b *SearchBuilder) WhereNotHasProperty(property string) *SearchBuilder {
	b.current = append(b.current, Filter{
		PropertyName: property,
		Operator:     OpNotHasProperty,
	})
	return b
}

// Or finalizes the current filter group and starts a new one. The groups
// are ORed together in the final query.
func (b *SearchBuilder) Or() *SearchBuilder {
	if len(b.current) > 0 {
		b.groups = append(b.groups, FilterGroup{Filters: b.current})
		b.current = nil
	}
	return b
}

// SortBy adds a sort criterion. Multiple calls add multiple sort keys
// in order of priority.
func (b *SearchBuilder) SortBy(property, direction string) *SearchBuilder {
	b.sorts = append(b.sorts, Sort{
		PropertyName: property,
		Direction:    direction,
	})
	return b
}

// Select specifies which properties to return in the results.
func (b *SearchBuilder) Select(properties ...string) *SearchBuilder {
	b.properties = append(b.properties, properties...)
	return b
}

// Limit sets the maximum number of results (max 100).
func (b *SearchBuilder) Limit(n int) *SearchBuilder {
	b.limit = n
	return b
}

// After sets the paging cursor for fetching the next page.
func (b *SearchBuilder) After(cursor string) *SearchBuilder {
	b.after = cursor
	return b
}

// Build returns the constructed [SearchRequest].
func (b *SearchBuilder) Build() *SearchRequest {
	// Flush any pending filters in the current group.
	groups := b.groups
	if len(b.current) > 0 {
		groups = append(groups, FilterGroup{Filters: b.current})
	}

	return &SearchRequest{
		FilterGroups: groups,
		Sorts:        b.sorts,
		Properties:   b.properties,
		Limit:        b.limit,
		After:        b.after,
	}
}

package crm

import "github.com/josiah-hester/hubspot-go-sdk/hubspot"

// Service provides access to all CRM APIs. Create one by passing a
// [hubspot.Client] (or any [hubspot.Requester]) to [NewService].
//
//	client := hubspot.NewClient(hubspot.PrivateAppToken("pat-na1-xxxxx"))
//	crm := crm.NewService(client)
//	contact, err := crm.Contacts().Get(ctx, "123", nil)
type Service struct {
	r hubspot.Requester
}

// NewService creates a new CRM service.
func NewService(r hubspot.Requester) *Service {
	return &Service{r: r}
}

// Contacts returns a service for the Contacts object type.
func (s *Service) Contacts() *ObjectsService {
	return newObjectsService(s.r, "contacts")
}

// Companies returns a service for the Companies object type.
func (s *Service) Companies() *ObjectsService {
	return newObjectsService(s.r, "companies")
}

// Deals returns a service for the Deals object type.
func (s *Service) Deals() *ObjectsService {
	return newObjectsService(s.r, "deals")
}

// Tickets returns a service for the Tickets object type.
func (s *Service) Tickets() *ObjectsService {
	return newObjectsService(s.r, "tickets")
}

// Products returns a service for the Products object type.
func (s *Service) Products() *ObjectsService {
	return newObjectsService(s.r, "products")
}

// LineItems returns a service for the Line Items object type.
func (s *Service) LineItems() *ObjectsService {
	return newObjectsService(s.r, "line_items")
}

// Quotes returns a service for the Quotes object type.
func (s *Service) Quotes() *ObjectsService {
	return newObjectsService(s.r, "quotes")
}

// Calls returns a service for the Calls engagement type.
func (s *Service) Calls() *ObjectsService {
	return newObjectsService(s.r, "calls")
}

// Emails returns a service for the Emails engagement type.
func (s *Service) Emails() *ObjectsService {
	return newObjectsService(s.r, "emails")
}

// Meetings returns a service for the Meetings engagement type.
func (s *Service) Meetings() *ObjectsService {
	return newObjectsService(s.r, "meetings")
}

// Notes returns a service for the Notes engagement type.
func (s *Service) Notes() *ObjectsService {
	return newObjectsService(s.r, "notes")
}

// Tasks returns a service for the Tasks engagement type.
func (s *Service) Tasks() *ObjectsService {
	return newObjectsService(s.r, "tasks")
}

// Communications returns a service for the Communications object type.
func (s *Service) Communications() *ObjectsService {
	return newObjectsService(s.r, "communications")
}

// FeedbackSubmissions returns a service for the Feedback Submissions object type.
func (s *Service) FeedbackSubmissions() *ObjectsService {
	return newObjectsService(s.r, "feedback_submissions")
}

// Orders returns a service for the Orders object type.
func (s *Service) Orders() *ObjectsService {
	return newObjectsService(s.r, "orders")
}

// Invoices returns a service for the Invoices object type.
func (s *Service) Invoices() *ObjectsService {
	return newObjectsService(s.r, "invoices")
}

// Leads returns a service for the Leads object type.
func (s *Service) Leads() *ObjectsService {
	return newObjectsService(s.r, "leads")
}

// Appointments returns a service for the Appointments object type.
func (s *Service) Appointments() *ObjectsService {
	return newObjectsService(s.r, "appointments")
}

// Object returns a service for any CRM object type by its type name
// or ID. Use this for custom objects or any object type not covered by
// the named accessors above.
//
//	// By custom object name
//	pets := crm.Object("pets")
//
//	// By custom object ID
//	custom := crm.Object("2-12345")
func (s *Service) Object(objectType string) *ObjectsService {
	return newObjectsService(s.r, objectType)
}

// Properties returns a [PropertiesService] for managing property definitions
// on the given object type.
//
//	props := crm.NewService(client).Properties("contacts")
//	resp, err := props.List(ctx, nil)
func (s *Service) Properties(objectType string) *PropertiesService {
	return newPropertiesService(s.r, objectType)
}

// Schemas returns a [SchemasService] for managing custom object schemas.
//
//	schemas := crm.NewService(client).Schemas()
//	resp, err := schemas.List(ctx, nil)
func (s *Service) Schemas() *SchemasService {
	return newSchemasService(s.r)
}

// Associations returns an [AssociationsService] for managing associations
// between two CRM object types using the v4 Associations API.
//
//	assoc := crm.NewService(client).Associations("contacts", "companies")
//	result, err := assoc.Create(ctx, "contact-1", "company-1", types)
func (s *Service) Associations(fromObjectType, toObjectType string) *AssociationsService {
	return newAssociationsService(s.r, fromObjectType, toObjectType)
}

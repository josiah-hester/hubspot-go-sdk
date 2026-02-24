//go:build integration

package crm_test

import (
	"context"
	"os"
	"testing"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
	"github.com/josiah-hester/hubspot-go-sdk/hubspot/crm"
)

func TestIntegration_Contacts_Get(t *testing.T) {
	token := os.Getenv("HUBSPOT_TOKEN")
	if token == "" {
		t.Skip("HUBSPOT_TOKEN not set")
	}

	contactID := os.Getenv("HUBSPOT_TEST_CONTACT_ID")
	if contactID == "" {
		t.Skip("HUBSPOT_TEST_CONTACT_ID not set")
	}

	client := hubspot.NewClient(hubspot.PrivateAppToken(token))
	contacts := crm.NewService(client).Contacts()

	contact, err := contacts.Get(context.Background(), contactID, &crm.GetOptions{
		Properties: []string{"email", "firstname", "lastname"},
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	t.Logf("Contact %s: %v", contact.ID, contact.Properties)

	if contact.ID != contactID {
		t.Errorf("ID = %q, want %q", contact.ID, contactID)
	}
}

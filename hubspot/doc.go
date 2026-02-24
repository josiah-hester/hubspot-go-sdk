// Package hubspot provides a Go client for the HubSpot API.
//
// Create a client with [NewClient] and a [TokenSource], then access
// services through the client's methods:
//
//	client := hubspot.NewClient(hubspot.PrivateAppToken("pat-na1-xxxxx"))
//	contact, err := client.CRM().Contacts().Get(ctx, "123", nil)
package hubspot

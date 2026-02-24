# hubspot-go-sdk

A Go client library for the [HubSpot API](https://developers.hubspot.com/docs/api/overview).

> **⚠️ Work in progress.** This SDK is under active development. The CRM APIs are the current focus.

## Install

```bash
go get github.com/josiah-hester/hubspot-go-sdk
```

Requires Go 1.22+.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/josiah-hester/hubspot-go-sdk/hubspot"
	"github.com/josiah-hester/hubspot-go-sdk/hubspot/crm"
)

func main() {
	client := hubspot.NewClient(hubspot.PrivateAppToken("pat-na1-xxxxx"))

	ctx := context.Background()

	// Get a contact
	contact, err := client.CRM().Contacts().Get(ctx, "123", &crm.GetOptions{
		Properties: []string{"email", "firstname", "lastname"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(contact.Properties["email"])
}
```

## Features

- **Single client, all services** — one `hubspot.Client` gives access to CRM, CMS, Marketing, and more
- **Generic CRM objects** — Contacts, Companies, Deals, Tickets, and all 40+ object types through one consistent interface
- **Automatic rate limiting** — proactive token bucket + adaptive tracking from HubSpot response headers
- **Automatic retries** — exponential backoff with jitter for 429s and server errors
- **Pagination iterator** — `ListAll()` handles cursor management automatically
- **Search builder** — fluent API for constructing CRM search queries
- **Multiple auth methods** — private app tokens and OAuth2 with auto-refresh

## Auth

```go
// Private app token (simplest)
client := hubspot.NewClient(hubspot.PrivateAppToken("pat-na1-xxxxx"))

// OAuth2 with automatic token refresh
client := hubspot.NewClient(hubspot.OAuthTokenSource(hubspot.OAuthConfig{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	RefreshToken: "your-refresh-token",
}))
```

## License

MIT — see [LICENSE](LICENSE).

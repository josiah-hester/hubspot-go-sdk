package hubspot

import "context"

// TokenSource provides authentication tokens for HubSpot API requests.
// Implementations must be safe for concurrent use.
type TokenSource interface {
	// Token returns a valid access token. Implementations may cache
	// tokens and handle refresh logic internally.
	Token(ctx context.Context) (*Token, error)
}

// Token holds an authentication credential for the HubSpot API.
type Token struct {
	// AccessToken is the token value sent in the Authorization header.
	AccessToken string

	// TokenType is the type prefix for the Authorization header
	// (typically "Bearer").
	TokenType string
}

// PrivateAppToken returns a [TokenSource] that provides a static private
// app access token. This is the simplest auth method for server-side
// integrations.
//
// Create a private app token in your HubSpot account under
// Settings > Integrations > Private Apps.
func PrivateAppToken(token string) TokenSource {
	return &staticTokenSource{
		token: Token{
			AccessToken: token,
			TokenType:   "Bearer",
		},
	}
}

type staticTokenSource struct {
	token Token
}

// Token returns the static token provided at construction.
func (s *staticTokenSource) Token(_ context.Context) (*Token, error) {
	return &s.token, nil
}

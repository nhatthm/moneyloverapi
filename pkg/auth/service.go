package auth

import "context"

// TokenProvider provides oauth2 token.
type TokenProvider interface {
	// Token provides a token.
	Token(ctx context.Context) (Token, error)
}

// TokenStorage persists or gets OAuthToken.
type TokenStorage interface {
	// Get gets OAuthToken from data source.
	Get(ctx context.Context, key string) (OAuthToken, error)
	// Set sets OAuthToken to data source.
	Set(ctx context.Context, key string, token OAuthToken) error
}

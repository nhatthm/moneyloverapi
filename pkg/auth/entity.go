package auth

import "time"

// Token is an oauth2 token.
type Token string

// OAuthToken contains all relevant information to access to the service and refresh the token.
type OAuthToken struct {
	AccessToken Token     `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// IsExpired checks whether the access token is expired or not.
func (t OAuthToken) IsExpired(timestamp time.Time) bool {
	return t.ExpiresAt.Before(timestamp)
}

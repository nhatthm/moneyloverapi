package moneyloverapi

import (
	"time"

	"go.nhat.io/clock"

	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

// WithBaseURL sets API Base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.config.baseURL = baseURL
	}
}

// WithTimeout sets API Timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.config.timeout = timeout
	}
}

// WithUsername sets username to login.
func WithUsername(username string) Option {
	return func(c *Client) {
		c.config.username = username
	}
}

// WithPassword sets password to login.
func WithPassword(password string) Option {
	return func(c *Client) {
		c.config.password = password
	}
}

// WithCredentials sets username and password to login.
func WithCredentials(username string, password string) Option {
	return func(c *Client) {
		c.config.username = username
		c.config.password = password
	}
}

// WithCredentialsProvider chains a new credentials provider.
func WithCredentialsProvider(provider CredentialsProvider) Option {
	return func(c *Client) {
		c.config.credentials.prepend(provider)
	}
}

// WithCredentialsProviderAtLast chains a new credentials provider at last position.
func WithCredentialsProviderAtLast(provider CredentialsProvider) Option {
	return func(c *Client) {
		c.config.credentials.append(provider)
	}
}

// WithTokenProvider chains a new token provider.
func WithTokenProvider(provider auth.TokenProvider) Option {
	return func(c *Client) {
		c.token.prepend(provider)
	}
}

// WithTokenStorage sets token storage for the internal apiTokenProvider.
func WithTokenStorage(storage auth.TokenStorage) Option {
	return func(c *Client) {
		c.config.tokenStorage = storage
	}
}

// WithClock sets the clock (for testing purpose).
func WithClock(clock clock.Clock) Option {
	return func(c *Client) {
		c.clock = clock
	}
}

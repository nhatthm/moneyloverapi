//+build integration

package moneyloverapi

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestIntegrationApiTokenProvider tests login functionalities.
//
// In order to run this tests, there env vars should be set:
// - MONEYLOVER_USERNAME: The username to login to MoneyLover, it should be an email address.
// - MONEYLOVER_PASSWORD: The password to login to moneyloverapi.
func TestIntegrationApiTokenProvider(t *testing.T) {
	t.Parallel()

	p := newAPITokenProvider(CredentialsFromEnv()).
		WithBaseURL(BaseURL)

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		token, err := p.Token(context.Background())

		assert.NotEmpty(t, token)
		assert.NoError(t, err)

		close(done)
	}()

	select {
	case <-done:
		return

	case <-ctx.Done():
		t.Fatal("timeout while getting access token")
	}
}

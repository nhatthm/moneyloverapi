package moneyloverapi_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/moneyloverapi"
	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

func TestInMemoryTokenStorage_GetMissingKey(t *testing.T) {
	t.Parallel()

	s := moneyloverapi.NewInMemoryTokenStorage()

	result, err := s.Get(context.Background(), "key")

	assert.Equal(t, auth.OAuthToken{}, result)
	assert.NoError(t, err)
}

func TestInMemoryTokenStorage(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	key := "key"
	originalToken := auth.OAuthToken{
		AccessToken: "access",
		ExpiresAt:   timestamp,
	}

	s := moneyloverapi.NewInMemoryTokenStorage()

	err := s.Set(context.Background(), key, originalToken)
	require.NoError(t, err)

	// Retrieve the token.
	token, err := s.Get(context.Background(), key)
	require.NoError(t, err)

	assert.Equal(t, originalToken, token)

	// Modify the token.
	newAccessToken := "foobar"
	token.AccessToken = auth.Token(newAccessToken)

	// Retrieve the token again.
	token, err = s.Get(context.Background(), key)
	require.NoError(t, err)

	assert.NotEqual(t, newAccessToken, token.AccessToken)
	assert.Equal(t, originalToken, token)
}

func TestInMemoryTokenStorage_SetError(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	originalToken := auth.OAuthToken{
		AccessToken: "access",
		ExpiresAt:   timestamp,
	}

	s := moneyloverapi.NewInMemoryTokenStorage()

	err := s.Set(context.Background(), "", originalToken)

	assert.Equal(t, moneyloverapi.ErrTokenKeyEmpty, err)
}

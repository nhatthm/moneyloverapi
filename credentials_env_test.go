package moneyloverapi_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/moneyloverapi"
)

const (
	envUsername = "MONEYLOVER_USERNAME"
	envPassword = "MONEYLOVER_PASSWORD"
)

func TestCredentialsFromEnv(t *testing.T) {
	currentUsername := os.Getenv(envUsername)
	currentPassword := os.Getenv(envPassword)

	t.Cleanup(func() {
		err := os.Setenv(envUsername, currentUsername)
		require.NoError(t, err)

		err = os.Setenv(envPassword, currentPassword)
		require.NoError(t, err)
	})

	t.Setenv(envUsername, "username")
	t.Setenv(envPassword, "password")

	p := moneyloverapi.CredentialsFromEnv()

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}

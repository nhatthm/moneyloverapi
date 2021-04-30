package moneyloverapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/moneyloverapi"
)

func TestCredentials(t *testing.T) {
	p := moneyloverapi.Credentials("username", "password")

	assert.Equal(t, "username", p.Username())
	assert.Equal(t, "password", p.Password())
}

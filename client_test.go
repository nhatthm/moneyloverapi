package moneyloverapi_test

import (
	"github.com/nhatthm/moneyloverapi/pkg/testkit"
)

var (
	mlUsername = "user@example.org"
	mlPassword = "123456"
)

func mockServer(mocks ...testkit.ServerOption) testkit.ServerMocker {
	return testkit.MockServer(mlUsername, mlPassword, mocks...)
}

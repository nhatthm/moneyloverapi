package moneyloverapi_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/nhatthm/moneyloverapi/pkg/wallet"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/moneyloverapi"
	"github.com/nhatthm/moneyloverapi/pkg/testkit"
	"github.com/nhatthm/moneyloverapi/pkg/transaction"
)

func TestClient_FindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	id1 := "ID1"
	id2 := "ID2"
	id3 := "ID3"

	transactionURL := "/transaction/list"

	withFindAllTransactionsInRange := func(result []transaction.Transaction) testkit.ServerOption {
		return testkit.WithFindAllTransactionsInRange("all", from, to, result)
	}

	testCases := []struct {
		scenario             string
		mockServer           testkit.ServerMocker
		expectedTransactions []transaction.Transaction
		expectedError        string
	}{
		{
			scenario: "server error",
			mockServer: mockServer(func(s *testkit.Server) {
				s.ExpectPost(transactionURL).
					ReturnCode(http.StatusInternalServerError)
			}),
			expectedError: "unexpected response status: 500 Internal Server Error",
		},
		{
			scenario: "success with an empty list",
			mockServer: mockServer(withFindAllTransactionsInRange(
				[]transaction.Transaction{},
			)),
			expectedTransactions: []transaction.Transaction{},
		},
		{
			scenario: "success with non empty list",
			mockServer: mockServer(withFindAllTransactionsInRange(
				[]transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
			)),
			expectedTransactions: []transaction.Transaction{{ID: id1}, {ID: id2}, {ID: id3}},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockServer(t)
			c := moneyloverapi.NewClient(
				moneyloverapi.WithBaseURL(s.URL()),
				moneyloverapi.WithCredentials(mlUsername, mlPassword),
			)

			result, err := c.FindAllTransactionsInRange(context.Background(), wallet.WalletAll, from, to)

			assert.Equal(t, tc.expectedTransactions, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

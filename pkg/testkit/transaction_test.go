package testkit

import (
	"testing"
	"time"

	"github.com/nhatthm/moneyloverapi/internal/api"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/moneyloverapi/pkg/transaction"
)

func TestWithFindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	id1 := "ID1"
	id2 := "ID2"
	id3 := "ID3"

	s := NewServer(t)

	WithFindAllTransactionsInRange("all", from, to,
		[]transaction.Transaction{
			{ID: id1},
			{ID: id2},
			{ID: id3},
		},
	)(s)

	expected := func() []*Request {
		s := NewServer(t)

		s.ExpectPost("/transaction/list").ReturnJSON(api.ListTransactionResponse{
			Error:  0,
			Msg:    "get_transaction_success",
			Action: "transaction_list",
			Data: api.ListTransactionResponseData{
				Daterange: &api.ListTransactionResponseDataDaterange{
					StartDate: from.Format("2006-01-02"),
					EndDate:   to.Format("2006-01-02"),
				},
				Transactions: []transaction.Transaction{
					{ID: id1},
					{ID: id2},
					{ID: id3},
				},
			},
		})

		return s.ExpectedRequests
	}()

	assert.Equal(t, len(expected), len(s.ExpectedRequests))

	for i, expected := range expected {
		actual := s.ExpectedRequests[i]

		expectedBody, err := expected.Handle(nil)
		assert.NoError(t, err)

		actualBody, err := actual.Handle(nil)
		assert.NoError(t, err)

		assert.Equal(t, expected.RequestURI, actual.RequestURI)
		assert.Equal(t, expectedBody, actualBody)
	}
}

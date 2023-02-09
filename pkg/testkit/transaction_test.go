package testkit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.nhat.io/httpmock"
	plannerMock "go.nhat.io/httpmock/mock/planner"
	"go.nhat.io/httpmock/planner"
	"go.nhat.io/matcher/v2"

	"github.com/nhatthm/moneyloverapi/internal/api"

	"github.com/nhatthm/moneyloverapi/pkg/transaction"
)

func TestWithFindAllTransactionsInRange(t *testing.T) {
	t.Parallel()

	type result struct {
		URI  string
		Body string
	}

	type expectation interface {
		httpmock.ExpectationHandler
		planner.Expectation
	}

	createTxsResponse := func(uri string, resp api.ListTransactionResponse) result {
		body, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		return result{
			URI:  uri,
			Body: string(body),
		}
	}

	from := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	to := time.Date(2020, 2, 2, 3, 4, 5, 0, time.UTC)
	id1 := "ID1"
	id2 := "ID2"
	id3 := "ID3"

	response := api.ListTransactionResponse{
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
	}

	actual := make([]expectation, 0)

	p := plannerMock.Mock(func(p *plannerMock.Planner) {
		p.On("Expect", mock.Anything).
			Run(func(args mock.Arguments) {
				actual = append(actual, args[0].(expectation))
			})
	})(t)

	s := NewServer(t).WithPlanner(p)

	WithFindAllTransactionsInRange("all", from, to,
		[]transaction.Transaction{
			{ID: id1},
			{ID: id2},
			{ID: id3},
		},
	)(s)

	expected := []result{
		createTxsResponse("/transaction/list", response),
	}

	assert.Equal(t, len(expected), len(actual))

	for i, actual := range actual {
		expected := expected[i]

		actualBody := handleExpectation(t, actual)

		assert.Equal(t, matcher.Exact(expected.URI), actual.URIMatcher())
		assert.Equal(t, expected.Body, actualBody)
	}
}

func handleExpectation(t *testing.T, e httpmock.ExpectationHandler) string {
	t.Helper()

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil) //nolint: errcheck

	err := e.Handle(rec, req, nil)
	require.NoError(t, err)

	return rec.Body.String()
}

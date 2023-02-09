package testkit

import (
	"time"

	"github.com/nhatthm/moneyloverapi/internal/api"
	"github.com/nhatthm/moneyloverapi/pkg/transaction"
)

func expectFindAllTransactionsInRange(s *Server, walletID string, from, to time.Time) Expectation {
	return s.ExpectPost("/transaction/list").
		WithBodyJSON(api.ListTransactionRequest{
			WalletID:  walletID,
			StartDate: from,
			EndDate:   to,
		})
}

// WithFindAllTransactionsInRange sets expectations for finding all transactions in a range.
func WithFindAllTransactionsInRange(
	walletID string,
	from time.Time,
	to time.Time,
	result []transaction.Transaction,
) ServerOption {
	return func(s *Server) {
		expectFindAllTransactionsInRange(s, walletID, from, to).
			ReturnJSON(api.ListTransactionResponse{
				Error:  0,
				Msg:    "get_transaction_success",
				Action: "transaction_list",
				Data: api.ListTransactionResponseData{
					Daterange: &api.ListTransactionResponseDataDaterange{
						StartDate: from.Format("2006-01-02"),
						EndDate:   to.Format("2006-01-02"),
					},
					Transactions: result,
				},
			})
	}
}

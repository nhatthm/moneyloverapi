package moneyloverapi

import (
	"context"
	"time"

	"github.com/nhatthm/moneyloverapi/pkg/wallet"

	"github.com/nhatthm/moneyloverapi/internal/api"
	"github.com/nhatthm/moneyloverapi/pkg/transaction"
)

var _ transaction.Finder = (*Client)(nil)

func (c *Client) findTransactions(ctx context.Context, req api.ListTransactionRequest) ([]transaction.Transaction, error) {
	res, err := c.api.PostTransactionList(ctx, api.PostTransactionListRequest{Body: &req})
	if err != nil {
		return nil, err
	}

	return res.ValueOK.Data.Transactions, nil
}

// FindAllTransactionsInRange finds all transactions in a time period.
func (c *Client) FindAllTransactionsInRange(ctx context.Context, wallet wallet.Wallet, from time.Time, to time.Time) ([]transaction.Transaction, error) {
	return c.findTransactions(ctx, api.ListTransactionRequest{
		WalletID:  wallet.ID,
		StartDate: from,
		EndDate:   to,
	})
}

package transaction

import (
	"context"
	"time"

	"github.com/nhatthm/moneyloverapi/pkg/wallet"
)

// Finder is a service to find transactions.
type Finder interface {
	// FindAllTransactionsInRange finds all transactions in a time period.
	FindAllTransactionsInRange(ctx context.Context, wallet wallet.Wallet, from time.Time, to time.Time) ([]Transaction, error)
}

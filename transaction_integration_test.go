// +build integration

package moneyloverapi

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nhatthm/moneyloverapi/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationFindAllTransactionsInRange tests transactions functionalities.
//
// In order to run this tests, there env vars should be set:
// - MONEYLOVER_USERNAME: The username to login to MoneyLover, it should be an email address.
// - MONEYLOVER_PASSWORD: The password to login to moneyloverapi.
// - MONEYLOVER_FROM: Transaction starting time. Default to 1 day ago (optional).
// - MONEYLOVER_TO: Transaction ending time. Default to now (optional).
func TestIntegrationFindAllTransactionsInRange(t *testing.T) {
	var from time.Time
	var to time.Time
	var err error

	if v, ok := os.LookupEnv("MONEYLOVER_TO"); ok {
		to, err = time.Parse(time.RFC3339, v)
		require.NoError(t, err)
	} else {
		to = time.Now()
	}

	if v, ok := os.LookupEnv("MONEYLOVER_FROM"); ok {
		from, err = time.Parse(time.RFC3339, v)
		require.NoError(t, err)
	} else {
		from = to.AddDate(0, 0, -1)
	}

	c := NewClient(
		WithBaseURL(BaseURL),
		WithTimeout(time.Minute),
	)

	done := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		transactions, err := c.FindAllTransactionsInRange(ctx, wallet.WalletAll, from, to)

		assert.NotEmpty(t, transactions)
		assert.NoError(t, err)

		close(done)
	}()

	select {
	case <-done:
		return

	case <-ctx.Done():
		t.Fatal("timeout while getting transactions")
	}
}

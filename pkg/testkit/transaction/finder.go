package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/nhatthm/moneyloverapi/pkg/transaction"
	"github.com/nhatthm/moneyloverapi/pkg/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// FinderMocker is Finder mocker.
type FinderMocker func(tb testing.TB) *Finder

// NoMockFinder is no mock Finder.
var NoMockFinder = MockFinder()

var _ transaction.Finder = (*Finder)(nil)

// Finder is a transaction.Finder.
type Finder struct {
	mock.Mock
}

// FindAllTransactionsInRange satisfies transaction.Finder.
func (f *Finder) FindAllTransactionsInRange(ctx context.Context, wallet wallet.Wallet, from time.Time, to time.Time) ([]transaction.Transaction, error) {
	result := f.Called(ctx, wallet, from, to)

	ret1 := result.Get(0)
	ret2 := result.Error(1)

	if ret1 == nil {
		return nil, ret2
	}

	return ret1.([]transaction.Transaction), ret2
}

// mockFinder mocks transaction.Finder interface.
func mockFinder(mocks ...func(f *Finder)) *Finder {
	f := &Finder{}

	for _, m := range mocks {
		m(f)
	}

	return f
}

// MockFinder creates Finder mock with cleanup to ensure all the expectations are met.
func MockFinder(mocks ...func(f *Finder)) FinderMocker {
	return func(tb testing.TB) *Finder {
		tb.Helper()

		f := mockFinder(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, f.Mock.AssertExpectations(tb))
		})

		return f
	}
}

package bank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper — reused across tests
func newTestAccount(t *testing.T, owner string, balance float64) *Account {
	t.Helper()
	acc, err := NewAccount(owner, balance)
	require.NoError(t, err)
	return acc
}

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name    string
		owner   string
		opening float64
		wantErr bool
	}{
		{"valid account", "Ahmed", 1000, false},
		{"zero opening balance", "Sara", 0, false},
		{"empty owner", "", 1000, true},
		{"negative opening balance", "Omar", -100, true},
		{"whitespace owner", "   ", 500, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := NewAccount(tt.owner, tt.opening)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, acc)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, acc)
			assert.Equal(t, tt.owner, acc.Owner())
			assert.Equal(t, tt.opening, acc.Balance())
		})
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		note        string
		wantBalance float64
		wantErr     bool
	}{
		{"valid deposit", 500, "salary", 1500, false},
		{"large deposit", 1_000_000, "jackpot", 1_001_000, false},
		{"zero amount", 0, "nothing", 0, true},
		{"negative amount", -100, "bad", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := newTestAccount(t, "Ahmed", 1000)
			err := acc.Deposit(tt.amount, tt.note)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidAmount)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantBalance, acc.Balance())
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		wantBalance float64
		wantErr     bool
		wantErrIs   error
	}{
		{"valid withdrawal", 300, 700, false, nil},
		{"exact balance", 1000, 0, false, nil},
		{"zero amount", 0, 0, true, ErrInvalidAmount},
		{"negative amount", -50, 0, true, ErrInvalidAmount},
		{"insufficient funds", 9999, 0, true, ErrInsufficientFunds},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := newTestAccount(t, "Ahmed", 1000)
			err := acc.Withdraw(tt.amount, "test")

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrIs != nil {
					assert.True(t, errors.Is(err, tt.wantErrIs),
						"expected errors.Is(%v, %v)", err, tt.wantErrIs)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantBalance, acc.Balance())
		})
	}
}

func TestTransferTo(t *testing.T) {
	t.Run("valid transfer", func(t *testing.T) {
		src := newTestAccount(t, "Ahmed", 1000)
		dst := newTestAccount(t, "Sara", 500)

		err := src.TransferTo(dst, 300)
		require.NoError(t, err)
		assert.Equal(t, 700.0, src.Balance())
		assert.Equal(t, 800.0, dst.Balance())
	})

	t.Run("insufficient funds", func(t *testing.T) {
		src := newTestAccount(t, "Ahmed", 100)
		dst := newTestAccount(t, "Sara", 0)

		err := src.TransferTo(dst, 500)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
		assert.Equal(t, 100.0, src.Balance()) // unchanged
		assert.Equal(t, 0.0, dst.Balance())   // unchanged
	})

	t.Run("nil target", func(t *testing.T) {
		src := newTestAccount(t, "Ahmed", 1000)
		err := src.TransferTo(nil, 100)
		require.Error(t, err)
	})

	t.Run("self transfer", func(t *testing.T) {
		acc := newTestAccount(t, "Ahmed", 1000)
		err := acc.TransferTo(acc, 100)
		require.Error(t, err)
	})

	t.Run("zero amount", func(t *testing.T) {
		src := newTestAccount(t, "Ahmed", 1000)
		dst := newTestAccount(t, "Sara", 0)
		err := src.TransferTo(dst, 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidAmount)
	})
}

func TestStatement(t *testing.T) {
	acc := newTestAccount(t, "Ahmed", 1000)
	acc.Deposit(500, "salary")
	acc.Withdraw(200, "rent")

	stmt := acc.Statement()
	assert.Contains(t, stmt, "Ahmed")
	assert.Contains(t, stmt, "DEPOSIT")
	assert.Contains(t, stmt, "WITHDRAWAL")
	assert.Contains(t, stmt, "1300") // final balance
}

func BenchmarkDeposit(b *testing.B) {
	acc, _ := NewAccount("bench", 1_000_000)
	b.ResetTimer()
	for b.Loop() {
		acc.Deposit(1, "bench")
	}
}

func BenchmarkWithdraw(b *testing.B) {
	acc, _ := NewAccount("bench", float64(b.N)*2)
	b.ResetTimer()
	for range b.N {
		acc.Withdraw(1, "bench")
	}
}

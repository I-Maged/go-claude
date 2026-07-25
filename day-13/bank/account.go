package bank

import (
	"day13/internal/validate"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient Funds")
	ErrInvalidAmount     = errors.New("Invalid Amount")
)

type TransactionKind string

const (
	KindDeposit    TransactionKind = "DEPOSIT"
	KindWithdrawal TransactionKind = "WITHDRAWAL"
	KindTransfer   TransactionKind = "TRANSFER"
)

type Transaction struct {
	Kind         TransactionKind
	Amount       float64
	At           time.Time
	Note         string
	BalanceAfter float64
}

func (t Transaction) String() string {
	sign := "+"
	if t.Kind == KindTransfer || t.Kind == KindWithdrawal {
		sign = "-"
	}

	return fmt.Sprintf("[%s] %s%8.2f  %-12s  balance:%8.2f  %s",
		t.At.Format("15:04:05"),
		sign,
		t.Amount,
		string(t.Kind),
		t.BalanceAfter,
		t.Note,
	)
}

type Account struct {
	owner        string
	balance      float64
	transactions []*Transaction
}

func NewAccount(owner string, opening float64) (*Account, error) {
	if err := validate.NotEmpty("owner", owner); err != nil {
		return nil, err
	}

	if err := validate.NonNegative("balance", opening); err != nil {
		return nil, err
	}

	a := &Account{owner: owner, balance: opening}
	if opening > 0 {
		a.record(KindDeposit, opening, "opening balance")
	}
	return a, nil
}

func (a *Account) record(kind TransactionKind, amount float64, note string) {
	a.transactions = append(a.transactions, &Transaction{
		Kind:         kind,
		Amount:       amount,
		At:           time.Now(),
		Note:         note,
		BalanceAfter: a.balance,
	})
}

func (a *Account) Deposit(amount float64, note string) error {
	if err := validate.Positive("deposit amount", amount); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAmount, err)
	}

	a.balance += amount
	a.record(KindDeposit, amount, note)

	return nil
}

func (a *Account) Withdraw(amount float64, note string) error {
	if err := validate.Positive("withdraw amount", amount); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAmount, err)
	}

	if amount > a.balance {
		return fmt.Errorf("%w: have %.2f, need %.2f", ErrInsufficientFunds, a.balance, amount)
	}

	a.record(KindWithdrawal, amount, note)
	a.balance -= amount

	return nil
}

func (a *Account) TransferTo(target *Account, amount float64) error {
	if target == nil {
		return errors.New("Target account cannot be nil")
	}

	if err := validate.NonNegative("transfer amount", amount); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAmount, err)
	}

	if amount > a.balance {
		return fmt.Errorf("%w: have %.2f need %.2f", ErrInsufficientFunds, a.balance, amount)
	}

	a.balance -= amount
	a.record(KindTransfer, amount, fmt.Sprintf("→ %s", target.owner))

	target.balance += amount
	target.record(KindDeposit, amount, fmt.Sprintf("← %s", a.owner))

	return nil
}

func (a *Account) Balance() float64 { return a.balance }
func (a *Account) Owner() string    { return a.owner }

func (a *Account) Statement() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "=== %s ===\n", a.owner)
	for _, t := range a.transactions {
		fmt.Fprintf(&sb, "  %s\n", t.String())
	}

	fmt.Fprintf(&sb, "  Balance: $%.2f\n", a.balance)
	return sb.String()
}

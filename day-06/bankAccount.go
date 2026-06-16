package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
)

type Transaction struct {
	Type      TransactionType
	Amount    float64
	Timestamp time.Time
	Note      string
}

func (t Transaction) String() string {
	sign := "+"
	if t.Type == Withdrawal || t.Type == Transfer {
		sign = "-"
	}
	return fmt.Sprintf("[%s] %s%.2f  %-12s %s",
		t.Timestamp.Format("15:04:05"),
		sign,
		t.Amount,
		string(t.Type),
		t.Note,
	)
}

type BankAccount struct {
	owner        string
	balance      float64
	transactions []Transaction
}

// Constructors — Go's convention
func NewBankAccount(owner string, initialDeposit float64) (*BankAccount, error) {
	if strings.TrimSpace(owner) == "" {
		return nil, errors.New("owner name cannot be empty")
	}

	if initialDeposit < 0 {
		return nil, errors.New("initial deposit cannot be negative")
	}

	acc := &BankAccount{owner: owner}
	if initialDeposit > 0 {
		acc.balance = initialDeposit
		acc.transactions = append(acc.transactions, Transaction{
			Type:      Deposit,
			Amount:    initialDeposit,
			Timestamp: time.Now(),
			Note:      "initial deposit",
		})
	}

	return acc, nil
}

func (b *BankAccount) Deposit(amount float64, note string) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}

	b.balance += amount
	b.transactions = append(b.transactions, Transaction{
		Type:      Deposit,
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      note,
	})

	return nil
}

func (b *BankAccount) Withdraw(amount float64, note string) error {
	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}

	if amount > b.balance {
		return fmt.Errorf("insufficient funds: have %.2f, need %.2f", b.balance, amount)
	}

	b.balance -= amount
	b.transactions = append(b.transactions, Transaction{
		Type:      Withdrawal,
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      note,
	})

	return nil
}

func (b *BankAccount) TransferTo(target *BankAccount, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	if amount > b.balance {
		return fmt.Errorf("insufficient funds for transfer: have %.2f, need %.2f",
			b.balance, amount)
	}

	b.balance -= amount
	senderNote := fmt.Sprintf("transfer to %s", target.owner)
	b.transactions = append(b.transactions, Transaction{
		Type:      Transfer,
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      senderNote,
	})

	target.balance += amount
	targetNote := fmt.Sprintf("Transfer from %s", b.owner)
	target.transactions = append(target.transactions, Transaction{
		Type:      Transfer,
		Amount:    amount,
		Timestamp: time.Now(),
		Note:      targetNote,
	})

	return nil
}

func (b *BankAccount) Balance() float64 { return b.balance }

func (b *BankAccount) Owner() string { return b.owner }

func (a *BankAccount) Statement() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "=== Statement for %s ===\n", a.owner)

	sb.WriteString(strings.Repeat("-", 55))
	sb.WriteString("\n")

	if len(a.transactions) == 0 {
		sb.WriteString("No transactions.\n")
	} else {
		for _, t := range a.transactions {
			sb.WriteString(t.String())
			sb.WriteString("\n")
		}
	}
	sb.WriteString(strings.Repeat("-", 55))
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "Current balance: $%.2f\n", a.balance)

	return sb.String()
}

type Bank struct {
	Name     string
	accounts map[string]*BankAccount
}

func NewBank(name string) *Bank {
	return &Bank{
		Name:     name,
		accounts: make(map[string]*BankAccount),
	}
}

func (k *Bank) OpenAccount(owner string, initialDeposit float64) (*BankAccount, error) {
	if _, exists := k.accounts[owner]; exists {
		return nil, fmt.Errorf("account for %s already exists", owner)
	}

	acc, err := NewBankAccount(owner, initialDeposit)
	if err != nil {
		return nil, err
	}
	k.accounts[owner] = acc

	return acc, nil
}

func (k *Bank) TotalAssets() float64 {
	total := 0.0
	for _, acc := range k.accounts {
		total += acc.Balance()
	}

	return total
}

func (k *Bank) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s Bank Summary ===\n", k.Name))
	for _, acc := range k.accounts {
		sb.WriteString(fmt.Sprintf("  %-15s $%.2f\n", acc.Owner(), acc.Balance()))
	}
	sb.WriteString(fmt.Sprintf("  Total assets: $%.2f\n", k.TotalAssets()))
	return sb.String()
}

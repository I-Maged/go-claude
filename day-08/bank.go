package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// --- Types ---

type TransactionType string

const (
	TypeDeposit  TransactionType = "DEPOSIT"
	TypeWithdraw TransactionType = "WITHDRAW"
	TypeTransfer TransactionType = "TRANSFER"
)

type Transaction struct {
	Type         TransactionType
	Amount       float64
	At           time.Time
	Note         string
	BalanceAfter float64
}

func (t Transaction) String() string {
	sign := "+"
	if t.Type == TypeWithdraw || t.Type == TypeTransfer {
		sign = "-"
	}

	return fmt.Sprintf("[%s] %s%8.2f  %-12s  balance: %8.2f  | %s",
		t.At.Format("15:04:05.000"),
		sign, t.Amount,
		string(t.Type),
		t.BalanceAfter,
		t.Note,
	)
}

// --- BankAccount ---

type BankAccount struct {
	owner        string
	balance      float64
	transactions []*Transaction
}

func NewBankAccount(owner string, initialDeposit float64) (*BankAccount, error) {
	if owner == "" {
		return nil, errors.New("owner cannot be empty")
	}
	if initialDeposit < 0 {
		return nil, errors.New("initial deposit cannot be negative")
	}

	ba := &BankAccount{owner: owner, balance: initialDeposit}

	if initialDeposit > 0 {
		ba.record(TypeDeposit, initialDeposit, "initial deposit")
	}

	return ba, nil
}

func (ba *BankAccount) record(ttype TransactionType, amount float64, note string) {
	ba.transactions = append(ba.transactions, &Transaction{
		Type:         ttype,
		Amount:       amount,
		At:           time.Now(),
		Note:         note,
		BalanceAfter: ba.balance,
	})
}

func (ba *BankAccount) Deposit(amount float64, note string) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive, got %.2f", amount)
	}

	ba.balance += amount
	ba.record(TypeDeposit, amount, note)

	return nil
}

func (ba *BankAccount) Withdraw(amount float64, note string) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw amount cannot be negative, got %.2f", amount)
	}

	if ba.balance < amount {
		return fmt.Errorf("insufficient funds. balance:%.2f, requested:%.2f", ba.balance, amount)
	}

	ba.balance -= amount
	ba.record(TypeWithdraw, amount, note)

	return nil
}

func (ba *BankAccount) TransferTo(target *BankAccount, amount float64) error {
	if target == nil {
		return errors.New("target account doesnot exist")
	}
	if ba == target {
		return errors.New("cannot transfer to the same account")
	}
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive. got %.2f", amount)
	}
	if ba.balance < amount {
		return fmt.Errorf("insufficient funds. balance:%.2f, requested:%.2f", ba.balance, amount)
	}

	ba.balance -= amount
	ba.record(TypeTransfer, amount, fmt.Sprintf("→ %s", target.owner))

	target.balance += amount
	target.record(TypeTransfer, amount, fmt.Sprintf("← %s", ba.owner))

	return nil
}

// Helper functions — Go has no built-in way to get pointer to a literal
func (ba *BankAccount) Balance() float64 { return ba.balance }
func (ba *BankAccount) Owner() string    { return ba.owner }

func (ba *BankAccount) LastTransaction() *Transaction {
	if len(ba.transactions) == 0 {
		return nil
	}
	return ba.transactions[len(ba.transactions)-1]
}

func (ba *BankAccount) ModifyLastNote(newNote string) error {
	t := ba.LastTransaction()
	if t == nil {
		return errors.New("no transactions to modify")
	}
	t.Note = newNote
	return nil
}

func (ba *BankAccount) Statement() {
	fmt.Printf("\n=== Statement: %s ===\n", ba.owner)
	fmt.Println(strings.Repeat("─", 72))
	if len(ba.transactions) == 0 {
		fmt.Println("  No transactions.")
	}
	for _, t := range ba.transactions {
		fmt.Println(" ", t)
	}
	fmt.Println(strings.Repeat("─", 72))
	fmt.Printf("  Current balance: $%.2f\n", ba.balance)
}

// --- Pointer benchmark ---
// largeCopy is a big struct to show copy cost
type LargeStruct struct {
	data [1000]float64
}

func processByValue(s LargeStruct) float64 {
	total := 0.0
	for _, v := range s.data {
		total += v
	}
	return total
}

func processByPointer(s *LargeStruct) float64 {
	total := 0.0
	for _, v := range s.data {
		total += v
	}
	return total
}

func benchmarkCopying() {
	s := LargeStruct{}
	for i := range s.data {
		s.data[i] = float64(i)
	}

	const iterations = 100_000

	start := time.Now()
	total := 0.0
	for range iterations {
		total += processByValue(s) // copies 8000 bytes each call
	}
	valueTime := time.Since(start)

	start = time.Now()
	total = 0.0
	for range iterations {
		total += processByPointer(&s) // copies 8 bytes (pointer) each call
	}
	pointerTime := time.Since(start)

	fmt.Printf("\n=== Benchmark: value vs pointer receiver ===\n")
	fmt.Printf("Struct size:     %d bytes\n", 1000*8)
	fmt.Printf("Iterations:      %d\n", iterations)
	fmt.Printf("By value:        %v\n", valueTime)
	fmt.Printf("By pointer:      %v\n", pointerTime)
	speedup := float64(valueTime) / float64(pointerTime)
	fmt.Printf("Pointer speedup: %.1fx faster\n", speedup)
}

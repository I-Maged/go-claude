package main

import "fmt"

func main() {
	ahmed, err := NewBankAccount("Ahmed", 5000)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sara, err := NewBankAccount("Sara", 2000)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ahmed.Deposit(1500, "salary")
	ahmed.Withdraw(300, "groceries")
	ahmed.TransferTo(sara, 1000)
	sara.Deposit(500, "freelance")
	sara.Withdraw(200, "utilities")

	ahmed.ModifyLastNote("transfer to Sara — rent share")

	ahmed.Statement()
	sara.Statement()

	fmt.Println("\n=== Nil pointer safety ===")
	var nilAccount *BankAccount
	if err := ahmed.TransferTo(nilAccount, 100); err != nil {
		fmt.Println("Transfer to nil:", err)
	}
	if err := ahmed.TransferTo(ahmed, 100); err != nil {
		fmt.Println("Self transfer:", err)
	}

	fmt.Println("\n=== Last transaction pointer ===")
	last := ahmed.LastTransaction()
	if last != nil {
		fmt.Printf("Last transaction: %s\n", last.Note)
		fmt.Printf("Amount: $%.2f\n", last.Amount)
	}

	benchmarkCopying()
}

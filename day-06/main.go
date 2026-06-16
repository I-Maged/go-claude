package main

import "fmt"

func main() {
	bank := NewBank("Go National Bank")

	ahmed, err := bank.OpenAccount("Ahmed", 5000)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sara, err := bank.OpenAccount("Sara", 3000)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Transactions
	ahmed.Deposit(1500, "salary")
	ahmed.Withdraw(200, "groceries")
	ahmed.Withdraw(800, "rent")
	ahmed.TransferTo(sara, 500)

	sara.Deposit(2000, "freelance payment")
	sara.Withdraw(150, "utilities")

	// Try error cases
	fmt.Println("--- Error cases ---")
	if err := ahmed.Withdraw(99999, "buy a yacht"); err != nil {
		fmt.Println("Error:", err)
	}
	if err := ahmed.Deposit(-100, "negative"); err != nil {
		fmt.Println("Error:", err)
	}
	if _, err := bank.OpenAccount("Ahmed", 0); err != nil {
		fmt.Println("Error:", err)
	}

	// Print statements
	fmt.Println()
	fmt.Println(ahmed.Statement())
	fmt.Println(sara.Statement())
	fmt.Println(bank.Summary())
}

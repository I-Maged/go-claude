package main

import (
	"day13/bank"
	"day13/cache"
	"day13/jobqueue"
	"errors"
	"fmt"
	"time"
)

func demoBank() {
	fmt.Println("=== Bank ===")
	ahmed, err := bank.NewAccount("Ahmed", 5000)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	sara, _ := bank.NewAccount("Sara", 2000)

	ahmed.Deposit(1500, "salary")
	ahmed.Withdraw(300, "groceries")
	ahmed.TransferTo(sara, 1000)

	fmt.Print(ahmed.Statement())
	fmt.Print(sara.Statement())

	// errors.Is works across package boundary
	err = ahmed.Withdraw(99999, "yacht")
	if errors.Is(err, bank.ErrInsufficientFunds) {
		fmt.Println("caught:", err)
	}
}

func demoCache() {
	fmt.Println("\n=== Cache ===")
	store := cache.NewStore(500 * time.Millisecond)
	defer store.Close()

	store.Set("session:1", "user-ahmed", 1*time.Second)
	store.Set("config:env", "production", 10*time.Second)

	if v, ok := store.Get("session:1"); ok {
		fmt.Println("session:1 =", v)
	}
	if v, ok := store.Get("config:env"); ok {
		fmt.Println("config:env =", v)
	}

	fmt.Println("cache size:", store.Len())
}

func demoJobQueue() {
	fmt.Println("\n=== Job Queue ===")
	q := jobqueue.New(3, 10)

	for i := 1; i <= 10; i++ {
		id := i
		err := q.Submit(jobqueue.Job{
			ID:      id,
			Payload: fmt.Sprintf("task-%d", id),
			Execute: func() error {
				time.Sleep(50 * time.Millisecond)
				return nil
			},
		})
		if err != nil {
			fmt.Println("rejected:", err)
		}
	}

	time.Sleep(500 * time.Millisecond)
	q.Shutdown()
	q.Stats()
}

func main() {
	demoBank()
	demoCache()
	demoJobQueue()
}

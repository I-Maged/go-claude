package main

import (
	"fmt"
	"math/rand"
	"time"
)

func GuessingGame() {
	fmt.Println("=====Guessing Games=====")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	secret := rng.Intn(10) + 1
	attempts := 0
	maxAttempts := 4

	fmt.Println("I'm thinking of number between 1 & 10")
	fmt.Printf("You have %d attempts\n", maxAttempts)

	for attempts < maxAttempts {
		attempts++
		fmt.Printf("Attempt %d/%d - Enter your guess:\n", attempts, maxAttempts)

		var guess int
		fmt.Scan(&guess)

		switch {
		case guess < 1 || guess > 10:
			fmt.Println("Out of range. Enter a number between 1 & 10")
			attempts--
		case guess > secret:
			fmt.Println("Too high")
		case guess < secret:
			fmt.Println("Too low")
		case guess == secret:
			fmt.Printf("Correct! You got it in %d attempt(s)!\n", attempts)
		}
	}
	fmt.Printf("Out of attempts. The number was %d.\n", secret)
}

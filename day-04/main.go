package main

import (
	"fmt"
)

func main() {
	min, max := minMax([]int{7, 5, 3, 9, 1})
	fmt.Printf("min=%d, max=%d\n", min, max)

	fmt.Println("=== Go Calculator ===")
	fmt.Println()

	calculate("10 + 5", 10, 5, add)
	calculate("10 - 5", 10, 5, subtract)
	calculate("10 * 5", 10, 5, multiply)
	calculate("10 / 5", 10, 5, divide)
	calculate("10 / 0", 10, 0, divide)
	calculate("2 ^ 8", 2, 8, power)

	fmt.Println("Wrapping single-arg functions:")
	result, err := sqrt(144)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("√144 = %.4f\n", result)
	}
}

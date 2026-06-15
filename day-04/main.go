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
	calculate("10 / 0", 10, 0, divide) // triggers error
	calculate("2 ^ 8", 2, 8, power)

	fmt.Println("Wrapping single-arg functions:")
	result, err := squareRoot(144)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("√144 = %.4f\n", result)
	}

	result, err = squareRoot(-9)
	if err != nil {
		fmt.Println("Error:", err) // triggers error
	} else {
		fmt.Printf("√-9 = %.4f\n", result)
	}

	fmt.Println("\n=== Closure: multiplier factory ===")
	triple := func(x float64) float64 { return x * 3 }
	half := func(x float64) float64 { return x / 2 }

	for _, n := range []float64{4, 8, 12} {
		fmt.Printf("triple(%.0f)=%.0f  half(%.0f)=%.1f\n", n, triple(n), n, half(n))
	}
}

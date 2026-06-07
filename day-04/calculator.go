package main

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrDivisionByZero = errors.New("Can't divide by 0")
	ErrNegativeSqrt   = errors.New("Can't take square root of a negative number")
)

func add(a, b float64) (float64, error) {
	return a + b, nil
}

func subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func multiply(a, b float64) (float64, error) {
	return a * b, nil
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

func sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return math.Sqrt(a), nil
}

func power(a, b float64) (float64, error) {
	return math.Pow(a, b), nil
}

func calculate(label string, a, b float64, op func(float64, float64) (float64, error)) {
	defer fmt.Println("---")

	result, err := op(a, b)
	if err != nil {
		fmt.Printf("%-20s ERROR: %v\n", label, err)
		return
	}
	fmt.Printf("%-20s = %.4f\n", label, result)
}

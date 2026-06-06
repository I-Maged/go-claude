package main

import (
	"fmt"
)

func main() {
	if result, err := divide(6, 3); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Division result = %f\n", result)
	}

	min, max := minMax([]int{7, 5, 3, 9, 1})
	fmt.Printf("min=%d, max=%d\n", min, max)
}

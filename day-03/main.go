package main

import (
	"fmt"
	"runtime"
	"strconv"
)

func main() {
	if n, err := strconv.Atoi("123"); err == nil {
		fmt.Printf("Parsed number: %d\n", n)
	} else {
		fmt.Printf("Error: %v\n", err)
	}

outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer
			}
			fmt.Printf("i==%d, & j==%d\n", 1, j)
		}
	}

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("Mac")
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Printf("Other: %s\n", os)
	}

	names := []string{"maged", "Ibrahim"}
	for _, name := range names {
		fmt.Println(name)
	}

	fmt.Printf("%s\n", fizzBuzz(7))

	GuessingGame()
}

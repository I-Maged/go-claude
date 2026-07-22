package main

import (
	"fmt"
	"sync"
)

func unsafeCounter(n int) int {
	counter := 0
	var wg sync.WaitGroup

	for range n {
		wg.Add(1)
		go func() {
			counter++
			wg.Done()
		}()
	}

	wg.Wait()

	return counter
}

func safeCounter(n int) int {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			counter++
		}()
	}

	wg.Wait()
	return counter
}

func main() {
	const n = 1000

	fmt.Println("=== Without Mutex ===")
	for i := range 5 {
		result := unsafeCounter(n)
		fmt.Printf("run %d: got %d (expected %d, off by %d)\n",
			i+1, result, n, n-result)
	}

	fmt.Println("\n=== With Mutex ===")
	for i := range 5 {
		result := safeCounter(n)
		fmt.Printf("run %d: got %d (expected %d, off by %d)\n",
			i+1, result, n, n-result)
	}
}

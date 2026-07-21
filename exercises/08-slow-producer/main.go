package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var ErrTimeout = errors.New("operation timed out")

func simulateWork(min, max int) int {
	duration := time.Duration(min+rand.Intn(max-min)) * time.Millisecond
	time.Sleep(duration)
	return rand.Intn(1000)
}

func fetchWithTimeout(workFn func() int, limit time.Duration) (int, error) {
	// buffered: goroutine can exit even after timeout
	resultCh := make(chan int, 1)

	go func() {
		resultCh <- workFn()
	}()

	timer := time.NewTimer(limit)
	defer timer.Stop()

	select {
	case result := <-resultCh:
		return result, nil
	case <-timer.C:
		return 0, ErrTimeout
	}
}

func main() {
	const timeout = 500 * time.Millisecond
	fast, slow := 0, 0

	for i := 1; i <= 10; i++ {
		result, err := fetchWithTimeout(func() int {
			return simulateWork(100, 900)
		}, timeout)

		if err != nil {
			fmt.Printf("trial %2d: TIMEOUT\n", i)
			slow++
		} else {
			fmt.Printf("trial %2d: got %d\n", i, result)
			fast++
		}
	}

	fmt.Printf("\n%d succeeded, %d timed out\n", fast, slow)
}

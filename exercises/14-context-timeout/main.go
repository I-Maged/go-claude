package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func simulateHTTPCall(duration time.Duration) string {
	time.Sleep(duration)
	return `{"status": "ok", "data": "some response"}`
}

func fetchWithContext(ctx context.Context, callDuration time.Duration) (string, error) {

	ch := make(chan string, 1)
	// TODO: launch simulateHTTPCall in a goroutine,
	//       send its return value into the result channel
	go func() {
		result := simulateHTTPCall(callDuration)
		ch <- result
	}()

	select {
	case result := <-ch:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	// Case 1: call takes 500ms, timeout is 1s — should SUCCEED
	fmt.Println("=== Case 1: fast call (500ms) with 1s timeout ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()

	result, err := fetchWithContext(ctx1, 500*time.Millisecond)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("success:", result)
	}

	// Case 2: call takes 2s, timeout is 1s — should TIMEOUT
	fmt.Println("\n=== Case 2: slow call (2s) with 1s timeout ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	result, err = fetchWithContext(ctx2, 2*time.Second)
	if err != nil {
		fmt.Println("error:", err)
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Request timed out")
		}
	} else {
		fmt.Println("success:", result)
	}

	// Case 3: context already cancelled before call starts
	fmt.Println("\n=== Case 3: pre-cancelled context ===")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 1*time.Second)
	cancel3()

	result, err = fetchWithContext(ctx3, 500*time.Millisecond)
	if err != nil {
		fmt.Println("error:", err)
		fmt.Println("is DeadlineExceeded:", errors.Is(err, context.DeadlineExceeded))
		fmt.Println("is Canceled:", errors.Is(err, context.Canceled))
	} else {
		fmt.Println("success:", result)
	}
}

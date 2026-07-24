package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d shutting down: %v\n", id, ctx.Err())
			return

		case t := <-ticker.C:
			fmt.Printf("worker %d doing work at %v\n", id, t.Format("15:04:05.000"))
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("\nmain: cancelling context...")

	cancel()
	wg.Wait()

	fmt.Println("main: all workers stopped")
}

package main

import (
	"fmt"
	"sync"
	"time"
)

func OffByOne() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			fmt.Println("worker", n, "done")
		}(i)
	}

	wg.Wait()
	fmt.Println("all workers finished")
}

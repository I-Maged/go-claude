package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func processJob() {
	time.Sleep(500 * time.Millisecond)
}

func main() {
	const totalJobs = 100
	const maxConcurrent = 5
	var active atomic.Int64
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	for i := 1; i <= totalJobs; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() {
				<-sem
			}()

			current := active.Add(1)
			fmt.Printf("job %3d started (active: %d)\n", id, current)

			processJob()

			current = active.Add(-1)
			fmt.Printf("job %3d finished (active: %d)\n", id, current)
		}(i)
	}

	wg.Wait()
	fmt.Println("all jobs done")
}

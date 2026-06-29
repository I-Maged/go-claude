package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func makeSource(id, count int) <-chan int {
	src := make(chan int)

	go func() {
		defer close(src)
		for range count {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			src <- rand.Intn(100)
		}
		fmt.Printf("source %d closed\n", id)
	}()

	return src
}

func fanIn(sources ...<-chan int) <-chan int {
	merged := make(chan int)
	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for c := range ch {
				merged <- c
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	// rand.Seed(time.Now().UnixNano())
	// rng := rand.New(rand.NewSource(5))

	s1 := makeSource(1, 5)
	s2 := makeSource(2, 5)
	s3 := makeSource(3, 5)

	merged := fanIn(s1, s2, s3)

	total := 0
	count := 0

	for v := range merged {
		fmt.Printf("received: %d\n", v)
		count++
		total += v
	}

	fmt.Printf("received %d values, total sum: %d\n", count, total)
}

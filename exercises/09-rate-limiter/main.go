package main

import (
	"fmt"
	"time"
)

type Request struct {
	ID        int
	ArrivedAt time.Time
}

func rateLimiter(in <-chan Request, rate time.Duration) <-chan Request {
	out := make(chan Request)

	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()

		// Let the first request through immediately
		first := true
		for req := range in {
			if first {
				out <- req
				first = false
				continue
			}
			<-ticker.C
			out <- req
		}

		close(out)
	}()

	return out
}

func main() {
	const numRequests = 20
	const rateLimit = 200 * time.Millisecond

	in := make(chan Request, numRequests)
	for i := 1; i <= numRequests; i++ {
		in <- Request{ID: i, ArrivedAt: time.Now()}
	}
	close(in)

	limited := rateLimiter(in, rateLimit)

	start := time.Now()
	for req := range limited {
		elapsed := time.Since(start)
		fmt.Printf(
			"req %2d | arrived: +%3dms | allowed: +%4dms\n",
			req.ID,
			req.ArrivedAt.Sub(start).Milliseconds(),
			elapsed.Milliseconds(),
		)
	}

	fmt.Printf("\ntotal time: %v\n", time.Since(start).Round(time.Millisecond))
	fmt.Printf("expected:   ~%v\n", time.Duration(numRequests-1)*rateLimit)
}

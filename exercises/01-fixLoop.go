package main

import (
	"fmt"
	"sync"
)

func FixLoop() {
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}

	wg.Wait()
}

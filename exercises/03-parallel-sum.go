package main

import (
	"fmt"
	"sync"
)

func parallelSum(nums []int, numWorkers int) int {
	n := len(nums)
	chunkSize := n / numWorkers

	results := make([]int, numWorkers) // TODO: what goes in here, and why this size?
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		// TODO: compute start and end indices for this worker's chunk
		// careful: the LAST worker must go all the way to n, not just (w+1)*chunkSize,
		// in case n doesn't divide evenly by numWorkers

		start := chunkSize * w
		end := 0
		if w == numWorkers-1 {
			end = n
		} else {
			end = chunkSize * (w + 1)
		}

		wg.Add(1)
		go func(workerIndex, start, end int) {
			defer wg.Done()
			sum := 0
			// TODO: sum nums[start:end]
			// TODO: write the result into this worker's own slot in `results`
			for _, num := range nums[start:end] {
				sum += num
			}
			results[workerIndex] = sum
		}(w, start, end)
	}

	wg.Wait()

	fmt.Println("results:", results)

	// TODO: combine all partial sums into one final total
	total := 0
	for _, result := range results {
		total += result
	}
	return total
}

func TestParallelSum() {
	nums := make([]int, 1_000_000)
	for i := range nums {
		nums[i] = i + 1 // 1, 2, 3, ..., 1000000
	}

	result := parallelSum(nums, 4)
	fmt.Println("Parallel sum:", result)

	// Sanity check against the known formula for sum of 1..n
	expected := 1_000_000 * (1_000_000 + 1) / 2
	fmt.Println("Expected:    ", expected)
	fmt.Println("Match:", result == expected)
}

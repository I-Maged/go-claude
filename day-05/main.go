package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	text := `Go is an open source programming language that makes it easy
    to build simple reliable and efficient software. Go is expressive
    concise clean and efficient. Its concurrency mechanisms make it easy
    to write programs that get the most out of multicore and networked
    machines while its novel type system enables flexible and modular
    program construction. Go compiles quickly to machine code yet it
    has the convenience of garbage collection and the power of
    run-time reflection. It is a fast statically typed compiled
    language that feels like a dynamically typed interpreted language.
    Go is fun. Go is fast. Go is easy.`

	freq := wordFrequency(text)

	fmt.Println("=== Word Frequency Counter ===")
	fmt.Printf("Total unique words: %d\n\n", len(freq))

	// Top 10 most frequent
	top := topN(freq, 10)
	fmt.Println("Top 10 words:")
	printFreqTable(freq, top)

	fmt.Println()

	// Words that appear exactly once
	hapax := make([]string, 0)
	for word, count := range freq {
		if count == 1 {
			hapax = append(hapax, word)
		}
	}
	sort.Strings(hapax)
	fmt.Printf("Words appearing once (%d total):\n", len(hapax))
	fmt.Println(strings.Join(hapax, ", "))

	// copy() demo — safe independent copy
	fmt.Println("\n=== copy() demo ===")
	original := []int{1, 2, 3, 4, 5}
	cloned := make([]int, len(original))
	copy(cloned, original)
	cloned[0] = 999
	fmt.Printf("original: %v\n", original) // unchanged
	fmt.Printf("cloned:   %v\n", cloned)

	// Slice operations demo
	fmt.Println("\n=== Slice operations ===")
	nums := []int{10, 20, 30, 40, 50}
	fmt.Printf("nums:        %v\n", nums)
	fmt.Printf("nums[1:3]:   %v\n", nums[1:3])
	fmt.Printf("nums[:2]:    %v\n", nums[:2])
	fmt.Printf("nums[3:]:    %v\n", nums[3:])

	// Remove element at index 2 (30) — idiomatic Go
	i := 2
	nums = append(nums[:i], nums[i+1:]...)
	fmt.Printf("after remove index 2: %v\n", nums)
}

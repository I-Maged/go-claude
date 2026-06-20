package main

import "fmt"

func main() {
	// pipeline.TestPipeline()
	result := generate(10)
	evens := filterEvens(result)
	squared := square(evens)

	for v := range squared {
		fmt.Println(v)
	}
}

package main

func generate(n int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()

	return out
}

func filterEvens(ch <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range ch {
			if v%2 == 0 {
				out <- v
			}
		}
	}()

	return out
}

func square(ch <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range ch {
			out <- v * v
		}
	}()

	return out
}

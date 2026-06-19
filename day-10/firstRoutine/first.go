package firstroutine

import (
	"fmt"
	"sync"
	"time"
)

func SayHello() {
	fmt.Println("Hello Go Routine")
}

func TestRoutine() {
	go SayHello()

	fmt.Println("Hello from Test")
	time.Sleep(100 * time.Millisecond)
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %v starting\n", id)
	fmt.Printf("Worker %v is done\n", id)
}

func TestWorker() {
	fmt.Println("=== Workers starting  ===")
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("=== All workers finished  ===")
}

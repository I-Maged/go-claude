package main

import (
	"fmt"
	"sync"
)

func ping(pingCh chan<- string, pongCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for range 5 {
		pingCh <- "Ping"
		msg := <-pongCh
		fmt.Printf("Received: %s\n", msg)
	}
}

func pong(pingCh <-chan string, pongCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for range 5 {
		msg := <-pingCh
		fmt.Printf("Received: %s\n", msg)
		pongCh <- "Pong"
	}
}

func main() {
	pingCh := make(chan string)
	pongCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go ping(pingCh, pongCh, &wg)
	go pong(pingCh, pongCh, &wg)

	wg.Wait()

	close(pingCh)
	close(pongCh)

	fmt.Println("done")
}

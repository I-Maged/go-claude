package main

import (
	"fmt"
	"sync"
)

func Ping(pingCh chan<- string, pongCh <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for range 5 {
		pingCh <- "Ping"
		msg := <-pongCh
		fmt.Printf("received: %s\n", msg)
	}
}

func Pong(pingCh <-chan string, pongCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for range 5 {
		msg := <-pingCh
		fmt.Printf("received: %s\n", msg)
		pongCh <- "Pong"
	}
}

func main() {
	pingCh := make(chan string)
	pongCh := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)

	go Ping(pingCh, pongCh, &wg)
	go Pong(pingCh, pongCh, &wg)

	wg.Wait()
	fmt.Println("done")
}

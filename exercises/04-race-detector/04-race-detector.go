package main

import (
	"fmt"
	"maps"
	"sync"
)

var pages = []string{"/home", "/about", "/contact", "/blog"}

type VisitCounter struct {
	mu         sync.Mutex
	pageVisits map[string]int
}

func (vc *VisitCounter) Record(page string) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.pageVisits[page]++
}

func (vc *VisitCounter) Print() {
	vc.mu.Lock()

	snapshot := make(map[string]int, len(vc.pageVisits))
	maps.Copy(snapshot, vc.pageVisits)

	vc.mu.Unlock()

	fmt.Println("Visit Count:")

	for page, count := range snapshot {
		fmt.Printf("  %-12s %d\n", page, count)
	}
}

func main() {
	var wg sync.WaitGroup

	vc := &VisitCounter{
		pageVisits: map[string]int{
			"/home":    0,
			"/about":   0,
			"/contact": 0,
			"/blog":    0,
		},
	}

	for i := range 50 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			page := pages[id%len(pages)]
			vc.Record(page)
		}(i)
	}

	wg.Wait()
	vc.Print()
}

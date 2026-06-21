package main

import (
	"fmt"
	"sync"
)

// type Cache struct {
// 	mu   sync.RWMutex
// 	data map[string]string
// }

// func NewCache() *Cache {
// 	return &Cache{data: make(map[string]string)}
// }

// func (c *Cache) Get(key string) (string, bool) {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()
// 	v, ok := c.data[key]
// 	return v, ok
// }

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.data[key] = value
}

func main() {
	cache := NewCache()
	cache.Set("name", "Maged")

	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v, _ := cache.Get("name")
			fmt.Printf("reader %d got: %s\n", id, v)
		}(i)
	}

	wg.Wait()
}

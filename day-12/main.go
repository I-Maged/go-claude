package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	value     string
	expiresAt time.Time
}

func (e *cacheEntry) isExpired() bool {
	return time.Now().After(e.expiresAt)
}

type Cache struct {
	mu      sync.RWMutex
	data    map[string]cacheEntry
	stopped sync.Once
	cancel  context.CancelFunc
}

func NewCache(cleanUpInterval time.Duration) *Cache {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Cache{
		data:   make(map[string]cacheEntry),
		cancel: cancel,
	}
	go c.startCleanupLoop(ctx, cleanUpInterval)
	return c
}

func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.data[key]

	if !exists {
		return "", false
	}

	if entry.isExpired() {
		return "", false
	}

	return entry.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	delete(c.data, key)
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *Cache) startCleanupLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("cache: cleanup loop stopping")
			return
		case <-ticker.C:
			c.cleanupExpired()
		}
	}
}

func (c *Cache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	removed := 0
	for key, entry := range c.data {
		if entry.isExpired() {
			delete(c.data, key)
		}
	}
	if removed > 0 {
		fmt.Printf("cache: cleaned up %d expired entries\n", removed)
	}
}

func (c *Cache) Close() {
	c.stopped.Do(func() {
		c.cancel()
	})
}

func main() {
	cache := NewCache(500 * time.Millisecond)
	defer cache.Close()

	cache.Set("session:abc", "user-42", 1*time.Second)
	cache.Set("session:xyz", "user-99", 3*time.Second)
	cache.Set("config:env", "production", 10*time.Second)

	fmt.Println("=== Concurrent reads and writes ===")
	var wg sync.WaitGroup

	for i := range 20 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%4 == 0 {
				key := fmt.Sprintf("dynamic:%d", id)
				cache.Set(key, fmt.Sprintf("value-%d", id), 2*time.Second)
			} else {
				if v, ok := cache.Get("config:env"); ok {
					_ = v
				}
			}
		}(i)
	}
	wg.Wait()

	fmt.Printf("Cache size after concurrent access: %d\n\n", cache.Len())

	fmt.Println("=== Watching session:abc expire (TTL 1s) ===")
	for i := range 4 {
		v, ok := cache.Get("session:abc")
		fmt.Printf("t=%dms  session:abc -> value=%q found=%v\n", i*500, v, ok)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n=== Final state ===")
	fmt.Printf("config:env still present: %v\n", func() bool {
		_, ok := cache.Get("config:env")
		return ok
	}())
	fmt.Printf("Cache size: %d\n", cache.Len())

	cache.Close()
	cache.Close()
	fmt.Println("\nCache closed safely, even when Close() called twice")
}

package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

func (e entry) expired() bool {
	return time.Now().After(e.expiresAt)
}

// Store is a thread-safe in-memory key-value cache with TTL expiry
type Store struct {
	mu     sync.RWMutex
	data   map[string]entry
	cancel context.CancelFunc
	once   sync.Once
}

func NewStore(cleanupInterval time.Duration) *Store {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Store{
		data:   make(map[string]entry),
		cancel: cancel,
	}
	go s.cleanupLoop(ctx, cleanupInterval)
	return s
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[key]
	if !ok || e.expired() {
		return "", false
	}
	return e.value, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Store) Close() {
	s.once.Do(s.cancel)
}

func (s *Store) cleanupLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			removed := 0
			for k, e := range s.data {
				if e.expired() {
					delete(s.data, k)
					removed++
				}
			}
			s.mu.Unlock()
			if removed > 0 {
				fmt.Printf("cache: evicted %d expired entries\n", removed)
			}
		}
	}
}

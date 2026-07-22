package main

import (
	"fmt"
	"sync"
)

type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

func (kv *KVStore) Set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = value
}

func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[key]
	return val, ok
}

func (kv *KVStore) Delete(key string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, key)
}

func (kv *KVStore) Len() int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return len(kv.data)
}

func main() {
	kv := NewKVStore()
	var wg sync.WaitGroup

	for i := range 10 {
		kv.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	for i := range 100 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id%10)
			if val, ok := kv.Get(key); ok {
				fmt.Printf("reader %3d got: %s = %s\n", id, key, val)
			} else {
				fmt.Printf("reader %3d: %s not found\n", id, key)
			}
		}(i)
	}

	for i := range 2 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", id)
			newVal := fmt.Sprintf("newValue%d", id)
			kv.Set(key, newVal)
			fmt.Printf("writer %d wrote: %s = %s\n", id, key, newVal)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nFinal store size: %d\n", kv.Len())

	fmt.Println("Final values:")
	for i := range 10 {
		key := fmt.Sprintf("key%d", i)
		if val, ok := kv.Get(key); ok {
			fmt.Printf("  %s = %s\n", key, val)
		}
	}
}

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s := NewStore(100 * time.Millisecond)
	t.Cleanup(s.Close) // automatically called when test finishes
	return s
}

func TestSetAndGet(t *testing.T) {
	s := newTestStore(t)

	s.Set("key1", "value1", 1*time.Second)

	val, ok := s.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestGetMissing(t *testing.T) {
	s := newTestStore(t)

	val, ok := s.Get("nonexistent")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestExpiry(t *testing.T) {
	s := newTestStore(t)

	s.Set("short", "value", 50*time.Millisecond)

	// before expiry
	val, ok := s.Get("short")
	require.True(t, ok)
	assert.Equal(t, "value", val)

	// after expiry
	time.Sleep(100 * time.Millisecond)
	val, ok = s.Get("short")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)

	s.Set("key", "value", 1*time.Second)
	s.Delete("key")

	_, ok := s.Get("key")
	assert.False(t, ok)
}

func TestLen(t *testing.T) {
	s := newTestStore(t)

	assert.Equal(t, 0, s.Len())
	s.Set("a", "1", time.Second)
	s.Set("b", "2", time.Second)
	assert.Equal(t, 2, s.Len())
	s.Delete("a")
	assert.Equal(t, 1, s.Len())
}

func TestOverwrite(t *testing.T) {
	s := newTestStore(t)

	s.Set("key", "original", time.Second)
	s.Set("key", "updated", time.Second)

	val, ok := s.Get("key")
	require.True(t, ok)
	assert.Equal(t, "updated", val)
}

func TestConcurrentAccess(t *testing.T) {
	s := newTestStore(t)

	// pre-seed
	for i := range 10 {
		s.Set(string(rune('a'+i)), "value", time.Second)
	}

	done := make(chan struct{})

	// concurrent readers
	for range 50 {
		go func() {
			s.Get(string(rune('a' + 0)))
			done <- struct{}{}
		}()
	}

	// concurrent writers
	for i := range 5 {
		go func(n int) {
			s.Set(string(rune('a'+n)), "new", time.Second)
			done <- struct{}{}
		}(i)
	}

	// wait for all 55 goroutines
	for range 55 {
		<-done
	}
}

func BenchmarkSet(b *testing.B) {
	s := NewStore(time.Minute) // long cleanup interval for benchmark
	defer s.Close()
	b.ResetTimer()

	for i := range b.N {
		s.Set(string(rune(i)), "value", time.Second)
	}
}

func BenchmarkGet(b *testing.B) {
	s := NewStore(time.Minute)
	defer s.Close()
	s.Set("key", "value", time.Minute)
	b.ResetTimer()

	for b.Loop() {
		s.Get("key")
	}
}

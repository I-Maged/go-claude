package jobqueue

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestQueue(t *testing.T, workers, buf int) *Queue {
	t.Helper()
	q := New(workers, buf)
	t.Cleanup(q.Shutdown)
	return q
}

func TestSubmitAndProcess(t *testing.T) {
	q := newTestQueue(t, 2, 10)

	var processed atomic.Int64
	for i := range 5 {
		err := q.Submit(Job{
			ID:      i,
			Payload: "task",
			Execute: func() error {
				processed.Add(1)
				return nil
			},
		})
		require.NoError(t, err)
	}

	// give workers time to process
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, int64(5), processed.Load())
}

func TestQueueFull(t *testing.T) {
	// 0 workers — nothing drains the channel
	// buffer of 2 — fills immediately
	q := New(0, 2)
	defer q.Shutdown()

	err1 := q.Submit(Job{ID: 1, Execute: func() error { return nil }})
	err2 := q.Submit(Job{ID: 2, Execute: func() error { return nil }})
	err3 := q.Submit(Job{ID: 3, Execute: func() error { return nil }})

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	require.Error(t, err3)
	assert.ErrorIs(t, err3, ErrQueueFull)
}

func TestFailedJob(t *testing.T) {
	q := newTestQueue(t, 1, 5)

	var attempted atomic.Int64
	boom := errors.New("job failed")

	for i := range 3 {
		err := q.Submit(Job{
			ID: i,
			Execute: func() error {
				attempted.Add(1)
				return boom // always fails
			},
		})
		require.NoError(t, err)
	}

	time.Sleep(200 * time.Millisecond)

	// all 3 were attempted even though they failed
	assert.Equal(t, int64(3), attempted.Load())

	// stats reflect correct counts
	q.mu.Lock()
	assert.Equal(t, 3, q.submitted)
	assert.Equal(t, 0, q.processed) // none succeeded
	assert.Equal(t, 3, q.failed)
	q.mu.Unlock()
}

func TestShutdownDrainsInFlight(t *testing.T) {
	q := New(2, 10)

	var processed atomic.Int64
	for i := range 10 {
		q.Submit(Job{
			ID: i,
			Execute: func() error {
				time.Sleep(20 * time.Millisecond)
				processed.Add(1)
				return nil
			},
		})
	}

	// shutdown while jobs are still running
	q.Shutdown()

	// jobs that were already picked up by workers finish
	// (exact count depends on timing — just verify some ran)
	assert.Positive(t, processed.Load())
}

func TestSubmittedCount(t *testing.T) {
	q := newTestQueue(t, 2, 20)

	for i := range 10 {
		q.Submit(Job{ID: i, Execute: func() error { return nil }})
	}

	q.mu.Lock()
	assert.Equal(t, 10, q.submitted)
	q.mu.Unlock()
}

func BenchmarkSubmit(b *testing.B) {
	q := New(4, b.N+1) // buffer large enough to never block
	defer q.Shutdown()
	b.ResetTimer()

	for i := range b.N {
		q.Submit(Job{
			ID:      i,
			Execute: func() error { return nil },
		})
	}
}

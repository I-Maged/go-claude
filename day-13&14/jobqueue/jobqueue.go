package jobqueue

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrQueueFull is returned by Submit when the buffer is at capacity
var ErrQueueFull = errors.New("job queue full")

// Job is the unit of work the queue processes
type Job struct {
	ID      int
	Payload string
	// Execute is the actual work — caller provides this
	Execute func() error
}

// Queue distributes jobs across a fixed pool of workers
type Queue struct {
	jobCh     chan Job
	mu        sync.Mutex
	submitted int
	processed int
	failed    int
	wg        sync.WaitGroup
	cancel    context.CancelFunc
}

func New(numWorkers, bufferSize int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobCh:  make(chan Job, bufferSize),
		cancel: cancel,
	}
	for range numWorkers {
		q.wg.Add(1)
		go q.worker(ctx)
	}
	return q
}

func (q *Queue) worker(ctx context.Context) {
	defer q.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-q.jobCh:
			if !ok {
				return
			}
			err := job.Execute()
			q.mu.Lock()
			if err != nil {
				q.failed++
				fmt.Printf("  ✗ job %d failed: %v\n", job.ID, err)
			} else {
				q.processed++
				fmt.Printf("  ✓ job %d done: %s\n", job.ID, job.Payload)
			}
			q.mu.Unlock()
		}
	}
}

func (q *Queue) Submit(job Job) error {
	select {
	case q.jobCh <- job:
		q.mu.Lock()
		q.submitted++
		q.mu.Unlock()
		return nil
	default:
		return fmt.Errorf("%w: job %d", ErrQueueFull, job.ID)
	}
}

func (q *Queue) Shutdown() {
	q.cancel()
	q.wg.Wait()
	close(q.jobCh)
}

func (q *Queue) Stats() {
	q.mu.Lock()
	defer q.mu.Unlock()
	fmt.Printf("submitted=%d processed=%d failed=%d\n",
		q.submitted, q.processed, q.failed)
}

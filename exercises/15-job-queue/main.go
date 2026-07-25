package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

func (j Job) Execute() {
	time.Sleep(time.Duration(50+j.ID%50) * time.Millisecond)
	fmt.Printf("  ✓ job %2d processed: %s\n", j.ID, j.Payload)
}

type JobQueue struct {
	jobCh     chan Job
	submitted int
	mu        sync.Mutex
	wg        sync.WaitGroup
	cancel    context.CancelFunc
}

func NewJobQueue(numWorkers, bufferSize int) *JobQueue {
	ctx, cancel := context.WithCancel(context.Background())

	q := &JobQueue{
		jobCh:  make(chan Job, bufferSize),
		cancel: cancel,
	}

	for range numWorkers {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-q.jobCh:
					if !ok {
						return
					}
					job.Execute()
				}
			}
		}()
	}

	return q
}

func (q *JobQueue) Submit(job Job) error {
	select {
	case q.jobCh <- job:
		q.mu.Lock()
		q.submitted++
		q.mu.Unlock()
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *JobQueue) Shutdown() {
	fmt.Println("\nShutdown: signalling workers...")
	q.cancel()
	q.wg.Wait()
	close(q.jobCh)
	fmt.Println("Shutdown: all workers stopped")
}

func (q *JobQueue) Stats() {
	q.mu.Lock()
	defer q.mu.Unlock()
	fmt.Printf("Stats: %d jobs submitted\n", q.submitted)
}

func main() {
	fmt.Println("=== Starting job queue (4 workers, buffer 10) ===")
	queue := NewJobQueue(4, 10)

	fmt.Println("\n--- Burst 1: jobs 1-8 ---")
	for i := 1; i <= 8; i++ {
		err := queue.Submit(Job{ID: i, Payload: fmt.Sprintf("task-%d", i)})
		if err != nil {
			fmt.Printf("  job %d rejected: %v\n", i, err)
		}
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n--- Burst 2: jobs 9-16 ---")
	for i := 9; i <= 16; i++ {
		err := queue.Submit(Job{ID: i, Payload: fmt.Sprintf("task-%d", i)})
		if err != nil {
			fmt.Printf("  job %d rejected: %v\n", i, err)
		}
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n--- Burst 3: jobs 17-20 ---")
	for i := 17; i <= 20; i++ {
		err := queue.Submit(Job{ID: i, Payload: fmt.Sprintf("task-%d", i)})
		if err != nil {
			fmt.Printf("  job %d rejected: %v\n", i, err)
		}
	}

	time.Sleep(300 * time.Millisecond)
	queue.Shutdown()
	queue.Stats()
}

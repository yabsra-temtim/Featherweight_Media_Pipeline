package worker

import (
	"context"
	"log"
	"sync"
)

// JobProcessor defines the minimal contract required by worker pool jobs.
type JobProcessor interface {
	ProcessJob(ctx context.Context, jobID string) error
}

type Pool struct {
	processor   JobProcessor
	queue       chan string
	workerCount int
	wg          sync.WaitGroup
}

func NewPool(workerCount, queueSize int, processor JobProcessor) *Pool {
	return &Pool{
		processor:   processor,
		queue:       make(chan string, queueSize),
		workerCount: workerCount,
	}
}

// Start boots up the worker goroutines
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
	log.Printf("Worker pool started with %d workers", p.workerCount)
}

// Stop gracefully shuts down the workers
func (p *Pool) Stop() {
	close(p.queue) // Stop accepting new jobs
	p.wg.Wait()    // Wait for current jobs to finish
	log.Println("Worker pool stopped")
}

// Submit adds a job to the queue. Returns false if the queue is full (caller should return 503).
func (p *Pool) Submit(jobID string) bool {
	select {
	case p.queue <- jobID:
		return true
	default:
		return false
	}
}

// worker is the actual goroutine that processes jobs
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for jobID := range p.queue {
		log.Printf("[Worker %d] Processing job %s", id, jobID)
		if err := p.processor.ProcessJob(ctx, jobID); err != nil {
			log.Printf("[Worker %d] Failed job %s: %v", id, jobID, err)
		} else {
			log.Printf("[Worker %d] Completed job %s", id, jobID)
		}
	}
}

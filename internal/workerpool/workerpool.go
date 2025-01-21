package workerpool

import (
	"sync"
)

// Job represents a task to be processed by a worker
type Job struct {
	Task func() any // The task function to execute
	Resp chan<- any // Channel to send the result back
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	jobQueue chan Job       // Channel for jobs
	wg       sync.WaitGroup // WaitGroup to wait for all workers to finish
}

// NewWorkerPool creates a new worker pool with a fixed number of workers
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		jobQueue: make(chan Job, 100), // Buffer for 100 jobs
	}

	// Start the workers
	pool.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go pool.worker()
	}

	return pool
}

// worker processes jobs from the job queue
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for job := range p.jobQueue {
		// Execute the task and send the result back
		result := job.Task()
		job.Resp <- result
	}
}

// Submit adds a new job to the worker pool
func (p *WorkerPool) Submit(task func() any) <-chan any {
	resp := make(chan any, 1)
	p.jobQueue <- Job{
		Task: task,
		Resp: resp,
	}
	return resp
}

// Close waits for all workers to finish and closes the job queue
func (p *WorkerPool) Close() {
	close(p.jobQueue)
	p.wg.Wait()
}

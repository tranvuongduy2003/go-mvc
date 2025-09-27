package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/jobs"
)

// Worker represents a single worker that processes jobs
type Worker struct {
	id       string
	queue    jobs.JobQueue
	handlers map[string]jobs.JobHandler
	metrics  jobs.JobMetrics
	running  bool
	shutdown chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

// NewWorker creates a new worker
func NewWorker(id string, queue jobs.JobQueue) *Worker {
	return NewWorkerWithJobMetrics(id, queue, nil)
}

// NewWorkerWithJobMetrics creates a new worker with job metrics
func NewWorkerWithJobMetrics(id string, queue jobs.JobQueue, metrics jobs.JobMetrics) *Worker {
	if id == "" {
		id = uuid.New().String()
	}

	return &Worker{
		id:       id,
		queue:    queue,
		handlers: make(map[string]jobs.JobHandler),
		metrics:  metrics,
		shutdown: make(chan struct{}),
	}
}

// GetWorkerID returns the unique identifier of the worker
func (w *Worker) GetWorkerID() string {
	return w.id
}

// IsRunning returns whether the worker is currently running
func (w *Worker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

// RegisterHandler registers a job handler for a specific job type
func (w *Worker) RegisterHandler(handler jobs.JobHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[handler.GetJobType()] = handler
}

// Start begins processing jobs from the queue
func (w *Worker) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return fmt.Errorf("worker %s is already running", w.id)
	}
	w.running = true
	w.mu.Unlock()

	w.wg.Add(1)
	go w.run(ctx)

	return nil
}

// Stop gracefully stops the worker
func (w *Worker) Stop(ctx context.Context) error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return fmt.Errorf("worker %s is not running", w.id)
	}
	w.running = false
	w.mu.Unlock()

	// Signal shutdown
	close(w.shutdown)

	// Wait for worker to finish with timeout
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("worker %s shutdown timeout", w.id)
	}
}

// run is the main worker loop
func (w *Worker) run(ctx context.Context) {
	defer w.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.shutdown:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processJob(ctx)
		}
	}
}

// processJob processes a single job from the queue
func (w *Worker) processJob(ctx context.Context) {
	// Try to dequeue a job
	job, err := w.queue.Dequeue(ctx)
	if err != nil {
		// Log error but continue
		return
	}

	if job == nil {
		// No job available
		return
	}

	// Start timing for metrics
	start := time.Now()

	// Find handler for this job type
	w.mu.RLock()
	handler, exists := w.handlers[job.GetType()]
	w.mu.RUnlock()

	if !exists {
		// No handler found, fail the job
		w.queue.NackJob(ctx, job, fmt.Errorf("no handler found for job type: %s", job.GetType()))
		// Record failure metrics if available
		if w.metrics != nil {
			w.metrics.IncrementJobsProcessed(job.GetType(), false)
		}
		return
	}

	// Execute the job
	if err := handler.Execute(ctx, job); err != nil {
		// Job failed, nack it
		w.queue.NackJob(ctx, job, err)
		// Record failure metrics if available
		if w.metrics != nil {
			w.metrics.IncrementJobsProcessed(job.GetType(), false)
			w.metrics.IncrementJobRetries(job.GetType())
		}
		return
	}

	// Job succeeded, ack it
	w.queue.AckJob(ctx, job)
	// Record success metrics if available
	if w.metrics != nil {
		duration := time.Since(start)
		w.metrics.ObserveJobDuration(job.GetType(), duration)
		w.metrics.IncrementJobsProcessed(job.GetType(), true)
	}
}

// WorkerPool manages multiple workers
type WorkerPool struct {
	workers     map[string]*Worker
	queue       jobs.JobQueue
	workerCount int
	mu          sync.RWMutex
	running     bool
	shutdown    chan struct{}
	wg          sync.WaitGroup

	// Stats
	stats   WorkerPoolStats
	statsMu sync.RWMutex
}

// WorkerPoolStats provides statistics about the worker pool
type WorkerPoolStats struct {
	ActiveWorkers      int   `json:"active_workers"`
	TotalJobsProcessed int64 `json:"total_jobs_processed"`
	SuccessfulJobs     int64 `json:"successful_jobs"`
	FailedJobs         int64 `json:"failed_jobs"`
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(queue jobs.JobQueue, workerCount int) *WorkerPool {
	if workerCount <= 0 {
		workerCount = 1
	}

	return &WorkerPool{
		workers:     make(map[string]*Worker),
		queue:       queue,
		workerCount: workerCount,
		shutdown:    make(chan struct{}),
	}
}

// Start starts all workers in the pool
func (wp *WorkerPool) Start(ctx context.Context) error {
	wp.mu.Lock()
	if wp.running {
		wp.mu.Unlock()
		return fmt.Errorf("worker pool is already running")
	}
	wp.running = true
	wp.mu.Unlock()

	// Create and start workers
	for i := 0; i < wp.workerCount; i++ {
		workerID := fmt.Sprintf("worker-%d", i+1)
		worker := NewWorker(workerID, wp.queue)

		wp.mu.Lock()
		wp.workers[workerID] = worker
		wp.mu.Unlock()

		if err := worker.Start(ctx); err != nil {
			// If any worker fails to start, stop all started workers
			wp.Stop(ctx)
			return fmt.Errorf("failed to start worker %s: %w", workerID, err)
		}
	}

	wp.updateStats()
	return nil
}

// Stop stops all workers gracefully
func (wp *WorkerPool) Stop(ctx context.Context) error {
	wp.mu.Lock()
	if !wp.running {
		wp.mu.Unlock()
		return fmt.Errorf("worker pool is not running")
	}
	wp.running = false
	wp.mu.Unlock()

	// Signal shutdown
	close(wp.shutdown)

	// Stop all workers
	var wg sync.WaitGroup
	wp.mu.RLock()
	workers := make([]*Worker, 0, len(wp.workers))
	for _, worker := range wp.workers {
		workers = append(workers, worker)
	}
	wp.mu.RUnlock()

	for _, worker := range workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.Stop(ctx)
		}(worker)
	}

	// Wait for all workers to stop
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		wp.mu.Lock()
		wp.workers = make(map[string]*Worker)
		wp.mu.Unlock()
		wp.updateStats()
		return nil
	case <-ctx.Done():
		return fmt.Errorf("worker pool shutdown timeout")
	}
}

// AddWorker adds a worker to the pool
func (wp *WorkerPool) AddWorker(worker *Worker) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		// If pool is running, start the worker
		go func() {
			ctx := context.Background()
			worker.Start(ctx)
		}()
	}

	wp.workers[worker.GetWorkerID()] = worker
	wp.updateStats()
}

// RemoveWorker removes a worker from the pool
func (wp *WorkerPool) RemoveWorker(workerID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if worker, exists := wp.workers[workerID]; exists {
		// Stop the worker if it's running
		if worker.IsRunning() {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				worker.Stop(ctx)
			}()
		}
		delete(wp.workers, workerID)
		wp.updateStats()
	}
}

// RegisterHandler registers a handler with all workers
func (wp *WorkerPool) RegisterHandler(handler jobs.JobHandler) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	for _, worker := range wp.workers {
		worker.RegisterHandler(handler)
	}
}

// GetWorkerCount returns the number of active workers
func (wp *WorkerPool) GetWorkerCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return len(wp.workers)
}

// GetStats returns worker pool statistics
func (wp *WorkerPool) GetStats() WorkerPoolStats {
	wp.statsMu.RLock()
	defer wp.statsMu.RUnlock()
	return wp.stats
}

// updateStats updates the worker pool statistics
func (wp *WorkerPool) updateStats() {
	wp.statsMu.Lock()
	defer wp.statsMu.Unlock()

	wp.mu.RLock()
	activeCount := 0
	for _, worker := range wp.workers {
		if worker.IsRunning() {
			activeCount++
		}
	}
	wp.mu.RUnlock()

	wp.stats.ActiveWorkers = activeCount
}

// IncrementSuccessfulJobs increments the successful jobs counter
func (wp *WorkerPool) IncrementSuccessfulJobs() {
	wp.statsMu.Lock()
	defer wp.statsMu.Unlock()
	wp.stats.SuccessfulJobs++
	wp.stats.TotalJobsProcessed++
}

// IncrementFailedJobs increments the failed jobs counter
func (wp *WorkerPool) IncrementFailedJobs() {
	wp.statsMu.Lock()
	defer wp.statsMu.Unlock()
	wp.stats.FailedJobs++
	wp.stats.TotalJobsProcessed++
}

// WorkerWithMetrics wraps a worker to collect metrics
type WorkerWithMetrics struct {
	*Worker
	pool *WorkerPool
}

// NewWorkerWithMetrics creates a new worker with metrics collection
func NewWorkerWithMetrics(id string, queue jobs.JobQueue, pool *WorkerPool) *WorkerWithMetrics {
	worker := NewWorker(id, queue)
	return &WorkerWithMetrics{
		Worker: worker,
		pool:   pool,
	}
}

// processJob overrides the parent method to collect metrics
func (wm *WorkerWithMetrics) processJob(ctx context.Context) {
	// Try to dequeue a job
	job, err := wm.queue.Dequeue(ctx)
	if err != nil {
		return
	}

	if job == nil {
		return
	}

	// Find handler for this job type
	wm.mu.RLock()
	handler, exists := wm.handlers[job.GetType()]
	wm.mu.RUnlock()

	if !exists {
		// No handler found, fail the job
		wm.queue.NackJob(ctx, job, fmt.Errorf("no handler found for job type: %s", job.GetType()))
		wm.pool.IncrementFailedJobs()
		return
	}

	// Execute the job
	if err := handler.Execute(ctx, job); err != nil {
		// Job failed, nack it
		wm.queue.NackJob(ctx, job, err)
		wm.pool.IncrementFailedJobs()
		return
	}

	// Job succeeded, ack it
	wm.queue.AckJob(ctx, job)
	wm.pool.IncrementSuccessfulJobs()
}

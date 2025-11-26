package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

type Worker struct {
	id       string
	queue    job.JobQueue
	handlers map[string]job.JobHandler
	metrics  job.JobMetrics
	running  bool
	shutdown chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

func NewWorker(id string, queue job.JobQueue) *Worker {
	return NewWorkerWithJobMetrics(id, queue, nil)
}

func NewWorkerWithJobMetrics(id string, queue job.JobQueue, metrics job.JobMetrics) *Worker {
	if id == "" {
		id = uuid.New().String()
	}

	return &Worker{
		id:       id,
		queue:    queue,
		handlers: make(map[string]job.JobHandler),
		metrics:  metrics,
		shutdown: make(chan struct{}),
	}
}

func (w *Worker) GetWorkerID() string {
	return w.id
}

func (w *Worker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

func (w *Worker) RegisterHandler(handler job.JobHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[handler.GetJobType()] = handler
}

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

func (w *Worker) Stop(ctx context.Context) error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return fmt.Errorf("worker %s is not running", w.id)
	}
	w.running = false
	w.mu.Unlock()

	close(w.shutdown)

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

func (w *Worker) processJob(ctx context.Context) {
	job, err := w.queue.Dequeue(ctx)
	if err != nil {
		return
	}

	if job == nil {
		return
	}

	start := time.Now()

	w.mu.RLock()
	handler, exists := w.handlers[job.GetType()]
	w.mu.RUnlock()

	if !exists {
		w.queue.NackJob(ctx, job, fmt.Errorf("no handler found for job type: %s", job.GetType()))
		if w.metrics != nil {
			w.metrics.IncrementJobsProcessed(job.GetType(), false)
		}
		return
	}

	if err := handler.Execute(ctx, job); err != nil {
		w.queue.NackJob(ctx, job, err)
		if w.metrics != nil {
			w.metrics.IncrementJobsProcessed(job.GetType(), false)
			w.metrics.IncrementJobRetries(job.GetType())
		}
		return
	}

	w.queue.AckJob(ctx, job)
	if w.metrics != nil {
		duration := time.Since(start)
		w.metrics.ObserveJobDuration(job.GetType(), duration)
		w.metrics.IncrementJobsProcessed(job.GetType(), true)
	}
}

type WorkerPool struct {
	workers     map[string]*Worker
	queue       job.JobQueue
	workerCount int
	mu          sync.RWMutex
	running     bool
	shutdown    chan struct{}
	wg          sync.WaitGroup

	stats   WorkerPoolStats
	statsMu sync.RWMutex
}

type WorkerPoolStats struct {
	ActiveWorkers      int   `json:"active_workers"`
	TotalJobsProcessed int64 `json:"total_jobs_processed"`
	SuccessfulJobs     int64 `json:"successful_jobs"`
	FailedJobs         int64 `json:"failed_jobs"`
}

func NewWorkerPool(queue job.JobQueue, workerCount int) *WorkerPool {
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

func (wp *WorkerPool) Start(ctx context.Context) error {
	wp.mu.Lock()
	if wp.running {
		wp.mu.Unlock()
		return fmt.Errorf("worker pool is already running")
	}
	wp.running = true
	wp.mu.Unlock()

	for i := 0; i < wp.workerCount; i++ {
		workerID := fmt.Sprintf("worker-%d", i+1)
		worker := NewWorker(workerID, wp.queue)

		wp.mu.Lock()
		wp.workers[workerID] = worker
		wp.mu.Unlock()

		if err := worker.Start(ctx); err != nil {
			wp.Stop(ctx)
			return fmt.Errorf("failed to start worker %s: %w", workerID, err)
		}
	}

	wp.updateStats()
	return nil
}

func (wp *WorkerPool) Stop(ctx context.Context) error {
	wp.mu.Lock()
	if !wp.running {
		wp.mu.Unlock()
		return fmt.Errorf("worker pool is not running")
	}
	wp.running = false
	wp.mu.Unlock()

	close(wp.shutdown)

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

func (wp *WorkerPool) AddWorker(worker *Worker) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		go func() {
			ctx := context.Background()
			worker.Start(ctx)
		}()
	}

	wp.workers[worker.GetWorkerID()] = worker
	wp.updateStats()
}

func (wp *WorkerPool) RemoveWorker(workerID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if worker, exists := wp.workers[workerID]; exists {
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

func (wp *WorkerPool) RegisterHandler(handler job.JobHandler) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	for _, worker := range wp.workers {
		worker.RegisterHandler(handler)
	}
}

func (wp *WorkerPool) GetWorkerCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return len(wp.workers)
}

func (wp *WorkerPool) GetStats() WorkerPoolStats {
	wp.statsMu.RLock()
	defer wp.statsMu.RUnlock()
	return wp.stats
}

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

func (wp *WorkerPool) IncrementSuccessfulJobs() {
	wp.statsMu.Lock()
	defer wp.statsMu.Unlock()
	wp.stats.SuccessfulJobs++
	wp.stats.TotalJobsProcessed++
}

func (wp *WorkerPool) IncrementFailedJobs() {
	wp.statsMu.Lock()
	defer wp.statsMu.Unlock()
	wp.stats.FailedJobs++
	wp.stats.TotalJobsProcessed++
}

type WorkerWithMetrics struct {
	*Worker
	pool *WorkerPool
}

func NewWorkerWithMetrics(id string, queue job.JobQueue, pool *WorkerPool) *WorkerWithMetrics {
	worker := NewWorker(id, queue)
	return &WorkerWithMetrics{
		Worker: worker,
		pool:   pool,
	}
}

func (wm *WorkerWithMetrics) processJob(ctx context.Context) {
	job, err := wm.queue.Dequeue(ctx)
	if err != nil {
		return
	}

	if job == nil {
		return
	}

	wm.mu.RLock()
	handler, exists := wm.handlers[job.GetType()]
	wm.mu.RUnlock()

	if !exists {
		wm.queue.NackJob(ctx, job, fmt.Errorf("no handler found for job type: %s", job.GetType()))
		wm.pool.IncrementFailedJobs()
		return
	}

	if err := handler.Execute(ctx, job); err != nil {
		wm.queue.NackJob(ctx, job, err)
		wm.pool.IncrementFailedJobs()
		return
	}

	wm.queue.AckJob(ctx, job)
	wm.pool.IncrementSuccessfulJobs()
}

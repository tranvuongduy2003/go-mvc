# Background Jobs System

A comprehensive background job processing system built with Clean Architecture principles for the Go-MVC project.

## Overview

This background jobs system provides:

- **Queue-based Job Processing**: Redis-backed job queues with priority support
- **Worker Pool Management**: Configurable worker pools with concurrency control
- **Job Scheduling**: Support for delayed and scheduled job execution
- **Metrics and Monitoring**: Prometheus-based metrics collection
- **Clean Architecture**: Proper separation of concerns with ports and adapters pattern
- **Dependency Injection**: Uber FX-based lifecycle management

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Background Jobs System                    │
├─────────────────────────────────────────────────────────────┤
│ Application Layer                                           │
│  ├── Services (BackgroundJobService)                        │
│  └── Commands (SubmitJob, ScheduleJob, CancelJob)          │
├─────────────────────────────────────────────────────────────┤
│ Domain Layer                                               │
│  ├── Entities (BaseJob, EmailJob, FileProcessingJob)       │
│  ├── Value Objects (JobStatus, JobPriority, JobPayload)    │
│  └── Factories (JobFactory, JobValidator)                  │
├─────────────────────────────────────────────────────────────┤
│ Infrastructure Layer                                        │
│  ├── Job Queue (Redis-based with priority queues)          │
│  ├── Worker Pool (Concurrent job processing)               │
│  ├── Scheduler (Delayed and recurring jobs)                │
│  ├── Metrics (Prometheus metrics collection)               │
│  └── Handlers (EmailHandler, FileHandler, CleanupHandler)  │
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Job Queue System (`internal/core/ports/jobs/`)

**Interfaces:**
- `Job`: Core job interface with ID, type, payload, priority, and status
- `JobQueue`: Queue operations (Enqueue, Dequeue, Ack, Nack)
- `Worker`: Individual worker for job processing
- `WorkerPool`: Manages multiple workers
- `Scheduler`: Handles delayed and scheduled jobs
- `JobMetrics`: Metrics collection interface

### 2. Domain Models (`internal/core/domain/jobs/`)

**Job Types:**
- `BaseJob`: Foundation job implementation
- `EmailJob`: Email sending jobs
- `FileProcessingJob`: File processing and manipulation
- `DataCleanupJob`: Data cleanup and maintenance
- `NotificationJob`: Push notifications and alerts

**Factories:**
- `JobFactory`: Creates jobs with validation
- `JobValidator`: Validates job data and constraints

### 3. Redis Job Queue (`internal/adapters/jobs/redis/`)

**Features:**
- Priority-based job queuing (Low, Normal, High, Critical)
- Delayed job execution
- Retry mechanisms with exponential backoff
- Atomic operations for reliability
- Dead letter queues for failed jobs

**Redis Keys Structure:**
```
jobs:queue:{queue_name}:{priority}     # Priority queues
jobs:delayed                           # Delayed jobs (sorted set)
jobs:processing:{worker_id}            # Currently processing jobs
jobs:data:{job_id}                     # Job data storage
jobs:failed:{queue_name}              # Failed jobs
```

### 4. Worker System (`internal/adapters/jobs/worker/`)

**Worker Pool:**
- Configurable number of workers
- Graceful shutdown with job completion
- Worker health monitoring
- Job handler registration
- Statistics collection

**Worker Features:**
- Concurrent job processing
- Handler-based job execution
- Automatic retry on failures
- Metrics integration

### 5. Job Scheduler (`internal/adapters/jobs/scheduler/`)

**Capabilities:**
- One-time scheduled jobs
- Recurring job patterns
- Cron-like scheduling
- Redis persistence
- Timezone support

### 6. Metrics System (`internal/adapters/jobs/metrics/`)

**Prometheus Metrics:**
- `go_mvc_jobs_enqueued_total`: Jobs enqueued by type and priority
- `go_mvc_jobs_processed_total`: Jobs processed by status
- `go_mvc_jobs_duration_seconds`: Job processing time
- `go_mvc_jobs_queue_size`: Current queue sizes
- `go_mvc_jobs_active_workers`: Active worker count
- `go_mvc_jobs_retries_total`: Job retry attempts

### 7. Job Handlers (`internal/adapters/jobs/handlers/`)

#### Email Handler
- Welcome emails
- Password reset emails
- Notification emails
- Bulk email campaigns
- Template-based emails

#### File Processing Handler
- Image processing (resize, optimize, format conversion)
- Video processing (compression, format conversion)
- Document processing (PDF generation, text extraction)
- Batch file operations
- File validation

#### Cleanup Handler
- User data cleanup (inactive, unverified, deleted users)
- Session cleanup (expired, orphaned sessions)
- File cleanup (temp files, orphaned files)
- Log cleanup (old logs, compression, rotation)
- Full system maintenance

## Usage Examples

### 1. Basic Job Submission

```go
// Submit an email job
jobID, err := jobService.SubmitJob(ctx, "email", jobs.JobPayload{
    "email_type": "welcome",
    "to":         "user@example.com",
    "username":   "John Doe",
})
```

### 2. Job with Options

```go
// Submit a high-priority job with delay
opts := jobs.JobOptions{
    Priority:    jobs.PriorityHigh,
    Delay:       &time.Duration(5 * time.Minute),
    Queue:       "urgent",
    MaxRetries:  3,
}

jobID, err := jobService.SubmitJobWithOptions(ctx, "email", payload, opts)
```

### 3. Scheduled Jobs

```go
// Schedule a job for specific time
scheduledTime := time.Now().Add(1 * time.Hour)
opts := jobs.JobOptions{
    ScheduledAt: &scheduledTime,
    Priority:    jobs.PriorityNormal,
}

jobID, err := jobService.SubmitJobWithOptions(ctx, "data_cleanup", payload, opts)
```

### 4. Custom Job Handler

```go
type CustomJobHandler struct {
    service CustomService
    metrics jobs.JobMetrics
}

func (h *CustomJobHandler) Execute(ctx context.Context, job jobs.Job) error {
    start := time.Now()
    defer func() {
        h.metrics.ObserveJobDuration(job.GetType(), time.Since(start))
    }()
    
    // Process the job
    err := h.service.ProcessJob(ctx, job)
    
    // Record metrics
    success := err == nil
    h.metrics.IncrementJobsProcessed(job.GetType(), success)
    
    return err
}

func (h *CustomJobHandler) GetJobType() string {
    return "custom_job"
}
```

## Configuration

### Job System Configuration

```yaml
jobs:
  enabled: true
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
  worker_pool:
    worker_count: 5
    max_retries: 3
    retry_delay: "5s"
  scheduler:
    enabled: true
    cleanup_interval: "1m"
  metrics:
    enabled: true
    collection_interval: "30s"
```

### Queue Configuration

```yaml
queues:
  default:
    priority_levels: ["low", "normal", "high", "critical"]
    max_jobs: 1000
  email:
    priority_levels: ["normal", "high"]
    max_jobs: 5000
  file_processing:
    priority_levels: ["low", "normal"]
    max_jobs: 100
  cleanup:
    priority_levels: ["low"]
    max_jobs: 50
```

## Monitoring and Observability

### Metrics Dashboard

The system provides comprehensive metrics through Prometheus:

1. **Job Throughput**: Jobs processed per second by type and status
2. **Queue Health**: Queue sizes and processing rates
3. **Worker Performance**: Active workers and processing times
4. **Error Rates**: Failed jobs and retry patterns
5. **Business Metrics**: Custom metrics for specific job types

### Grafana Dashboard

A pre-configured Grafana dashboard is available at:
`configs/grafana/dashboards/background-jobs-dashboard.json`

Key panels:
- Job Processing Rate
- Queue Sizes by Priority
- Worker Pool Utilization
- Processing Time Distribution
- Error Rate by Job Type

### Health Checks

Health check endpoints:
- `/health/jobs`: Overall job system health
- `/health/queues`: Queue status and sizes
- `/health/workers`: Worker pool status

## Running the Example

1. **Start Redis**:
   ```bash
   docker run -d -p 6379:6379 redis:alpine
   ```

2. **Run the example**:
   ```bash
   cd examples/background_jobs
   go run main.go
   ```

3. **Monitor metrics**:
   ```bash
   curl http://localhost:8080/metrics
   ```

## Development and Testing

### Adding New Job Types

1. **Create domain model**:
   ```go
   // internal/core/domain/jobs/my_job.go
   type MyJob struct {
       *BaseJob
       SpecificField string
   }
   ```

2. **Create job handler**:
   ```go
   // internal/adapters/jobs/handlers/my_handler.go
   type MyJobHandler struct {
       service MyService
       metrics jobs.JobMetrics
   }
   ```

3. **Register with factory**:
   ```go
   // Update JobFactory to handle new job type
   ```

4. **Add tests**:
   ```go
   // Test job creation, validation, and processing
   ```

### Testing Jobs

```go
func TestMyJobHandler(t *testing.T) {
    // Create mock services
    mockService := &MockMyService{}
    mockMetrics := &MockJobMetrics{}
    
    // Create handler
    handler := NewMyJobHandler(mockService, mockMetrics)
    
    // Create test job
    job := &MyJob{
        BaseJob: createBaseJob("my_job", payload),
        SpecificField: "test_value",
    }
    
    // Execute and verify
    err := handler.Execute(context.Background(), job)
    assert.NoError(t, err)
    assert.True(t, mockService.ProcessJobCalled)
}
```

## Production Considerations

### Scaling

1. **Horizontal Scaling**: Run multiple application instances
2. **Worker Pool Tuning**: Adjust worker count based on job types
3. **Queue Partitioning**: Use different Redis instances for different queues
4. **Priority Management**: Balance high-priority and bulk jobs

### Reliability

1. **Redis Clustering**: Use Redis cluster for high availability
2. **Job Persistence**: Configure Redis persistence (AOF/RDB)
3. **Dead Letter Queues**: Monitor and handle failed jobs
4. **Circuit Breakers**: Implement circuit breakers for external services

### Monitoring

1. **Alerting**: Set up alerts for queue growth and failure rates
2. **Logging**: Comprehensive logging with structured data
3. **Tracing**: Distributed tracing for job execution paths
4. **Health Checks**: Regular health checks for all components

### Security

1. **Redis Security**: Use authentication and encryption
2. **Job Validation**: Validate all job payloads
3. **Resource Limits**: Implement job execution timeouts
4. **Access Control**: Limit job submission permissions

## Troubleshooting

### Common Issues

1. **Jobs Not Processing**:
   - Check worker pool status
   - Verify job handler registration
   - Check Redis connectivity

2. **High Memory Usage**:
   - Monitor queue sizes
   - Check for job data leaks
   - Review Redis memory usage

3. **Slow Job Processing**:
   - Profile job handlers
   - Check external service dependencies
   - Review worker pool configuration

4. **Failed Jobs Accumulating**:
   - Check dead letter queues
   - Review job validation logic
   - Monitor external service health

### Debug Commands

```bash
# Check Redis queues
redis-cli LLEN "jobs:queue:default:normal"

# Monitor job processing
redis-cli MONITOR

# Check metrics
curl http://localhost:8080/metrics | grep jobs_

# Worker pool status
curl http://localhost:8080/health/workers
```

## Future Enhancements

- [ ] Job dependencies and workflows
- [ ] Job result storage and retrieval
- [ ] Web UI for job monitoring
- [ ] Job archival and cleanup policies
- [ ] Multi-tenant job isolation
- [ ] Job execution sandboxing
- [ ] Advanced scheduling patterns
- [ ] Job result notifications

## Contributing

1. Follow the existing Clean Architecture patterns
2. Add comprehensive tests for new features
3. Update metrics for new job types
4. Document configuration options
5. Include example usage in the documentation
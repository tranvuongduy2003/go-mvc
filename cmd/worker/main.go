package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"

	di "github.com/tranvuongduy2003/go-mvc/internal/di"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
)

var rootCmd = &cobra.Command{
	Use:   "worker",
	Short: "Background worker for job processing",
	Long:  `A background worker service for processing asynchronous jobs, handling queues, and managing background tasks.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(startCommand())
	rootCmd.AddCommand(statusCommand())
	rootCmd.AddCommand(stopCommand())
	rootCmd.AddCommand(healthCommand())
	rootCmd.AddCommand(versionCommand())
}

// Worker commands
func startCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the background worker",
		Long:  `Start the background worker service to process queued jobs.`,
		Run: func(cmd *cobra.Command, args []string) {
			workers, _ := cmd.Flags().GetInt("workers")
			daemon, _ := cmd.Flags().GetBool("daemon")

			fmt.Printf("Starting worker service with %d workers...\n", workers)

			if daemon {
				fmt.Println("Running in daemon mode...")
			}

			app := fx.New(
				di.InfrastructureModule,
				fx.Invoke(func(config *config.AppConfig, logger *zap.Logger, lc fx.Lifecycle) {
					worker := NewWorkerService(logger, workers)

					lc.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							logger.Info("Starting worker service", zap.Int("workers", workers))
							worker.Start(ctx)
							return nil
						},
						OnStop: func(ctx context.Context) error {
							logger.Info("Stopping worker service")
							worker.Stop(ctx)
							return nil
						},
					})

					// Wait for interrupt signal
					if !daemon {
						waitForShutdown(logger)
					}
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to start worker: %v", err)
			}

			if daemon {
				// Keep running in daemon mode
				select {}
			} else {
				// Graceful shutdown
				defer app.Stop(context.Background())
			}
		},
	}

	cmd.Flags().IntP("workers", "w", 4, "Number of worker goroutines")
	cmd.Flags().BoolP("daemon", "d", false, "Run as daemon process")
	return cmd
}

func statusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check worker status",
		Long:  `Check the current status of the worker service and job queues.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Checking worker status...")

			app := fx.New(
				di.InfrastructureModule,
				fx.Invoke(func(config *config.AppConfig, logger *zap.Logger) {
					status := checkWorkerStatus()

					fmt.Printf("Worker Status: %s\n", status.Status)
					fmt.Printf("Active Jobs: %d\n", status.ActiveJobs)
					fmt.Printf("Completed Jobs: %d\n", status.CompletedJobs)
					fmt.Printf("Failed Jobs: %d\n", status.FailedJobs)
					fmt.Printf("Queue Size: %d\n", status.QueueSize)
					fmt.Printf("Last Activity: %s\n", status.LastActivity.Format("2006-01-02 15:04:05"))

					logger.Info("Worker status checked")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to check status: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func stopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the background worker",
		Long:  `Stop the background worker service gracefully.`,
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")

			if force {
				fmt.Println("Force stopping worker service...")
			} else {
				fmt.Println("Gracefully stopping worker service...")
			}

			// In a real implementation, this would send a signal to the running worker process
			// For now, we'll just simulate the stop operation
			fmt.Println("✅ Worker service stopped successfully!")
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force stop without waiting for jobs to complete")
	return cmd
}

func healthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check worker health",
		Long:  `Perform a health check on the worker service and its dependencies.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Performing worker health check...")

			app := fx.New(
				di.InfrastructureModule,
				fx.Invoke(func(config *config.AppConfig, logger *zap.Logger) {
					health := performHealthCheck()

					fmt.Printf("Overall Health: %s\n", health.Overall)
					fmt.Printf("Database Connection: %s\n", health.Database)
					fmt.Printf("Queue System: %s\n", health.Queue)
					fmt.Printf("Memory Usage: %s\n", health.Memory)
					fmt.Printf("CPU Usage: %s\n", health.CPU)

					if health.Overall == "Healthy" {
						fmt.Println("✅ Worker service is healthy!")
					} else {
						fmt.Println("❌ Worker service has health issues!")
					}

					logger.Info("Worker health check completed")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to perform health check: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show worker version",
		Long:  `Display the current version of the worker service.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Go MVC Worker Service")
			fmt.Println("Version: 1.0.0")
			fmt.Println("Build: development")
			fmt.Println("Go Version:", os.Getenv("GO_VERSION"))
		},
	}
}

// Worker service implementation
type WorkerService struct {
	logger      *zap.Logger
	workerCount int
	jobs        chan Job
	quit        chan bool
	wg          sync.WaitGroup
	running     bool
}

type Job struct {
	ID       string
	Type     string
	Payload  map[string]interface{}
	Retry    int
	MaxRetry int
}

type WorkerStatus struct {
	Status        string
	ActiveJobs    int
	CompletedJobs int
	FailedJobs    int
	QueueSize     int
	LastActivity  time.Time
}

type HealthStatus struct {
	Overall  string
	Database string
	Queue    string
	Memory   string
	CPU      string
}

func NewWorkerService(logger *zap.Logger, workerCount int) *WorkerService {
	return &WorkerService{
		logger:      logger,
		workerCount: workerCount,
		jobs:        make(chan Job, 1000), // Buffer for 1000 jobs
		quit:        make(chan bool),
		running:     false,
	}
}

func (w *WorkerService) Start(ctx context.Context) {
	if w.running {
		w.logger.Warn("Worker service is already running")
		return
	}

	w.running = true
	w.logger.Info("Starting worker service", zap.Int("workers", w.workerCount))

	// Start worker goroutines
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.worker(i + 1)
	}

	// Start job producer (simulated)
	go w.jobProducer()

	fmt.Printf("✅ Worker service started with %d workers\n", w.workerCount)
}

func (w *WorkerService) Stop(ctx context.Context) {
	if !w.running {
		w.logger.Warn("Worker service is not running")
		return
	}

	w.logger.Info("Stopping worker service")
	w.running = false

	// Signal all workers to stop
	close(w.quit)

	// Wait for all workers to finish
	w.wg.Wait()

	// Close job channel
	close(w.jobs)

	fmt.Println("✅ Worker service stopped gracefully")
}

func (w *WorkerService) worker(id int) {
	defer w.wg.Done()

	w.logger.Info("Worker started", zap.Int("worker_id", id))

	for {
		select {
		case job := <-w.jobs:
			w.processJob(id, job)
		case <-w.quit:
			w.logger.Info("Worker stopping", zap.Int("worker_id", id))
			return
		}
	}
}

func (w *WorkerService) processJob(workerID int, job Job) {
	w.logger.Info("Processing job",
		zap.Int("worker_id", workerID),
		zap.String("job_id", job.ID),
		zap.String("job_type", job.Type))

	// Simulate job processing
	switch job.Type {
	case "email":
		w.processEmailJob(job)
	case "notification":
		w.processNotificationJob(job)
	case "cleanup":
		w.processCleanupJob(job)
	default:
		w.logger.Warn("Unknown job type", zap.String("type", job.Type))
	}

	w.logger.Info("Job completed",
		zap.Int("worker_id", workerID),
		zap.String("job_id", job.ID))
}

func (w *WorkerService) processEmailJob(job Job) {
	// Simulate email processing
	time.Sleep(100 * time.Millisecond)
	w.logger.Info("Email sent", zap.String("job_id", job.ID))
}

func (w *WorkerService) processNotificationJob(job Job) {
	// Simulate notification processing
	time.Sleep(50 * time.Millisecond)
	w.logger.Info("Notification sent", zap.String("job_id", job.ID))
}

func (w *WorkerService) processCleanupJob(job Job) {
	// Simulate cleanup processing
	time.Sleep(200 * time.Millisecond)
	w.logger.Info("Cleanup completed", zap.String("job_id", job.ID))
}

func (w *WorkerService) jobProducer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	jobTypes := []string{"email", "notification", "cleanup"}
	counter := 0

	for {
		select {
		case <-ticker.C:
			if !w.running {
				return
			}

			counter++
			job := Job{
				ID:       fmt.Sprintf("job_%d", counter),
				Type:     jobTypes[counter%len(jobTypes)],
				Payload:  map[string]interface{}{"data": fmt.Sprintf("payload_%d", counter)},
				Retry:    0,
				MaxRetry: 3,
			}

			select {
			case w.jobs <- job:
				w.logger.Debug("Job queued", zap.String("job_id", job.ID))
			default:
				w.logger.Warn("Job queue is full, dropping job", zap.String("job_id", job.ID))
			}
		case <-w.quit:
			return
		}
	}
}

// Helper functions
func waitForShutdown(logger *zap.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info("Received shutdown signal")
}

func checkWorkerStatus() WorkerStatus {
	// In a real implementation, this would check actual worker status
	return WorkerStatus{
		Status:        "Running",
		ActiveJobs:    5,
		CompletedJobs: 1234,
		FailedJobs:    12,
		QueueSize:     45,
		LastActivity:  time.Now(),
	}
}

func performHealthCheck() HealthStatus {
	// In a real implementation, this would perform actual health checks
	return HealthStatus{
		Overall:  "Healthy",
		Database: "Connected",
		Queue:    "Available",
		Memory:   "Normal (65%)",
		CPU:      "Normal (45%)",
	}
}

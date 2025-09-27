package handlers

import (
	"context"
	"fmt"
	"time"

	domainjobs "github.com/tranvuongduy2003/go-mvc/internal/core/domain/jobs"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/jobs"
)

// DataCleanupJobHandler handles data cleanup and maintenance jobs
type DataCleanupJobHandler struct {
	userService    UserCleanupService
	sessionService SessionCleanupService
	fileService    FileCleanupService
	logService     LogCleanupService
	metrics        jobs.JobMetrics
}

// Cleanup service interfaces
type UserCleanupService interface {
	CleanupInactiveUsers(ctx context.Context, inactiveDays int) (int, error)
	CleanupUnverifiedUsers(ctx context.Context, unverifiedDays int) (int, error)
	CleanupDeletedUsers(ctx context.Context) (int, error)
}

type SessionCleanupService interface {
	CleanupExpiredSessions(ctx context.Context) (int, error)
	CleanupOrphanedSessions(ctx context.Context) (int, error)
}

type FileCleanupService interface {
	CleanupTempFiles(ctx context.Context, olderThanHours int) (int, error)
	CleanupOrphanedFiles(ctx context.Context) (int, error)
	CleanupLargeFiles(ctx context.Context, sizeLimitMB int) (int, error)
}

type LogCleanupService interface {
	CleanupOldLogs(ctx context.Context, retentionDays int) (int, error)
	CompressOldLogs(ctx context.Context, compressDays int) (int, error)
	RotateLogs(ctx context.Context) error
}

// NewDataCleanupJobHandler creates a new data cleanup job handler
func NewDataCleanupJobHandler(
	userService UserCleanupService,
	sessionService SessionCleanupService,
	fileService FileCleanupService,
	logService LogCleanupService,
	metrics jobs.JobMetrics,
) *DataCleanupJobHandler {
	return &DataCleanupJobHandler{
		userService:    userService,
		sessionService: sessionService,
		fileService:    fileService,
		logService:     logService,
		metrics:        metrics,
	}
}

// Execute processes a data cleanup job
func (h *DataCleanupJobHandler) Execute(ctx context.Context, job jobs.Job) error {
	start := time.Now()
	defer func() {
		if h.metrics != nil {
			h.metrics.ObserveJobDuration(job.GetType(), time.Since(start))
		}
	}()

	// Cast to DataCleanupJob
	cleanupJob, ok := job.(*domainjobs.DataCleanupJob)
	if !ok {
		err := fmt.Errorf("expected DataCleanupJob, got %T", job)
		if h.metrics != nil {
			h.metrics.IncrementJobsProcessed(job.GetType(), false)
		}
		return err
	}

	// Extract cleanup data
	payload := cleanupJob.GetPayload()
	cleanupType, _ := payload["cleanup_type"].(string)

	var err error
	var itemsProcessed int

	// Handle different cleanup types
	switch cleanupType {
	case "users":
		itemsProcessed, err = h.cleanupUsers(ctx, payload)
	case "sessions":
		itemsProcessed, err = h.cleanupSessions(ctx, payload)
	case "files":
		itemsProcessed, err = h.cleanupFiles(ctx, payload)
	case "logs":
		itemsProcessed, err = h.cleanupLogs(ctx, payload)
	case "full_system":
		itemsProcessed, err = h.performFullSystemCleanup(ctx, payload)
	case "database":
		itemsProcessed, err = h.cleanupDatabase(ctx, payload)
	default:
		err = fmt.Errorf("unknown cleanup type: %s", cleanupType)
	}

	// Record metrics
	if h.metrics != nil {
		success := err == nil
		h.metrics.IncrementJobsProcessed(job.GetType(), success)

		// Custom business metrics
		if businessMetrics, ok := h.metrics.(*BusinessCleanupMetrics); ok {
			businessMetrics.RecordCleanupJob(cleanupType, itemsProcessed, success)
		}
	}

	if err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	fmt.Printf("🧹 Cleanup completed: %s (%d items processed)\n", cleanupType, itemsProcessed)
	return nil
}

// GetJobType returns the job type this handler processes
func (h *DataCleanupJobHandler) GetJobType() string {
	return "data_cleanup"
}

// Private helper methods for different cleanup types

func (h *DataCleanupJobHandler) cleanupUsers(ctx context.Context, payload jobs.JobPayload) (int, error) {
	userCleanupType, _ := payload["user_cleanup_type"].(string)

	switch userCleanupType {
	case "inactive":
		days := getIntFromPayload(payload, "inactive_days", 90)
		return h.userService.CleanupInactiveUsers(ctx, days)
	case "unverified":
		days := getIntFromPayload(payload, "unverified_days", 30)
		return h.userService.CleanupUnverifiedUsers(ctx, days)
	case "deleted":
		return h.userService.CleanupDeletedUsers(ctx)
	default:
		// Clean up all user-related data
		total := 0

		// Clean inactive users (default 90 days)
		if count, err := h.userService.CleanupInactiveUsers(ctx, 90); err != nil {
			return 0, fmt.Errorf("failed to cleanup inactive users: %w", err)
		} else {
			total += count
		}

		// Clean unverified users (default 30 days)
		if count, err := h.userService.CleanupUnverifiedUsers(ctx, 30); err != nil {
			return total, fmt.Errorf("failed to cleanup unverified users: %w", err)
		} else {
			total += count
		}

		// Clean deleted users
		if count, err := h.userService.CleanupDeletedUsers(ctx); err != nil {
			return total, fmt.Errorf("failed to cleanup deleted users: %w", err)
		} else {
			total += count
		}

		return total, nil
	}
}

func (h *DataCleanupJobHandler) cleanupSessions(ctx context.Context, payload jobs.JobPayload) (int, error) {
	sessionCleanupType, _ := payload["session_cleanup_type"].(string)

	switch sessionCleanupType {
	case "expired":
		return h.sessionService.CleanupExpiredSessions(ctx)
	case "orphaned":
		return h.sessionService.CleanupOrphanedSessions(ctx)
	default:
		// Clean up all session-related data
		total := 0

		// Clean expired sessions
		if count, err := h.sessionService.CleanupExpiredSessions(ctx); err != nil {
			return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
		} else {
			total += count
		}

		// Clean orphaned sessions
		if count, err := h.sessionService.CleanupOrphanedSessions(ctx); err != nil {
			return total, fmt.Errorf("failed to cleanup orphaned sessions: %w", err)
		} else {
			total += count
		}

		return total, nil
	}
}

func (h *DataCleanupJobHandler) cleanupFiles(ctx context.Context, payload jobs.JobPayload) (int, error) {
	fileCleanupType, _ := payload["file_cleanup_type"].(string)

	switch fileCleanupType {
	case "temp":
		hours := getIntFromPayload(payload, "temp_file_hours", 24)
		return h.fileService.CleanupTempFiles(ctx, hours)
	case "orphaned":
		return h.fileService.CleanupOrphanedFiles(ctx)
	case "large":
		sizeLimitMB := getIntFromPayload(payload, "size_limit_mb", 100)
		return h.fileService.CleanupLargeFiles(ctx, sizeLimitMB)
	default:
		// Clean up all file-related data
		total := 0

		// Clean temp files (older than 24 hours)
		if count, err := h.fileService.CleanupTempFiles(ctx, 24); err != nil {
			return 0, fmt.Errorf("failed to cleanup temp files: %w", err)
		} else {
			total += count
		}

		// Clean orphaned files
		if count, err := h.fileService.CleanupOrphanedFiles(ctx); err != nil {
			return total, fmt.Errorf("failed to cleanup orphaned files: %w", err)
		} else {
			total += count
		}

		return total, nil
	}
}

func (h *DataCleanupJobHandler) cleanupLogs(ctx context.Context, payload jobs.JobPayload) (int, error) {
	logCleanupType, _ := payload["log_cleanup_type"].(string)

	switch logCleanupType {
	case "old":
		days := getIntFromPayload(payload, "retention_days", 30)
		return h.logService.CleanupOldLogs(ctx, days)
	case "compress":
		days := getIntFromPayload(payload, "compress_days", 7)
		count, err := h.logService.CompressOldLogs(ctx, days)
		return count, err
	case "rotate":
		err := h.logService.RotateLogs(ctx)
		return 1, err // Return 1 as a placeholder for rotation
	default:
		// Perform all log maintenance
		total := 0

		// Rotate logs first
		if err := h.logService.RotateLogs(ctx); err != nil {
			return 0, fmt.Errorf("failed to rotate logs: %w", err)
		}
		total++

		// Compress old logs (older than 7 days)
		if count, err := h.logService.CompressOldLogs(ctx, 7); err != nil {
			return total, fmt.Errorf("failed to compress old logs: %w", err)
		} else {
			total += count
		}

		// Clean old logs (older than 30 days)
		if count, err := h.logService.CleanupOldLogs(ctx, 30); err != nil {
			return total, fmt.Errorf("failed to cleanup old logs: %w", err)
		} else {
			total += count
		}

		return total, nil
	}
}

func (h *DataCleanupJobHandler) performFullSystemCleanup(ctx context.Context, payload jobs.JobPayload) (int, error) {
	total := 0

	// User cleanup
	if count, err := h.cleanupUsers(ctx, map[string]interface{}{}); err != nil {
		return total, fmt.Errorf("user cleanup failed: %w", err)
	} else {
		total += count
	}

	// Session cleanup
	if count, err := h.cleanupSessions(ctx, map[string]interface{}{}); err != nil {
		return total, fmt.Errorf("session cleanup failed: %w", err)
	} else {
		total += count
	}

	// File cleanup
	if count, err := h.cleanupFiles(ctx, map[string]interface{}{}); err != nil {
		return total, fmt.Errorf("file cleanup failed: %w", err)
	} else {
		total += count
	}

	// Log cleanup
	if count, err := h.cleanupLogs(ctx, map[string]interface{}{}); err != nil {
		return total, fmt.Errorf("log cleanup failed: %w", err)
	} else {
		total += count
	}

	return total, nil
}

func (h *DataCleanupJobHandler) cleanupDatabase(ctx context.Context, payload jobs.JobPayload) (int, error) {
	// This could include database-specific cleanup operations
	// like optimizing tables, updating statistics, etc.

	// Simulate database cleanup operations
	time.Sleep(2 * time.Second)

	// In a real implementation, this would:
	// - Optimize database tables
	// - Update table statistics
	// - Clean up database logs
	// - Defragment indexes

	return 10, nil // Placeholder for database operations performed
}

// DataCleanupJobFactory creates data cleanup jobs with proper validation
type DataCleanupJobFactory struct{}

// NewDataCleanupJobFactory creates a new data cleanup job factory
func NewDataCleanupJobFactory() *DataCleanupJobFactory {
	return &DataCleanupJobFactory{}
}

// CreateUserCleanupJob creates a user cleanup job
func (f *DataCleanupJobFactory) CreateUserCleanupJob(cleanupType string, params map[string]interface{}) (jobs.Job, error) {
	payload := jobs.JobPayload{
		"cleanup_type":      "users",
		"user_cleanup_type": cleanupType,
	}

	// Add type-specific parameters
	for key, value := range params {
		payload[key] = value
	}

	opts := jobs.JobOptions{
		Priority: jobs.PriorityLow, // Cleanup jobs typically have lower priority
		Queue:    "cleanup",
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJobWithOptions("data_cleanup", payload, opts)
}

// CreateScheduledCleanupJob creates a recurring cleanup job
func (f *DataCleanupJobFactory) CreateScheduledCleanupJob(cleanupType string) (jobs.Job, error) {
	payload := jobs.JobPayload{
		"cleanup_type": cleanupType,
	}

	opts := jobs.JobOptions{
		Priority: jobs.PriorityLow,
		Queue:    "scheduled_cleanup",
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJobWithOptions("data_cleanup", payload, opts)
}

// CreateFullSystemCleanupJob creates a comprehensive system cleanup job
func (f *DataCleanupJobFactory) CreateFullSystemCleanupJob() (jobs.Job, error) {
	payload := jobs.JobPayload{
		"cleanup_type": "full_system",
	}

	opts := jobs.JobOptions{
		Priority: jobs.PriorityLow,
		Queue:    "system_maintenance",
	}

	factory := domainjobs.NewJobFactory()
	return factory.CreateJobWithOptions("data_cleanup", payload, opts)
}

// Mock implementations for demonstration

// MockUserCleanupService simulates user cleanup operations
type MockUserCleanupService struct{}

func NewMockUserCleanupService() *MockUserCleanupService {
	return &MockUserCleanupService{}
}

func (m *MockUserCleanupService) CleanupInactiveUsers(ctx context.Context, inactiveDays int) (int, error) {
	time.Sleep(500 * time.Millisecond)
	count := 15 // Simulate finding 15 inactive users
	fmt.Printf("🧹 Cleaned up %d inactive users (inactive > %d days)\n", count, inactiveDays)
	return count, nil
}

func (m *MockUserCleanupService) CleanupUnverifiedUsers(ctx context.Context, unverifiedDays int) (int, error) {
	time.Sleep(300 * time.Millisecond)
	count := 8 // Simulate finding 8 unverified users
	fmt.Printf("🧹 Cleaned up %d unverified users (unverified > %d days)\n", count, unverifiedDays)
	return count, nil
}

func (m *MockUserCleanupService) CleanupDeletedUsers(ctx context.Context) (int, error) {
	time.Sleep(200 * time.Millisecond)
	count := 3 // Simulate finding 3 soft-deleted users
	fmt.Printf("🧹 Cleaned up %d soft-deleted users\n", count)
	return count, nil
}

// MockSessionCleanupService simulates session cleanup operations
type MockSessionCleanupService struct{}

func NewMockSessionCleanupService() *MockSessionCleanupService {
	return &MockSessionCleanupService{}
}

func (m *MockSessionCleanupService) CleanupExpiredSessions(ctx context.Context) (int, error) {
	time.Sleep(400 * time.Millisecond)
	count := 45 // Simulate finding 45 expired sessions
	fmt.Printf("🧹 Cleaned up %d expired sessions\n", count)
	return count, nil
}

func (m *MockSessionCleanupService) CleanupOrphanedSessions(ctx context.Context) (int, error) {
	time.Sleep(300 * time.Millisecond)
	count := 12 // Simulate finding 12 orphaned sessions
	fmt.Printf("🧹 Cleaned up %d orphaned sessions\n", count)
	return count, nil
}

// MockFileCleanupService simulates file cleanup operations
type MockFileCleanupService struct{}

func NewMockFileCleanupService() *MockFileCleanupService {
	return &MockFileCleanupService{}
}

func (m *MockFileCleanupService) CleanupTempFiles(ctx context.Context, olderThanHours int) (int, error) {
	time.Sleep(600 * time.Millisecond)
	count := 23 // Simulate finding 23 temp files
	fmt.Printf("🧹 Cleaned up %d temp files (older than %d hours)\n", count, olderThanHours)
	return count, nil
}

func (m *MockFileCleanupService) CleanupOrphanedFiles(ctx context.Context) (int, error) {
	time.Sleep(800 * time.Millisecond)
	count := 7 // Simulate finding 7 orphaned files
	fmt.Printf("🧹 Cleaned up %d orphaned files\n", count)
	return count, nil
}

func (m *MockFileCleanupService) CleanupLargeFiles(ctx context.Context, sizeLimitMB int) (int, error) {
	time.Sleep(1000 * time.Millisecond)
	count := 4 // Simulate finding 4 large files
	fmt.Printf("🧹 Cleaned up %d large files (> %d MB)\n", count, sizeLimitMB)
	return count, nil
}

// MockLogCleanupService simulates log cleanup operations
type MockLogCleanupService struct{}

func NewMockLogCleanupService() *MockLogCleanupService {
	return &MockLogCleanupService{}
}

func (m *MockLogCleanupService) CleanupOldLogs(ctx context.Context, retentionDays int) (int, error) {
	time.Sleep(400 * time.Millisecond)
	count := 156 // Simulate finding 156 old log entries
	fmt.Printf("🧹 Cleaned up %d old log entries (older than %d days)\n", count, retentionDays)
	return count, nil
}

func (m *MockLogCleanupService) CompressOldLogs(ctx context.Context, compressDays int) (int, error) {
	time.Sleep(700 * time.Millisecond)
	count := 34 // Simulate compressing 34 log files
	fmt.Printf("🧹 Compressed %d log files (older than %d days)\n", count, compressDays)
	return count, nil
}

func (m *MockLogCleanupService) RotateLogs(ctx context.Context) error {
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("🧹 Log rotation completed\n")
	return nil
}

// BusinessCleanupMetrics extends the basic metrics with cleanup-specific metrics
type BusinessCleanupMetrics struct {
	jobs.JobMetrics
	cleanupCounts  map[string]int64 // cleanup_type -> count
	itemsProcessed map[string]int64 // cleanup_type -> total items processed
}

// NewBusinessCleanupMetrics creates enhanced cleanup metrics
func NewBusinessCleanupMetrics(baseMetrics jobs.JobMetrics) *BusinessCleanupMetrics {
	return &BusinessCleanupMetrics{
		JobMetrics:     baseMetrics,
		cleanupCounts:  make(map[string]int64),
		itemsProcessed: make(map[string]int64),
	}
}

// RecordCleanupJob records cleanup-specific metrics
func (m *BusinessCleanupMetrics) RecordCleanupJob(cleanupType string, itemsProcessed int, success bool) {
	if success {
		m.cleanupCounts[cleanupType]++
		m.itemsProcessed[cleanupType] += int64(itemsProcessed)
	}
}

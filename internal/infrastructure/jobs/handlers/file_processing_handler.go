package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/job"
)

type FileProcessingJobHandler struct {
	fileService    FileProcessingService
	storageService StorageService
	metrics        job.JobMetrics
}

type FileProcessingService interface {
	ProcessImage(ctx context.Context, inputPath, outputPath string, options ImageProcessingOptions) error
	ProcessVideo(ctx context.Context, inputPath, outputPath string, options VideoProcessingOptions) error
	ProcessDocument(ctx context.Context, inputPath, outputPath string, options DocumentProcessingOptions) error
	ValidateFile(ctx context.Context, filePath string, fileType string) error
}

type StorageService interface {
	UploadFile(ctx context.Context, filePath, destination string) error
	DownloadFile(ctx context.Context, source, destination string) error
	DeleteFile(ctx context.Context, filePath string) error
	GetFileInfo(ctx context.Context, filePath string) (*FileInfo, error)
}

type ImageProcessingOptions struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Quality  int    `json:"quality"`
	Format   string `json:"format"`
	Optimize bool   `json:"optimize"`
}

type VideoProcessingOptions struct {
	Resolution string `json:"resolution"`
	Bitrate    string `json:"bitrate"`
	Format     string `json:"format"`
	Compress   bool   `json:"compress"`
}

type DocumentProcessingOptions struct {
	ConvertTo   string `json:"convert_to"`
	ExtractText bool   `json:"extract_text"`
	GeneratePDF bool   `json:"generate_pdf"`
	Watermark   string `json:"watermark"`
}

type FileInfo struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFileProcessingJobHandler(fileService FileProcessingService, storageService StorageService, metrics job.JobMetrics) *FileProcessingJobHandler {
	return &FileProcessingJobHandler{
		fileService:    fileService,
		storageService: storageService,
		metrics:        metrics,
	}
}

func (h *FileProcessingJobHandler) Execute(ctx context.Context, excutedJob job.Job) error {
	start := time.Now()
	defer func() {
		if h.metrics != nil {
			h.metrics.ObserveJobDuration(excutedJob.GetType(), time.Since(start))
		}
	}()

	fileJob, ok := excutedJob.(*job.FileProcessingJob)
	if !ok {
		err := fmt.Errorf("expected FileProcessingJob, got %T", excutedJob)
		if h.metrics != nil {
			h.metrics.IncrementJobsProcessed(excutedJob.GetType(), false)
		}
		return err
	}

	payload := fileJob.GetPayload()
	processingType, _ := payload["processing_type"].(string)
	inputPath, _ := payload["input_path"].(string)
	outputPath, _ := payload["output_path"].(string)
	fileType, _ := payload["file_type"].(string)

	if inputPath == "" {
		err := fmt.Errorf("input_path is required")
		if h.metrics != nil {
			h.metrics.IncrementJobsProcessed(excutedJob.GetType(), false)
		}
		return err
	}

	var err error

	switch processingType {
	case "image":
		err = h.processImage(ctx, inputPath, outputPath, payload)
	case "video":
		err = h.processVideo(ctx, inputPath, outputPath, payload)
	case "document":
		err = h.processDocument(ctx, inputPath, outputPath, payload)
	case "upload":
		err = h.handleUpload(ctx, inputPath, payload)
	case "validation":
		err = h.handleValidation(ctx, inputPath, fileType)
	case "batch":
		err = h.handleBatchProcessing(ctx, payload)
	default:
		err = fmt.Errorf("unknown processing type: %s", processingType)
	}

	if h.metrics != nil {
		success := err == nil
		h.metrics.IncrementJobsProcessed(excutedJob.GetType(), success)

		if businessMetrics, ok := h.metrics.(*BusinessFileMetrics); ok {
			businessMetrics.RecordFileProcessingJob(processingType, fileType, success)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	return nil
}

func (h *FileProcessingJobHandler) GetJobType() string {
	return "file_processing"
}

func (h *FileProcessingJobHandler) processImage(ctx context.Context, inputPath, outputPath string, payload job.JobPayload) error {
	options := ImageProcessingOptions{
		Width:    getIntFromPayload(payload, "width", 0),
		Height:   getIntFromPayload(payload, "height", 0),
		Quality:  getIntFromPayload(payload, "quality", 90),
		Format:   getStringFromPayload(payload, "format", "jpg"),
		Optimize: getBoolFromPayload(payload, "optimize", true),
	}

	if err := h.fileService.ValidateFile(ctx, inputPath, "image"); err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}

	return h.fileService.ProcessImage(ctx, inputPath, outputPath, options)
}

func (h *FileProcessingJobHandler) processVideo(ctx context.Context, inputPath, outputPath string, payload job.JobPayload) error {
	options := VideoProcessingOptions{
		Resolution: getStringFromPayload(payload, "resolution", "720p"),
		Bitrate:    getStringFromPayload(payload, "bitrate", "1000k"),
		Format:     getStringFromPayload(payload, "format", "mp4"),
		Compress:   getBoolFromPayload(payload, "compress", true),
	}

	if err := h.fileService.ValidateFile(ctx, inputPath, "video"); err != nil {
		return fmt.Errorf("video validation failed: %w", err)
	}

	return h.fileService.ProcessVideo(ctx, inputPath, outputPath, options)
}

func (h *FileProcessingJobHandler) processDocument(ctx context.Context, inputPath, outputPath string, payload job.JobPayload) error {
	options := DocumentProcessingOptions{
		ConvertTo:   getStringFromPayload(payload, "convert_to", "pdf"),
		ExtractText: getBoolFromPayload(payload, "extract_text", false),
		GeneratePDF: getBoolFromPayload(payload, "generate_pdf", false),
		Watermark:   getStringFromPayload(payload, "watermark", ""),
	}

	if err := h.fileService.ValidateFile(ctx, inputPath, "document"); err != nil {
		return fmt.Errorf("document validation failed: %w", err)
	}

	return h.fileService.ProcessDocument(ctx, inputPath, outputPath, options)
}

func (h *FileProcessingJobHandler) handleUpload(ctx context.Context, inputPath string, payload job.JobPayload) error {
	destination, _ := payload["destination"].(string)
	if destination == "" {
		return fmt.Errorf("destination is required for upload")
	}

	return h.storageService.UploadFile(ctx, inputPath, destination)
}

func (h *FileProcessingJobHandler) handleValidation(ctx context.Context, inputPath, fileType string) error {
	return h.fileService.ValidateFile(ctx, inputPath, fileType)
}

func (h *FileProcessingJobHandler) handleBatchProcessing(ctx context.Context, payload job.JobPayload) error {
	inputPaths, ok := payload["input_paths"].([]string)
	if !ok {
		return fmt.Errorf("input_paths must be a string array for batch processing")
	}

	processingType, _ := payload["batch_processing_type"].(string)
	outputDir, _ := payload["output_directory"].(string)

	for _, inputPath := range inputPaths {
		outputPath := filepath.Join(outputDir, filepath.Base(inputPath))

		switch processingType {
		case "image":
			if err := h.processImage(ctx, inputPath, outputPath, payload); err != nil {
				return fmt.Errorf("failed to process image %s: %w", inputPath, err)
			}
		case "video":
			if err := h.processVideo(ctx, inputPath, outputPath, payload); err != nil {
				return fmt.Errorf("failed to process video %s: %w", inputPath, err)
			}
		case "document":
			if err := h.processDocument(ctx, inputPath, outputPath, payload); err != nil {
				return fmt.Errorf("failed to process document %s: %w", inputPath, err)
			}
		default:
			return fmt.Errorf("unsupported batch processing type: %s", processingType)
		}
	}

	return nil
}

func getStringFromPayload(payload job.JobPayload, key, defaultValue string) string {
	if val, ok := payload[key].(string); ok {
		return val
	}
	return defaultValue
}

func getIntFromPayload(payload job.JobPayload, key string, defaultValue int) int {
	if val, ok := payload[key].(int); ok {
		return val
	}
	if val, ok := payload[key].(float64); ok {
		return int(val)
	}
	return defaultValue
}

func getBoolFromPayload(payload job.JobPayload, key string, defaultValue bool) bool {
	if val, ok := payload[key].(bool); ok {
		return val
	}
	return defaultValue
}

type FileProcessingJobFactory struct{}

func NewFileProcessingJobFactory() *FileProcessingJobFactory {
	return &FileProcessingJobFactory{}
}

func (f *FileProcessingJobFactory) CreateImageProcessingJob(inputPath, outputPath string, options ImageProcessingOptions) (job.Job, error) {
	payload := job.JobPayload{
		"processing_type": "image",
		"input_path":      inputPath,
		"output_path":     outputPath,
		"width":           options.Width,
		"height":          options.Height,
		"quality":         options.Quality,
		"format":          options.Format,
		"optimize":        options.Optimize,
	}

	opts := job.JobOptions{
		Priority: job.PriorityNormal,
		Queue:    "file_processing",
	}

	factory := job.NewJobFactory()
	return factory.CreateJobWithOptions("file_processing", payload, opts)
}

func (f *FileProcessingJobFactory) CreateBatchProcessingJob(inputPaths []string, outputDir, processingType string) (job.Job, error) {
	payload := job.JobPayload{
		"processing_type":       "batch",
		"input_paths":           inputPaths,
		"output_directory":      outputDir,
		"batch_processing_type": processingType,
	}

	opts := job.JobOptions{
		Priority: job.PriorityLow, // Batch jobs have lower priority
		Queue:    "batch_processing",
	}

	factory := job.NewJobFactory()
	return factory.CreateJobWithOptions("file_processing", payload, opts)
}

type MockFileProcessingService struct{}

func NewMockFileProcessingService() *MockFileProcessingService {
	return &MockFileProcessingService{}
}

func (m *MockFileProcessingService) ProcessImage(ctx context.Context, inputPath, outputPath string, options ImageProcessingOptions) error {
	processingTime := 500 * time.Millisecond
	if options.Width > 1000 || options.Height > 1000 {
		processingTime = 2 * time.Second
	}

	time.Sleep(processingTime)

	if strings.Contains(inputPath, "corrupt") {
		return fmt.Errorf("image file appears to be corrupt: %s", inputPath)
	}

	fmt.Printf("🖼️  Image processed: %s -> %s (quality: %d%%)\n", inputPath, outputPath, options.Quality)
	return nil
}

func (m *MockFileProcessingService) ProcessVideo(ctx context.Context, inputPath, outputPath string, options VideoProcessingOptions) error {
	time.Sleep(3 * time.Second)

	if strings.Contains(inputPath, "unsupported") {
		return fmt.Errorf("unsupported video format: %s", inputPath)
	}

	fmt.Printf("🎥 Video processed: %s -> %s (%s)\n", inputPath, outputPath, options.Resolution)
	return nil
}

func (m *MockFileProcessingService) ProcessDocument(ctx context.Context, inputPath, outputPath string, options DocumentProcessingOptions) error {
	time.Sleep(800 * time.Millisecond)

	fmt.Printf("📄 Document processed: %s -> %s (format: %s)\n", inputPath, outputPath, options.ConvertTo)
	return nil
}

func (m *MockFileProcessingService) ValidateFile(ctx context.Context, filePath, fileType string) error {
	time.Sleep(100 * time.Millisecond)

	if strings.Contains(filePath, "missing") {
		return fmt.Errorf("file not found: %s", filePath)
	}

	if !m.isValidFileType(filePath, fileType) {
		return fmt.Errorf("invalid file type for %s processing: %s", fileType, filePath)
	}

	return nil
}

func (m *MockFileProcessingService) isValidFileType(filePath, fileType string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch fileType {
	case "image":
		return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp"
	case "video":
		return ext == ".mp4" || ext == ".avi" || ext == ".mkv" || ext == ".mov" || ext == ".wmv"
	case "document":
		return ext == ".pdf" || ext == ".doc" || ext == ".docx" || ext == ".txt" || ext == ".rtf"
	default:
		return true
	}
}

type MockStorageService struct{}

func NewMockStorageService() *MockStorageService {
	return &MockStorageService{}
}

func (m *MockStorageService) UploadFile(ctx context.Context, filePath, destination string) error {
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("☁️  File uploaded: %s -> %s\n", filePath, destination)
	return nil
}

func (m *MockStorageService) DownloadFile(ctx context.Context, source, destination string) error {
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("⬇️  File downloaded: %s -> %s\n", source, destination)
	return nil
}

func (m *MockStorageService) DeleteFile(ctx context.Context, filePath string) error {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("🗑️  File deleted: %s\n", filePath)
	return nil
}

func (m *MockStorageService) GetFileInfo(ctx context.Context, filePath string) (*FileInfo, error) {
	time.Sleep(50 * time.Millisecond)

	return &FileInfo{
		Name:      filepath.Base(filePath),
		Size:      1024 * 1024, // 1MB simulation
		MimeType:  "application/octet-stream",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}, nil
}

type BusinessFileMetrics struct {
	job.JobMetrics
	fileProcessingCounts map[string]map[string]int64 // processing_type -> file_type -> count
}

func NewBusinessFileMetrics(baseMetrics job.JobMetrics) *BusinessFileMetrics {
	return &BusinessFileMetrics{
		JobMetrics:           baseMetrics,
		fileProcessingCounts: make(map[string]map[string]int64),
	}
}

func (m *BusinessFileMetrics) RecordFileProcessingJob(processingType, fileType string, success bool) {
	if success {
		if m.fileProcessingCounts[processingType] == nil {
			m.fileProcessingCounts[processingType] = make(map[string]int64)
		}
		m.fileProcessingCounts[processingType][fileType]++
	}
}

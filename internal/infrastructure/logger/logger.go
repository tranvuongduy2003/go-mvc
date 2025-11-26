package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
)

type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

func NewLogger(cfg config.Logger) (*Logger, error) {
	zapConfig := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	zapConfig.Encoding = cfg.Encoding
	if cfg.Development {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(level)
	}

	zapConfig.OutputPaths = cfg.OutputPaths
	zapConfig.ErrorOutputPaths = cfg.ErrorPaths

	zapLogger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(zap.Any(key, value)),
		sugar:  l.sugar.With(key, value),
	}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{
		Logger: l.Logger.With(zapFields...),
		sugar:  l.sugar.With(fields),
	}
}

func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Logger.Panic(msg, fields...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

var global *Logger

func SetGlobal(l *Logger) {
	global = l
	zap.ReplaceGlobals(l.Logger)
}

func Global() *Logger {
	if global == nil {
		zapLogger, _ := zap.NewProduction()
		global = &Logger{
			Logger: zapLogger,
			sugar:  zapLogger.Sugar(),
		}
	}
	return global
}

func InitDefault() *Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	l := &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
	}

	SetGlobal(l)
	return l
}

const (
	RequestIDKey  = "request_id"
	UserIDKey     = "user_id"
	TraceIDKey    = "trace_id"
	ComponentKey  = "component"
	OperationKey  = "operation"
	DurationKey   = "duration"
	StatusCodeKey = "status_code"
)

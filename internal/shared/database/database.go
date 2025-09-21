package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	appLogger "github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// Manager manages database connections
type Manager struct {
	primary *gorm.DB
	replica *gorm.DB
	logger  *appLogger.Logger
}

// NewManager creates a new database manager
func NewManager(cfg config.Database, log *appLogger.Logger) (*Manager, error) {
	manager := &Manager{
		logger: log,
	}

	// Setup primary database connection
	primary, err := setupConnection(cfg.Primary, log, "primary")
	if err != nil {
		return nil, fmt.Errorf("failed to setup primary database: %w", err)
	}
	manager.primary = primary

	// Setup replica database connection if configured
	if cfg.Replica.Host != "" {
		replica, err := setupConnection(cfg.Replica, log, "replica")
		if err != nil {
			log.Warnf("Failed to setup replica database: %v", err)
			// Use primary as fallback for replica
			manager.replica = primary
		} else {
			manager.replica = replica
		}
	} else {
		// Use primary as replica if no replica is configured
		manager.replica = primary
	}

	return manager, nil
}

// setupConnection creates a database connection
func setupConnection(cfg config.DatabaseConnection, log *appLogger.Logger, name string) (*gorm.DB, error) {
	// Parse log level
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}

	// Create GORM config
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s database: %w", name, err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB from %s database: %w", name, err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping %s database: %w", name, err)
	}

	log.Infof("%s database connected successfully", name)
	return db, nil
}

// Primary returns the primary database connection
func (m *Manager) Primary() *gorm.DB {
	return m.primary
}

// Replica returns the replica database connection
func (m *Manager) Replica() *gorm.DB {
	return m.replica
}

// Close closes all database connections
func (m *Manager) Close() error {
	var errs []error

	if m.primary != nil {
		if sqlDB, err := m.primary.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close primary database: %w", err))
			}
		}
	}

	if m.replica != nil && m.replica != m.primary {
		if sqlDB, err := m.replica.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close replica database: %w", err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing databases: %v", errs)
	}

	m.logger.Info("All database connections closed successfully")
	return nil
}

// Transaction executes a function within a database transaction
func (m *Manager) Transaction(fn func(*gorm.DB) error) error {
	return m.primary.Transaction(fn)
}

// Health checks the health of database connections
func (m *Manager) Health() error {
	// Check primary database
	if sqlDB, err := m.primary.DB(); err != nil {
		return fmt.Errorf("failed to get primary database: %w", err)
	} else if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("primary database ping failed: %w", err)
	}

	// Check replica database if it's different from primary
	if m.replica != m.primary {
		if sqlDB, err := m.replica.DB(); err != nil {
			return fmt.Errorf("failed to get replica database: %w", err)
		} else if err := sqlDB.Ping(); err != nil {
			return fmt.Errorf("replica database ping failed: %w", err)
		}
	}

	return nil
}

// Migrate runs database migrations
func (m *Manager) Migrate(models ...interface{}) error {
	if err := m.primary.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	m.logger.Info("Database migration completed successfully")
	return nil
}

// Stats returns database connection statistics
func (m *Manager) Stats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Primary database stats
	if sqlDB, err := m.primary.DB(); err == nil {
		dbStats := sqlDB.Stats()
		stats["primary"] = map[string]interface{}{
			"open_connections":     dbStats.OpenConnections,
			"in_use":               dbStats.InUse,
			"idle":                 dbStats.Idle,
			"wait_count":           dbStats.WaitCount,
			"wait_duration":        dbStats.WaitDuration.String(),
			"max_idle_closed":      dbStats.MaxIdleClosed,
			"max_idle_time_closed": dbStats.MaxIdleTimeClosed,
			"max_lifetime_closed":  dbStats.MaxLifetimeClosed,
		}
	}

	// Replica database stats if different
	if m.replica != m.primary {
		if sqlDB, err := m.replica.DB(); err == nil {
			dbStats := sqlDB.Stats()
			stats["replica"] = map[string]interface{}{
				"open_connections":     dbStats.OpenConnections,
				"in_use":               dbStats.InUse,
				"idle":                 dbStats.Idle,
				"wait_count":           dbStats.WaitCount,
				"wait_duration":        dbStats.WaitDuration.String(),
				"max_idle_closed":      dbStats.MaxIdleClosed,
				"max_idle_time_closed": dbStats.MaxIdleTimeClosed,
				"max_lifetime_closed":  dbStats.MaxLifetimeClosed,
			}
		}
	}

	return stats, nil
}

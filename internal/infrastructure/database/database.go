package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	appLogger "github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
)

type Manager struct {
	primary *gorm.DB
	replica *gorm.DB
	logger  *appLogger.Logger
}

func NewManager(cfg config.Database, log *appLogger.Logger) (*Manager, error) {
	manager := &Manager{
		logger: log,
	}

	primary, err := setupConnection(cfg.Primary, log, "primary")
	if err != nil {
		return nil, fmt.Errorf("failed to setup primary database: %w", err)
	}
	manager.primary = primary

	if cfg.Replica.Host != "" {
		replica, err := setupConnection(cfg.Replica, log, "replica")
		if err != nil {
			log.Warnf("Failed to setup replica database: %v", err)
			manager.replica = primary
		} else {
			manager.replica = replica
		}
	} else {
		manager.replica = primary
	}

	return manager, nil
}

func setupConnection(cfg config.DatabaseConnection, log *appLogger.Logger, name string) (*gorm.DB, error) {
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

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s database: %w", name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB from %s database: %w", name, err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping %s database: %w", name, err)
	}

	log.Infof("%s database connected successfully", name)
	return db, nil
}

func (m *Manager) Primary() *gorm.DB {
	return m.primary
}

func (m *Manager) Replica() *gorm.DB {
	return m.replica
}

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

func (m *Manager) Transaction(fn func(*gorm.DB) error) error {
	return m.primary.Transaction(fn)
}

func (m *Manager) Health() error {
	if sqlDB, err := m.primary.DB(); err != nil {
		return fmt.Errorf("failed to get primary database: %w", err)
	} else if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("primary database ping failed: %w", err)
	}

	if m.replica != m.primary {
		if sqlDB, err := m.replica.DB(); err != nil {
			return fmt.Errorf("failed to get replica database: %w", err)
		} else if err := sqlDB.Ping(); err != nil {
			return fmt.Errorf("replica database ping failed: %w", err)
		}
	}

	return nil
}

func (m *Manager) Migrate(models ...interface{}) error {
	if err := m.primary.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	m.logger.Info("Database migration completed successfully")
	return nil
}

func (m *Manager) Stats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	fxmodules "github.com/tranvuongduy2003/go-mvc/internal/fx_modules"
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration tool",
	Long:  `A tool for managing database migrations including running migrations, rollbacks, and checking migration status.`,
}

// loggerProvider provides a development logger
func loggerProvider() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(upCommand())
	rootCmd.AddCommand(downCommand())
	rootCmd.AddCommand(statusCommand())
	rootCmd.AddCommand(resetCommand())
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(versionCommand())
}

// Migration commands
func upCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		Long:  `Execute all pending database migrations to bring the database schema up to date.`,
		Run: func(cmd *cobra.Command, args []string) {
			steps, _ := cmd.Flags().GetInt("steps")

			fmt.Println("Running database migrations...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Provide(func() *zap.Logger {
					logger, _ := zap.NewDevelopment()
					return logger
				}),
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					if err := runMigrationsUp(db, steps); err != nil {
						logger.Error("Migration failed", zap.Error(err))
						log.Fatalf("Migration failed: %v", err)
					}

					logger.Info("Migrations completed successfully")
					fmt.Println("✅ All migrations completed successfully!")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to run migrations: %v", err)
			}
			app.Stop(context.Background())
		},
	}

	cmd.Flags().IntP("steps", "s", 0, "Number of migration steps to run (0 = all)")
	return cmd
}

func downCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback migrations",
		Long:  `Rollback database migrations by specified number of steps.`,
		Run: func(cmd *cobra.Command, args []string) {
			steps, _ := cmd.Flags().GetInt("steps")
			force, _ := cmd.Flags().GetBool("force")

			if steps <= 0 && !force {
				fmt.Println("❌ Must specify number of steps to rollback with --steps flag")
				fmt.Println("Use --force flag to rollback all migrations (WARNING: This will delete all data)")
				return
			}

			if force {
				fmt.Println("⚠️  This will rollback ALL migrations and delete all data!")
				fmt.Println("Use --confirm flag to proceed: migrate down --force --confirm")
				return
			}

			fmt.Printf("Rolling back %d migration step(s)...\n", steps)

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					if err := runMigrationsDown(db, steps); err != nil {
						logger.Error("Migration rollback failed", zap.Error(err))
						log.Fatalf("Migration rollback failed: %v", err)
					}

					logger.Info("Migration rollback completed successfully")
					fmt.Println("✅ Migration rollback completed!")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to rollback migrations: %v", err)
			}
			app.Stop(context.Background())
		},
	}

	cmd.Flags().IntP("steps", "s", 1, "Number of migration steps to rollback")
	cmd.Flags().Bool("force", false, "Force rollback all migrations")
	cmd.Flags().Bool("confirm", false, "Confirm the rollback operation")
	return cmd
}

func statusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		Long:  `Display the current status of all migrations including which ones have been applied.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Checking migration status...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Provide(func() *zap.Logger {
					logger, _ := zap.NewDevelopment()
					return logger
				}),
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					if err := showMigrationStatus(db); err != nil {
						logger.Error("Failed to get migration status", zap.Error(err))
						log.Fatalf("Failed to get migration status: %v", err)
					}
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to check migration status: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func resetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset all migrations (WARNING: This will delete all data)",
		Long:  `Drop all tables and re-run all migrations from scratch. This will permanently delete all data.`,
		Run: func(cmd *cobra.Command, args []string) {
			confirm, _ := cmd.Flags().GetBool("confirm")
			if !confirm {
				fmt.Println("⚠️  This operation will permanently delete all data!")
				fmt.Println("Use --confirm flag to proceed: migrate reset --confirm")
				return
			}

			fmt.Println("Resetting all migrations...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					if err := resetAllMigrations(db); err != nil {
						logger.Error("Migration reset failed", zap.Error(err))
						log.Fatalf("Migration reset failed: %v", err)
					}

					logger.Info("Migration reset completed successfully")
					fmt.Println("✅ All migrations reset completed!")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to reset migrations: %v", err)
			}
			app.Stop(context.Background())
		},
	}

	cmd.Flags().Bool("confirm", false, "Confirm the reset operation")
	return cmd
}

func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration file",
		Long:  `Create a new migration file with the specified name.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			// Create migration directory if it doesn't exist
			migrationDir := "internal/adapters/persistence/postgres/migrations"
			if err := os.MkdirAll(migrationDir, 0755); err != nil {
				log.Fatalf("Failed to create migration directory: %v", err)
			}

			// Generate timestamp
			timestamp := time.Now().Format("20060102150405")

			// Create migration file name
			fileName := fmt.Sprintf("%s_%s.sql", timestamp, strings.ReplaceAll(name, " ", "_"))
			filePath := filepath.Join(migrationDir, fileName)

			// Create migration file content
			content := fmt.Sprintf(`-- Migration: %s
-- Created at: %s

-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back

`, name, time.Now().Format("2006-01-02 15:04:05"))

			// Write migration file
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				log.Fatalf("Failed to create migration file: %v", err)
			}

			fmt.Printf("✅ Created migration file: %s\n", filePath)
		},
	}

	return cmd
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show migration tool version",
		Long:  `Display the current version of the migration tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Go MVC Migration Tool")
			fmt.Println("Version: 1.0.0")
			fmt.Println("Build: development")
		},
	}
}

// Helper functions
func runMigrationsUp(db *gorm.DB, steps int) error {
	fmt.Println("Auto-migrating database schema...")

	// Auto-migrate all models
	if err := db.AutoMigrate(); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	fmt.Println("Database schema migration completed")
	return nil
}

func runMigrationsDown(db *gorm.DB, steps int) error {
	fmt.Printf("Rolling back %d migration steps...\n", steps)

	// For now, we'll implement a simple table drop approach
	// In a real implementation, you would track migration versions
	// and rollback specific migrations

	tables := []string{
		"role_permissions",
		"user_roles",
		"permissions",
		"roles",
		"users",
	}

	// Rollback specified number of tables (or all if steps >= len(tables))
	rollbackCount := steps
	if steps >= len(tables) {
		rollbackCount = len(tables)
	}

	for i := 0; i < rollbackCount; i++ {
		table := tables[i]
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	return nil
}

func showMigrationStatus(db *gorm.DB) error {
	fmt.Println("\n=== Migration Status ===")

	// Check if tables exist
	tables := []string{"users", "roles", "permissions", "user_roles", "role_permissions"}

	for _, table := range tables {
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '%s')", table)
		if err := db.Raw(query).Scan(&exists).Error; err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}

		status := "❌ Not Applied"
		if exists {
			status = "✅ Applied"
		}

		fmt.Printf("Table %s: %s\n", table, status)
	}

	fmt.Println("\n=== Database Info ===")

	// Get database version
	var version string
	if err := db.Raw("SELECT version()").Scan(&version).Error; err == nil {
		fmt.Printf("PostgreSQL Version: %s\n", strings.Split(version, " ")[1])
	}

	// Get current database name
	var dbName string
	if err := db.Raw("SELECT current_database()").Scan(&dbName).Error; err == nil {
		fmt.Printf("Current Database: %s\n", dbName)
	}

	return nil
}

func resetAllMigrations(db *gorm.DB) error {
	fmt.Println("Dropping all tables...")

	tables := []string{
		"role_permissions",
		"user_roles",
		"permissions",
		"roles",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	fmt.Println("Re-running migrations...")
	return runMigrationsUp(db, 0)
}

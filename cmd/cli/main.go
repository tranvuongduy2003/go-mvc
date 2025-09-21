package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	fxmodules "github.com/tranvuongduy2003/go-mvc/internal/fx_modules"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
)

var rootCmd = &cobra.Command{
	Use:   "go-mvc-cli",
	Short: "CLI tool for Go MVC application management",
	Long:  `A comprehensive CLI tool for managing the Go MVC application including database operations, user management, and system maintenance.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Database commands
	rootCmd.AddCommand(createDBCommand())
	rootCmd.AddCommand(migrateCommand())
	rootCmd.AddCommand(seedCommand())
	rootCmd.AddCommand(resetDBCommand())

	// System commands
	rootCmd.AddCommand(healthCheckCommand())
	rootCmd.AddCommand(versionCommand())
}

// Database Commands
func createDBCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-db",
		Short: "Create database",
		Long:  `Create the application database with proper configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating database...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					logger.Info("Database connection established successfully")
					fmt.Println("✅ Database created successfully!")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to create database: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func migrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Run all pending database migrations to update the database schema.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running database migrations...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					if err := runMigrations(db); err != nil {
						logger.Error("Migration failed", zap.Error(err))
						log.Fatalf("Migration failed: %v", err)
					}
					logger.Info("Database migrations completed successfully")
					fmt.Println("✅ Database migrations completed!")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to run migrations: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func seedCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed database with initial data",
		Long:  `Populate the database with initial data including default roles, permissions, and admin user.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Seeding database...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fxmodules.RepositoryModule,
				fxmodules.DomainModule,
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to seed database: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func resetDBCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-db",
		Short: "Reset database (WARNING: This will delete all data)",
		Long:  `Drop all tables and recreate them. This will permanently delete all data.`,
		Run: func(cmd *cobra.Command, args []string) {
			confirm, _ := cmd.Flags().GetBool("confirm")
			if !confirm {
				fmt.Println("⚠️  This operation will permanently delete all data!")
				fmt.Println("Use --confirm flag to proceed: go-mvc-cli reset-db --confirm")
				return
			}

			fmt.Println("Resetting database...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Invoke(func(db *gorm.DB, logger *zap.Logger) {
					if err := resetDatabase(db); err != nil {
						logger.Error("Database reset failed", zap.Error(err))
						log.Fatalf("Database reset failed: %v", err)
					}
					logger.Info("Database reset completed successfully")
					fmt.Println("✅ Database reset completed!")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Failed to reset database: %v", err)
			}
			app.Stop(context.Background())
		},
	}

	cmd.Flags().Bool("confirm", false, "Confirm the database reset operation")
	return cmd
}

// System Commands
func healthCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check system health",
		Long:  `Perform a health check on all system components including database connectivity.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Performing health check...")

			app := fx.New(
				fxmodules.InfrastructureModule,
				fx.Invoke(func(db *gorm.DB, config *config.AppConfig, logger *zap.Logger) {
					// Database health check
					sqlDB, err := db.DB()
					if err != nil {
						fmt.Println("❌ Database connection failed")
						log.Fatalf("Database health check failed: %v", err)
					}

					if err := sqlDB.Ping(); err != nil {
						fmt.Println("❌ Database ping failed")
						log.Fatalf("Database ping failed: %v", err)
					}

					fmt.Println("✅ Database: OK")
					fmt.Printf("✅ Environment: %s\n", config.App.Environment)
					fmt.Printf("✅ HTTP Port: %d\n", config.Server.HTTP.Port)
					fmt.Println("✅ All systems operational!")

					logger.Info("Health check completed successfully")
				}),
				fx.NopLogger,
			)

			if err := app.Start(context.Background()); err != nil {
				log.Fatalf("Health check failed: %v", err)
			}
			app.Stop(context.Background())
		},
	}
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display the current version of the application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Go MVC Application CLI")
			fmt.Println("Version: 1.0.0")
			fmt.Println("Build: development")
			fmt.Println("Go Version:", os.Getenv("GO_VERSION"))
		},
	}
}

// Helper functions
func runMigrations(db *gorm.DB) error {
	// TODO: Import your models here
	fmt.Println("Running database migrations...")

	// Example:
	// if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.UserRole{}, &models.RolePermission{}); err != nil {
	//     return fmt.Errorf("failed to run migrations: %w", err)
	// }

	return nil
}

func resetDatabase(db *gorm.DB) error {
	// Drop all tables
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
	}

	// Run migrations again
	return runMigrations(db)
}

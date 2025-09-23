#!/bin/bash

# Migration script for PostgreSQL database

set -e

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-go_mvc_dev}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}

# Migration directory
MIGRATION_DIR="$(dirname "$0")/../internal/adapters/persistence/postgres/migrations"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check database connection
check_connection() {
    print_status "Checking database connection..."
    if ! docker exec dev-postgres pg_isready -h localhost -p 5432 -U postgres > /dev/null 2>&1; then
        print_error "Cannot connect to PostgreSQL. Is the container running?"
        exit 1
    fi
    print_status "Database connection successful"
}

# Function to create database if not exists
create_database() {
    print_status "Creating database if not exists..."
    docker exec dev-postgres psql -U postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || print_warning "Database $DB_NAME already exists"
}

# Function to run migration
run_migration() {
    local migration_file=$1
    print_status "Running migration: $(basename "$migration_file")"
    
    if docker exec -i dev-postgres psql -U postgres -d "$DB_NAME" < "$migration_file"; then
        print_status "Migration $(basename "$migration_file") completed successfully"
    else
        print_error "Migration $(basename "$migration_file") failed"
        exit 1
    fi
}

# Function to run all up migrations
migrate_up() {
    print_status "Running database migrations..."
    
    check_connection
    create_database
    
    # Find all .up.sql files and sort them
    for migration_file in $(find "$MIGRATION_DIR" -name "*.up.sql" | sort); do
        run_migration "$migration_file"
    done
    
    print_status "All migrations completed successfully!"
}

# Function to run all down migrations
migrate_down() {
    print_status "Rolling back database migrations..."
    
    check_connection
    
    # Find all .down.sql files and sort them in reverse
    for migration_file in $(find "$MIGRATION_DIR" -name "*.down.sql" | sort -r); do
        run_migration "$migration_file"
    done
    
    print_status "All rollbacks completed successfully!"
}

# Function to show migration status
migrate_status() {
    print_status "Checking migration status..."
    
    check_connection
    
    print_status "Available migration files:"
    find "$MIGRATION_DIR" -name "*.up.sql" | sort | while read -r file; do
        echo "  - $(basename "$file")"
    done
    
    print_status "Current database tables:"
    docker exec dev-postgres psql -U postgres -d "$DB_NAME" -c "\dt" 2>/dev/null || print_warning "No tables found or database doesn't exist"
}

# Main script logic
case "${1:-up}" in
    "up")
        migrate_up
        ;;
    "down")
        migrate_down
        ;;
    "status")
        migrate_status
        ;;
    *)
        echo "Usage: $0 {up|down|status}"
        echo "  up     - Run all pending migrations"
        echo "  down   - Rollback all migrations"
        echo "  status - Show migration status"
        exit 1
        ;;
esac
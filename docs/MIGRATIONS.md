# Database Migrations Guide

This document provides comprehensive guidance for managing database schema changes using golang-migrate/migrate.

## 📋 Table of Contents
- [Overview](#overview)
- [Quick Start](#quick-start)
- [Migration Commands](#migration-commands)
- [Writing Migrations](#writing-migrations)
- [Best Practices](#best-practices)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)
- [Team Workflow](#team-workflow)

## 🔍 Overview

### Migration System
- **Tool**: [golang-migrate/migrate](https://github.com/golang-migrate/migrate) v4.19.0+
- **Database**: PostgreSQL
- **Location**: `internal/adapters/persistence/postgres/migrations/`
- **Format**: Timestamp-based naming with separate up/down SQL files
- **Integration**: Full Makefile integration for easy use

### File Structure
```
internal/adapters/persistence/postgres/migrations/
├── 20250923181241_create_users_table.up.sql
├── 20250923181241_create_users_table.down.sql
├── 20250923182421_add_user_avatar.up.sql
└── 20250923182421_add_user_avatar.down.sql
```

## 🚀 Quick Start

### 1. Setup Environment
```bash
# Start database
make docker-up-db

# Check migration status
make migrate-status
```

### 2. Apply Migrations
```bash
# Apply all pending migrations
make migrate-up
```

### 3. Create New Migration
```bash
# Create migration for new feature
make migrate-create name=add_product_table

# Edit the generated files
vim internal/adapters/persistence/postgres/migrations/[timestamp]_add_product_table.up.sql
vim internal/adapters/persistence/postgres/migrations/[timestamp]_add_product_table.down.sql

# Apply the migration
make migrate-up
```

## 🛠️ Migration Commands

### Basic Operations
```bash
# Apply all pending migrations
make migrate-up

# Show current migration status
make migrate-status
# Output: Migration Path: internal/adapters/persistence/postgres/migrations
#         Current Version: 20250923181241
#         Available Migrations: 2 migration files found

# Show current version
make migrate-version
# Output: 20250923181241

# Rollback one migration
make migrate-down-1

# Rollback all migrations (DANGEROUS!)
make migrate-down
```

### Creating Migrations
```bash
# Create new migration with descriptive name
make migrate-create name=create_products_table
make migrate-create name=add_user_phone_column
make migrate-create name=create_orders_index
make migrate-create name=update_user_email_constraint
```

### Advanced Operations
```bash
# Force migration to specific version (for fixing dirty state)
make migrate-force

# Drop all migrations and data (DANGER!)
make migrate-drop
```

## ✍️ Writing Migrations

### Up Migration Template
```sql
-- internal/adapters/persistence/postgres/migrations/[timestamp]_description.up.sql

-- Create table with proper constraints
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT check_price_positive CHECK (price > 0),
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);

-- Add trigger for automatic updated_at
CREATE TRIGGER update_products_updated_at 
    BEFORE UPDATE ON products 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert initial data if needed
INSERT INTO categories (id, name) VALUES 
    (gen_random_uuid(), 'Electronics'),
    (gen_random_uuid(), 'Books')
ON CONFLICT (name) DO NOTHING;
```

### Down Migration Template
```sql
-- internal/adapters/persistence/postgres/migrations/[timestamp]_description.down.sql

-- Remove trigger
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

-- Remove indexes
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_category_id;

-- Remove constraints (if they were added in this migration)
-- ALTER TABLE products DROP CONSTRAINT IF EXISTS check_price_positive;

-- Drop table
DROP TABLE IF EXISTS products;

-- Remove initial data (only if it was added in this migration)
-- DELETE FROM categories WHERE name IN ('Electronics', 'Books');
```

## 📋 Best Practices

### 1. Naming Conventions
```bash
# ✅ Good examples
make migrate-create name=create_users_table
make migrate-create name=add_user_avatar_column
make migrate-create name=create_products_index
make migrate-create name=update_user_email_constraint
make migrate-create name=remove_deprecated_status_column

# ❌ Bad examples
make migrate-create name=fix_bug
make migrate-create name=update_table
make migrate-create name=temp_change
make migrate-create name=v2
```

### 2. Safe Migration Patterns

#### Always Use Conditional Statements
```sql
-- ✅ Good - Safe to run multiple times
CREATE TABLE IF NOT EXISTS users (...);
DROP TABLE IF EXISTS old_table;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ❌ Bad - Will fail if run multiple times
CREATE TABLE users (...);
DROP TABLE old_table;
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
CREATE INDEX idx_users_email ON users(email);
```

#### Handle Existing Data Carefully
```sql
-- ✅ Good - Safe data migration
-- Add column with default value
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';

-- Update existing NULL values
UPDATE users SET status = 'active' WHERE status IS NULL;

-- Add NOT NULL constraint after data is clean
-- Note: This might be a separate migration for safety
-- ALTER TABLE users ALTER COLUMN status SET NOT NULL;
```

#### Foreign Key Management
```sql
-- ✅ Good - Handle dependencies properly
-- Up migration
ALTER TABLE posts ADD COLUMN IF NOT EXISTS author_id UUID;
CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
ALTER TABLE posts ADD CONSTRAINT fk_posts_author 
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL;

-- Down migration (reverse order)
ALTER TABLE posts DROP CONSTRAINT IF EXISTS fk_posts_author;
DROP INDEX IF EXISTS idx_posts_author_id;
ALTER TABLE posts DROP COLUMN IF EXISTS author_id;
```

### 3. Transaction Considerations
```sql
-- Most DDL statements in PostgreSQL are transactional
-- But some operations like CREATE INDEX CONCURRENTLY are not
-- Keep migrations simple and atomic when possible

-- For large tables, consider these patterns:
-- 1. Add column without NOT NULL constraint
-- 2. Populate data in batches (separate migration)
-- 3. Add NOT NULL constraint (separate migration)
```

### 4. Performance Considerations
```sql
-- For large tables, create indexes concurrently
-- Note: This cannot be done in a transaction
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_large_table_column 
    ON large_table(column_name);

-- For adding NOT NULL to large tables
-- 1. Add column with default
-- 2. Update in batches
-- 3. Add constraint with VALIDATE (separate migration)
```

## 🔄 Common Patterns

### 1. Creating Related Tables
```sql
-- Migration: create_user_profile_system.up.sql

-- Create parent table first
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    bio TEXT,
    avatar_url VARCHAR(500),
    website VARCHAR(255),
    location VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);

-- Trigger
CREATE TRIGGER update_user_profiles_updated_at 
    BEFORE UPDATE ON user_profiles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

### 2. Adding Enums
```sql
-- Create enum type
CREATE TYPE IF NOT EXISTS user_status AS ENUM ('active', 'inactive', 'suspended', 'deleted');

-- Use in table
ALTER TABLE users ADD COLUMN IF NOT EXISTS status user_status DEFAULT 'active';

-- Down migration
ALTER TABLE users DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS user_status;
```

### 3. Many-to-Many Relationships
```sql
-- Junction table with additional fields
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    granted_by UUID,
    granted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_granted_by FOREIGN KEY (granted_by) REFERENCES users(id),
    
    -- Prevent duplicate assignments
    UNIQUE(user_id, role_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_granted_by ON user_roles(granted_by);
CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON user_roles(expires_at) WHERE expires_at IS NOT NULL;
```

### 4. Data Migration with Transformation
```sql
-- Example: Split full_name into first_name and last_name
-- Up migration
ALTER TABLE users ADD COLUMN IF NOT EXISTS first_name VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_name VARCHAR(100);

-- Migrate existing data
UPDATE users 
SET 
    first_name = COALESCE(split_part(full_name, ' ', 1), ''),
    last_name = COALESCE(split_part(full_name, ' ', 2), '')
WHERE full_name IS NOT NULL;

-- Down migration (careful about data loss!)
UPDATE users 
SET full_name = CONCAT(first_name, ' ', last_name)
WHERE first_name IS NOT NULL OR last_name IS NOT NULL;

ALTER TABLE users DROP COLUMN IF EXISTS first_name;
ALTER TABLE users DROP COLUMN IF EXISTS last_name;
```

## 🔧 Troubleshooting

### 1. Dirty Migration State
```bash
# Symptoms: "Dirty database version X. Fix and force version."
# Check current state
make migrate-version
# Output: 20250923181241 (dirty)

# Find the issue in logs
make docker-logs

# Fix the migration file
vim internal/adapters/persistence/postgres/migrations/[timestamp]_name.up.sql

# Force clean the state
make migrate-force

# Try again
make migrate-up
```

### 2. Migration Fails
```bash
# Check database connection
make docker-ps

# Review migration syntax
cat internal/adapters/persistence/postgres/migrations/[timestamp]_name.up.sql

# Test SQL manually
docker exec -it dev-postgres psql -U postgres -d go_mvc_dev
# > \i /path/to/migration.sql

# Fix and retry
make migrate-up
```

### 3. Foreign Key Constraint Violations
```sql
-- Check existing data before adding constraints
SELECT column_name, COUNT(*) as null_count
FROM table_name 
WHERE foreign_key_column IS NULL
GROUP BY column_name;

-- Clean data first, then add constraint in separate migration
UPDATE table_name SET foreign_key_column = some_default_value 
WHERE foreign_key_column IS NULL;
```

### 4. Index Creation on Large Tables
```bash
# For large tables, use CONCURRENTLY (not in transaction)
# Create separate migration file just for index
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_large_table_column 
    ON large_table(column_name);

# Monitor progress
SELECT 
    schemaname, 
    tablename, 
    indexname, 
    indexdef 
FROM pg_indexes 
WHERE indexname = 'idx_large_table_column';
```

## 👥 Team Workflow

### 1. Starting Development
```bash
# Daily routine for developers
git pull origin main
make docker-up-db
make migrate-status
make migrate-up  # Apply any new migrations
make test        # Verify everything works
```

### 2. Creating Features
```bash
# Feature development workflow
git checkout -b feature/user-profiles

# Create migration for feature
make migrate-create name=create_user_profiles_table

# Edit migration files
# ... implement feature ...

# Test migration cycle
make migrate-up
make migrate-down-1  # Test rollback
make migrate-up      # Apply again

# Test with application
make test
make run

# Commit when ready
git add .
git commit -m "feat: add user profiles with migration"
```

### 3. Code Review Checklist
- [ ] Migration has both up and down files
- [ ] Uses IF EXISTS/IF NOT EXISTS appropriately
- [ ] Includes proper indexes for performance
- [ ] Down migration properly reverses changes
- [ ] Migration tested locally (up and down)
- [ ] No data loss in down migration (or documented)
- [ ] Follows naming conventions
- [ ] Large table changes use safe patterns

### 4. Deployment Workflow
```bash
# Production deployment
# 1. Deploy migration
make migrate-up DATABASE_URL="$PRODUCTION_URL"

# 2. Verify migration status
make migrate-status DATABASE_URL="$PRODUCTION_URL"

# 3. Deploy application
# ... deploy process ...

# 4. Verify application works
curl https://api.example.com/health
```

### 5. Rollback Procedure
```bash
# Emergency rollback
# 1. Rollback application first
# 2. Then rollback migration if needed
make migrate-down-1 DATABASE_URL="$PRODUCTION_URL"

# 3. Verify database state
make migrate-status DATABASE_URL="$PRODUCTION_URL"
```

## 📚 Additional Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migration Best Practices](https://docs.github.com/en/migrations)
- [Project Development Guide](./DEVELOPMENT.md#database-migrations)
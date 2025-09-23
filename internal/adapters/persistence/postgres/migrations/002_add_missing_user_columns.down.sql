-- Remove index
DROP INDEX IF EXISTS idx_users_is_active;

-- Remove added columns
ALTER TABLE users 
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS name;
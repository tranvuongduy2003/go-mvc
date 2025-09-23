-- Add missing columns to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS name VARCHAR(100) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Update timestamps with default values
ALTER TABLE users 
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

-- Update email column type and constraint
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(255);

-- Create index on is_active for filtering
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Update existing records to have default values
UPDATE users SET 
    name = COALESCE(name, 'Unknown'),
    password_hash = COALESCE(password_hash, ''),
    is_active = COALESCE(is_active, true)
WHERE name = '' OR password_hash = '' OR is_active IS NULL;
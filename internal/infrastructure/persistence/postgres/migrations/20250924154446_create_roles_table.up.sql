-- Create roles table for RBAC system
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 1,
    
    -- Constraints
    CONSTRAINT roles_name_check CHECK (LENGTH(name) >= 2 AND name ~ '^[A-Z][A-Z0-9_]*$'),
    CONSTRAINT roles_description_check CHECK (LENGTH(description) <= 255)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_is_active ON roles(is_active);
CREATE INDEX IF NOT EXISTS idx_roles_created_at ON roles(created_at);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_roles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_roles_updated_at ON roles;
CREATE TRIGGER trigger_update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_roles_updated_at();

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('ADMIN', 'System administrator with full access'),
    ('USER', 'Regular user with basic permissions'),
    ('MODERATOR', 'Content moderator with management permissions')
ON CONFLICT (name) DO NOTHING;

-- Add comments for documentation
COMMENT ON TABLE roles IS 'Roles table for RBAC - stores role definitions';
COMMENT ON COLUMN roles.id IS 'Primary key - UUID for role identification';
COMMENT ON COLUMN roles.name IS 'Role name - uppercase with underscores, unique';
COMMENT ON COLUMN roles.description IS 'Human-readable role description';
COMMENT ON COLUMN roles.is_active IS 'Whether the role is currently active';
COMMENT ON COLUMN roles.created_at IS 'Timestamp when role was created';
COMMENT ON COLUMN roles.updated_at IS 'Timestamp when role was last updated';
COMMENT ON COLUMN roles.version IS 'Optimistic locking version number';

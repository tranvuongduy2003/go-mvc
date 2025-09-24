-- Create permissions table for RBAC system
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 1,
    
    -- Constraints
    CONSTRAINT permissions_name_check CHECK (LENGTH(name) >= 3 AND name ~ '^[a-z][a-z0-9_]*:[a-z][a-z0-9_]*$'),
    CONSTRAINT permissions_resource_check CHECK (LENGTH(resource) >= 2 AND resource ~ '^[a-z][a-z0-9_]*$'),
    CONSTRAINT permissions_action_check CHECK (LENGTH(action) >= 2 AND action ~ '^[a-z][a-z0-9_]*$'),
    CONSTRAINT permissions_description_check CHECK (LENGTH(description) <= 255),
    CONSTRAINT permissions_resource_action_unique UNIQUE (resource, action)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions(resource);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions(action);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX IF NOT EXISTS idx_permissions_is_active ON permissions(is_active);
CREATE INDEX IF NOT EXISTS idx_permissions_created_at ON permissions(created_at);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_permissions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_permissions_updated_at ON permissions;
CREATE TRIGGER trigger_update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_permissions_updated_at();

-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
    -- User management permissions
    ('users:create', 'users', 'create', 'Create new users'),
    ('users:read', 'users', 'read', 'View user details'),
    ('users:update', 'users', 'update', 'Update user information'),
    ('users:delete', 'users', 'delete', 'Delete users'),
    ('users:list', 'users', 'list', 'List all users'),
    
    -- Role management permissions
    ('roles:create', 'roles', 'create', 'Create new roles'),
    ('roles:read', 'roles', 'read', 'View role details'),
    ('roles:update', 'roles', 'update', 'Update role information'),
    ('roles:delete', 'roles', 'delete', 'Delete roles'),
    ('roles:list', 'roles', 'list', 'List all roles'),
    
    -- Permission management permissions
    ('permissions:create', 'permissions', 'create', 'Create new permissions'),
    ('permissions:read', 'permissions', 'read', 'View permission details'),
    ('permissions:update', 'permissions', 'update', 'Update permission information'),
    ('permissions:delete', 'permissions', 'delete', 'Delete permissions'),
    ('permissions:list', 'permissions', 'list', 'List all permissions'),
    
    -- System permissions
    ('system:manage', 'system', 'manage', 'Full system management access')
ON CONFLICT (name) DO NOTHING;

-- Add comments for documentation
COMMENT ON TABLE permissions IS 'Permissions table for RBAC - stores permission definitions';
COMMENT ON COLUMN permissions.id IS 'Primary key - UUID for permission identification';
COMMENT ON COLUMN permissions.name IS 'Permission name in format resource:action';
COMMENT ON COLUMN permissions.resource IS 'Resource the permission applies to';
COMMENT ON COLUMN permissions.action IS 'Action that can be performed';
COMMENT ON COLUMN permissions.description IS 'Human-readable permission description';
COMMENT ON COLUMN permissions.is_active IS 'Whether the permission is currently active';
COMMENT ON COLUMN permissions.created_at IS 'Timestamp when permission was created';
COMMENT ON COLUMN permissions.updated_at IS 'Timestamp when permission was last updated';
COMMENT ON COLUMN permissions.version IS 'Optimistic locking version number';

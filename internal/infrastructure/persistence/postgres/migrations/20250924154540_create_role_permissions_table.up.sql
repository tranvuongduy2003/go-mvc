-- Create role_permissions junction table for many-to-many relationship between roles and permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    granted_by UUID, -- Who granted this permission to the role (optional, for audit trail)
    granted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 1,
    
    -- Foreign key constraints
    CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL,
    
    -- Unique constraint to prevent duplicate permission assignments
    CONSTRAINT role_permissions_role_permission_unique UNIQUE (role_id, permission_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_granted_by ON role_permissions(granted_by);
CREATE INDEX IF NOT EXISTS idx_role_permissions_granted_at ON role_permissions(granted_at);
CREATE INDEX IF NOT EXISTS idx_role_permissions_is_active ON role_permissions(is_active);
CREATE INDEX IF NOT EXISTS idx_role_permissions_active ON role_permissions(role_id, is_active) WHERE is_active = true;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_role_permissions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_role_permissions_updated_at ON role_permissions;
CREATE TRIGGER trigger_update_role_permissions_updated_at
    BEFORE UPDATE ON role_permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_role_permissions_updated_at();

-- Insert default role-permission assignments
DO $$
DECLARE
    admin_role_id UUID;
    user_role_id UUID;
    moderator_role_id UUID;
    perm_id UUID;
BEGIN
    -- Get role IDs
    SELECT id INTO admin_role_id FROM roles WHERE name = 'ADMIN';
    SELECT id INTO user_role_id FROM roles WHERE name = 'USER';
    SELECT id INTO moderator_role_id FROM roles WHERE name = 'MODERATOR';
    
    -- Admin gets all permissions
    FOR perm_id IN SELECT id FROM permissions WHERE is_active = true LOOP
        INSERT INTO role_permissions (role_id, permission_id, granted_at)
        VALUES (admin_role_id, perm_id, CURRENT_TIMESTAMP)
        ON CONFLICT (role_id, permission_id) DO NOTHING;
    END LOOP;
    
    -- User gets basic read permissions
    INSERT INTO role_permissions (role_id, permission_id, granted_at)
    SELECT user_role_id, id, CURRENT_TIMESTAMP
    FROM permissions 
    WHERE name IN ('users:read') AND is_active = true
    ON CONFLICT (role_id, permission_id) DO NOTHING;
    
    -- Moderator gets user management permissions
    INSERT INTO role_permissions (role_id, permission_id, granted_at)
    SELECT moderator_role_id, id, CURRENT_TIMESTAMP
    FROM permissions 
    WHERE name IN ('users:read', 'users:list', 'users:update') AND is_active = true
    ON CONFLICT (role_id, permission_id) DO NOTHING;
    
END $$;

-- Create view for easy permission checking
CREATE OR REPLACE VIEW user_permissions AS
SELECT DISTINCT
    ur.user_id,
    p.id as permission_id,
    p.name as permission_name,
    p.resource,
    p.action,
    p.description,
    r.name as role_name
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE ur.is_active = true
    AND r.is_active = true
    AND rp.is_active = true
    AND p.is_active = true
    AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP);

-- Add comments for documentation
COMMENT ON TABLE role_permissions IS 'Junction table linking roles to permissions for RBAC';
COMMENT ON COLUMN role_permissions.id IS 'Primary key - UUID for role-permission assignment identification';
COMMENT ON COLUMN role_permissions.role_id IS 'Foreign key to roles table';
COMMENT ON COLUMN role_permissions.permission_id IS 'Foreign key to permissions table';
COMMENT ON COLUMN role_permissions.granted_by IS 'User who granted this permission (for audit trail)';
COMMENT ON COLUMN role_permissions.granted_at IS 'Timestamp when permission was granted';
COMMENT ON COLUMN role_permissions.is_active IS 'Whether the permission assignment is currently active';
COMMENT ON COLUMN role_permissions.created_at IS 'Timestamp when assignment was created';
COMMENT ON COLUMN role_permissions.updated_at IS 'Timestamp when assignment was last updated';
COMMENT ON COLUMN role_permissions.version IS 'Optimistic locking version number';

COMMENT ON VIEW user_permissions IS 'View for easily checking user permissions through role assignments';

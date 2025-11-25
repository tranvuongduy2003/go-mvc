-- Create user_roles junction table for many-to-many relationship between users and roles
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    assigned_by UUID, -- Who assigned this role (optional, for audit trail)
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE, -- Optional expiration for temporary role assignments
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 1,
    
    -- Foreign key constraints
    CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL,
    
    -- Unique constraint to prevent duplicate role assignments
    CONSTRAINT user_roles_user_role_unique UNIQUE (user_id, role_id),
    
    -- Check constraint for expiration logic
    CONSTRAINT user_roles_expiration_check CHECK (expires_at IS NULL OR expires_at > assigned_at)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_assigned_by ON user_roles(assigned_by);
CREATE INDEX IF NOT EXISTS idx_user_roles_assigned_at ON user_roles(assigned_at);
CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON user_roles(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_user_roles_is_active ON user_roles(is_active);
CREATE INDEX IF NOT EXISTS idx_user_roles_active_unexpired ON user_roles(user_id, is_active, expires_at) WHERE is_active = true;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_user_roles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_user_roles_updated_at ON user_roles;
CREATE TRIGGER trigger_update_user_roles_updated_at
    BEFORE UPDATE ON user_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_user_roles_updated_at();

-- Create function to check if role assignment is still valid (not expired)
CREATE OR REPLACE FUNCTION is_user_role_valid(user_role_record user_roles)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN user_role_record.is_active = true 
           AND (user_role_record.expires_at IS NULL OR user_role_record.expires_at > CURRENT_TIMESTAMP);
END;
$$ language 'plpgsql';

-- Add comments for documentation
COMMENT ON TABLE user_roles IS 'Junction table linking users to roles for RBAC';
COMMENT ON COLUMN user_roles.id IS 'Primary key - UUID for user-role assignment identification';
COMMENT ON COLUMN user_roles.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN user_roles.role_id IS 'Foreign key to roles table';
COMMENT ON COLUMN user_roles.assigned_by IS 'User who assigned this role (for audit trail)';
COMMENT ON COLUMN user_roles.assigned_at IS 'Timestamp when role was assigned';
COMMENT ON COLUMN user_roles.expires_at IS 'Optional expiration timestamp for temporary assignments';
COMMENT ON COLUMN user_roles.is_active IS 'Whether the role assignment is currently active';
COMMENT ON COLUMN user_roles.created_at IS 'Timestamp when assignment was created';
COMMENT ON COLUMN user_roles.updated_at IS 'Timestamp when assignment was last updated';
COMMENT ON COLUMN user_roles.version IS 'Optimistic locking version number';

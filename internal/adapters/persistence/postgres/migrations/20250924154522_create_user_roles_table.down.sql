-- Drop user_roles table and related objects
DROP TRIGGER IF EXISTS trigger_update_user_roles_updated_at ON user_roles;
DROP FUNCTION IF EXISTS update_user_roles_updated_at();
DROP FUNCTION IF EXISTS is_user_role_valid(user_roles);
DROP INDEX IF EXISTS idx_user_roles_active_unexpired;
DROP INDEX IF EXISTS idx_user_roles_is_active;
DROP INDEX IF EXISTS idx_user_roles_expires_at;
DROP INDEX IF EXISTS idx_user_roles_assigned_at;
DROP INDEX IF EXISTS idx_user_roles_assigned_by;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP TABLE IF EXISTS user_roles CASCADE;

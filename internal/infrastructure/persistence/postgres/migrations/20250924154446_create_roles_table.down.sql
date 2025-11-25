-- Drop roles table and related objects
DROP TRIGGER IF EXISTS trigger_update_roles_updated_at ON roles;
DROP FUNCTION IF EXISTS update_roles_updated_at();
DROP INDEX IF EXISTS idx_roles_created_at;
DROP INDEX IF EXISTS idx_roles_is_active;
DROP INDEX IF EXISTS idx_roles_name;
DROP TABLE IF EXISTS roles CASCADE;

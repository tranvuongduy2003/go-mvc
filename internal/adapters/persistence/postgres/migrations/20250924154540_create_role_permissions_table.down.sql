-- Drop role_permissions table and related objects
DROP VIEW IF EXISTS user_permissions;
DROP TRIGGER IF EXISTS trigger_update_role_permissions_updated_at ON role_permissions;
DROP FUNCTION IF EXISTS update_role_permissions_updated_at();
DROP INDEX IF EXISTS idx_role_permissions_active;
DROP INDEX IF EXISTS idx_role_permissions_is_active;
DROP INDEX IF EXISTS idx_role_permissions_granted_at;
DROP INDEX IF EXISTS idx_role_permissions_granted_by;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP TABLE IF EXISTS role_permissions CASCADE;

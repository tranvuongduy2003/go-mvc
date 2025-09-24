-- Drop permissions table and related objects
DROP TRIGGER IF EXISTS trigger_update_permissions_updated_at ON permissions;
DROP FUNCTION IF EXISTS update_permissions_updated_at();
DROP INDEX IF EXISTS idx_permissions_created_at;
DROP INDEX IF EXISTS idx_permissions_is_active;
DROP INDEX IF EXISTS idx_permissions_resource_action;
DROP INDEX IF EXISTS idx_permissions_action;
DROP INDEX IF EXISTS idx_permissions_resource;
DROP INDEX IF EXISTS idx_permissions_name;
DROP TABLE IF EXISTS permissions CASCADE;

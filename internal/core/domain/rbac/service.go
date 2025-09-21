package rbac

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RBACService defines the RBAC service interface
type RBACService interface {
	// Role management
	CreateRole(ctx context.Context, name, description string, createdBy uuid.UUID) (*Role, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role, updatedBy uuid.UUID) error
	DeleteRole(ctx context.Context, roleID uuid.UUID, deletedBy uuid.UUID) error
	ListRoles(ctx context.Context, offset, limit int) ([]*Role, error)

	// Permission management
	CreatePermission(ctx context.Context, name, resource, action, description string, createdBy uuid.UUID) (*Permission, error)
	GetPermission(ctx context.Context, permissionID uuid.UUID) (*Permission, error)
	GetPermissionByResourceAndAction(ctx context.Context, resource, action string) (*Permission, error)
	UpdatePermission(ctx context.Context, permission *Permission, updatedBy uuid.UUID) error
	DeletePermission(ctx context.Context, permissionID uuid.UUID, deletedBy uuid.UUID) error
	ListPermissions(ctx context.Context, offset, limit int) ([]*Permission, error)

	// Role-Permission management
	AssignPermissionToRole(ctx context.Context, roleID, permissionID, assignedBy uuid.UUID) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID, removedBy uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)

	// User-Role management
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy uuid.UUID, expiresAt *time.Time) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID, removedBy uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*Role, error)
	GetActiveUserRoles(ctx context.Context, userID uuid.UUID) ([]*Role, error)

	// Authorization checks
	HasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)
	HasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error)
	HasAnyRole(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error)
	HasAllRoles(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error)

	// Advanced authorization
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]*Permission, error)
	CanAccessResource(ctx context.Context, userID uuid.UUID, resource string, actions []string) (map[string]bool, error)

	// Bulk operations
	AssignMultipleRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error
	AssignMultiplePermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, assignedBy uuid.UUID) error

	// System operations
	InitializeDefaultRolesAndPermissions(ctx context.Context, systemUserID uuid.UUID) error
	CleanupExpiredRoles(ctx context.Context) error
}

// RBACServiceImpl implements RBACService
type RBACServiceImpl struct {
	roleRepo           RoleRepository
	permissionRepo     PermissionRepository
	userRoleRepo       UserRoleRepository
	rolePermissionRepo RolePermissionRepository
}

// NewRBACService creates a new RBAC service
func NewRBACService(
	roleRepo RoleRepository,
	permissionRepo PermissionRepository,
	userRoleRepo UserRoleRepository,
	rolePermissionRepo RolePermissionRepository,
) RBACService {
	return &RBACServiceImpl{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		userRoleRepo:       userRoleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

// CreateRole creates a new role
func (s *RBACServiceImpl) CreateRole(ctx context.Context, name, description string, createdBy uuid.UUID) (*Role, error) {
	// Check if role already exists
	existingRole, err := s.roleRepo.GetByName(ctx, name)
	if err == nil && existingRole != nil {
		return nil, fmt.Errorf("role with name '%s' already exists", name)
	}

	role := NewRole(name, description, createdBy)
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// GetRole retrieves a role by ID
func (s *RBACServiceImpl) GetRole(ctx context.Context, roleID uuid.UUID) (*Role, error) {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

// GetRoleByName retrieves a role by name
func (s *RBACServiceImpl) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	role, err := s.roleRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	return role, nil
}

// CreatePermission creates a new permission
func (s *RBACServiceImpl) CreatePermission(ctx context.Context, name, resource, action, description string, createdBy uuid.UUID) (*Permission, error) {
	// Check if permission already exists
	existingPermission, err := s.permissionRepo.GetByResourceAndAction(ctx, resource, action)
	if err == nil && existingPermission != nil {
		return nil, fmt.Errorf("permission for resource '%s' and action '%s' already exists", resource, action)
	}

	permission := NewPermission(name, resource, action, description, createdBy)
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

// HasPermission checks if a user has a specific permission
func (s *RBACServiceImpl) HasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	// Get user's active roles
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Check each role for the permission
	for _, userRole := range userRoles {
		if userRole.IsValid() {
			rolePermissions, err := s.rolePermissionRepo.GetActiveRolePermissions(ctx, userRole.RoleID)
			if err != nil {
				continue
			}

			for _, rolePermission := range rolePermissions {
				permission, err := s.permissionRepo.GetByID(ctx, rolePermission.PermissionID)
				if err != nil {
					continue
				}

				if permission.Resource == resource && permission.Action == action && permission.IsActive {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// HasRole checks if a user has a specific role
func (s *RBACServiceImpl) HasRole(ctx context.Context, userID, roleID uuid.UUID) (bool, error) {
	return s.userRoleRepo.HasRole(ctx, userID, roleID)
}

// GetUserPermissions gets all permissions for a user
func (s *RBACServiceImpl) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]*Permission, error) {
	// Get user's active roles
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	permissionMap := make(map[uuid.UUID]*Permission)

	// Collect all permissions from all roles
	for _, userRole := range userRoles {
		if userRole.IsValid() {
			rolePermissions, err := s.rolePermissionRepo.GetActiveRolePermissions(ctx, userRole.RoleID)
			if err != nil {
				continue
			}

			for _, rolePermission := range rolePermissions {
				permission, err := s.permissionRepo.GetByID(ctx, rolePermission.PermissionID)
				if err != nil {
					continue
				}

				if permission.IsActive {
					permissionMap[permission.ID] = permission
				}
			}
		}
	}

	// Convert map to slice
	permissions := make([]*Permission, 0, len(permissionMap))
	for _, permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// AssignRoleToUser assigns a role to a user
func (s *RBACServiceImpl) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy uuid.UUID, expiresAt *time.Time) error {
	// Check if role exists
	_, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Check if assignment already exists
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check existing role assignment: %w", err)
	}

	if hasRole {
		return fmt.Errorf("user already has this role")
	}

	userRole := NewUserRole(userID, roleID, assignedBy, expiresAt)
	if err := s.userRoleRepo.AssignRoleToUser(ctx, userRole); err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// InitializeDefaultRolesAndPermissions creates default system roles and permissions
func (s *RBACServiceImpl) InitializeDefaultRolesAndPermissions(ctx context.Context, systemUserID uuid.UUID) error {
	// Create default permissions
	defaultPermissions := []struct {
		name, resource, action, description string
	}{
		{"Create User", ResourceUser, ActionCreate, "Create new users"},
		{"Read User", ResourceUser, ActionRead, "View user details"},
		{"Update User", ResourceUser, ActionUpdate, "Update user information"},
		{"Delete User", ResourceUser, ActionDelete, "Delete users"},
		{"List Users", ResourceUser, ActionList, "List all users"},

		{"Create Role", ResourceRole, ActionCreate, "Create new roles"},
		{"Read Role", ResourceRole, ActionRead, "View role details"},
		{"Update Role", ResourceRole, ActionUpdate, "Update role information"},
		{"Delete Role", ResourceRole, ActionDelete, "Delete roles"},
		{"List Roles", ResourceRole, ActionList, "List all roles"},

		{"System Admin", ResourceSystem, ActionAdmin, "Full system administration"},
		{"Export Reports", ResourceReport, ActionExport, "Export system reports"},
		{"View Audit", ResourceAudit, ActionRead, "View audit logs"},
	}

	createdPermissions := make(map[string]*Permission)
	for _, perm := range defaultPermissions {
		permission, err := s.CreatePermission(ctx, perm.name, perm.resource, perm.action, perm.description, systemUserID)
		if err != nil {
			// Permission might already exist, try to get it
			existing, getErr := s.permissionRepo.GetByResourceAndAction(ctx, perm.resource, perm.action)
			if getErr != nil {
				return fmt.Errorf("failed to create/get permission %s: %w", perm.name, err)
			}
			permission = existing
		}
		createdPermissions[perm.resource+":"+perm.action] = permission
	}

	// Create default roles
	defaultRoles := []struct {
		name, description string
		permissions       []string
	}{
		{
			name:        RoleAdmin,
			description: "System administrator with full access",
			permissions: []string{
				ResourceUser + ":" + ActionCreate,
				ResourceUser + ":" + ActionRead,
				ResourceUser + ":" + ActionUpdate,
				ResourceUser + ":" + ActionDelete,
				ResourceUser + ":" + ActionList,
				ResourceRole + ":" + ActionCreate,
				ResourceRole + ":" + ActionRead,
				ResourceRole + ":" + ActionUpdate,
				ResourceRole + ":" + ActionDelete,
				ResourceRole + ":" + ActionList,
				ResourceSystem + ":" + ActionAdmin,
				ResourceReport + ":" + ActionExport,
				ResourceAudit + ":" + ActionRead,
			},
		},
		{
			name:        RoleUser,
			description: "Regular user with limited access",
			permissions: []string{
				ResourceUser + ":" + ActionRead,
			},
		},
		{
			name:        RoleModerator,
			description: "Moderator with user management access",
			permissions: []string{
				ResourceUser + ":" + ActionRead,
				ResourceUser + ":" + ActionUpdate,
				ResourceUser + ":" + ActionList,
			},
		},
	}

	for _, role := range defaultRoles {
		createdRole, err := s.CreateRole(ctx, role.name, role.description, systemUserID)
		if err != nil {
			// Role might already exist, try to get it
			existing, getErr := s.roleRepo.GetByName(ctx, role.name)
			if getErr != nil {
				return fmt.Errorf("failed to create/get role %s: %w", role.name, err)
			}
			createdRole = existing
		}

		// Assign permissions to role
		for _, permKey := range role.permissions {
			if permission, exists := createdPermissions[permKey]; exists {
				err := s.AssignPermissionToRole(ctx, createdRole.ID, permission.ID, systemUserID)
				if err != nil {
					// Permission might already be assigned, continue
					continue
				}
			}
		}
	}

	return nil
}

// AssignPermissionToRole assigns a permission to a role
func (s *RBACServiceImpl) AssignPermissionToRole(ctx context.Context, roleID, permissionID, assignedBy uuid.UUID) error {
	// Check if role exists
	_, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Check if permission exists
	_, err = s.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	// Check if assignment already exists
	hasPermission, err := s.rolePermissionRepo.HasPermission(ctx, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to check existing permission assignment: %w", err)
	}

	if hasPermission {
		return nil // Already assigned
	}

	rolePermission := NewRolePermission(roleID, permissionID, assignedBy)
	if err := s.rolePermissionRepo.GrantPermissionToRole(ctx, rolePermission); err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}

// UpdateRole updates an existing role
func (s *RBACServiceImpl) UpdateRole(ctx context.Context, role *Role, updatedBy uuid.UUID) error {
	role.AuditLog.Update(updatedBy)
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

// DeleteRole deletes a role
func (s *RBACServiceImpl) DeleteRole(ctx context.Context, roleID uuid.UUID, deletedBy uuid.UUID) error {
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// ListRoles lists roles with pagination
func (s *RBACServiceImpl) ListRoles(ctx context.Context, offset, limit int) ([]*Role, error) {
	roles, err := s.roleRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}

// GetPermission retrieves a permission by ID
func (s *RBACServiceImpl) GetPermission(ctx context.Context, permissionID uuid.UUID) (*Permission, error) {
	permission, err := s.permissionRepo.GetByID(ctx, permissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	return permission, nil
}

// GetPermissionByResourceAndAction retrieves a permission by resource and action
func (s *RBACServiceImpl) GetPermissionByResourceAndAction(ctx context.Context, resource, action string) (*Permission, error) {
	permission, err := s.permissionRepo.GetByResourceAndAction(ctx, resource, action)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	return permission, nil
}

// UpdatePermission updates an existing permission
func (s *RBACServiceImpl) UpdatePermission(ctx context.Context, permission *Permission, updatedBy uuid.UUID) error {
	permission.AuditLog.Update(updatedBy)
	if err := s.permissionRepo.Update(ctx, permission); err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}
	return nil
}

// DeletePermission deletes a permission
func (s *RBACServiceImpl) DeletePermission(ctx context.Context, permissionID uuid.UUID, deletedBy uuid.UUID) error {
	if err := s.permissionRepo.Delete(ctx, permissionID); err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// ListPermissions lists permissions with pagination
func (s *RBACServiceImpl) ListPermissions(ctx context.Context, offset, limit int) ([]*Permission, error) {
	permissions, err := s.permissionRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	return permissions, nil
}

// RemovePermissionFromRole removes a permission from a role
func (s *RBACServiceImpl) RemovePermissionFromRole(ctx context.Context, roleID, permissionID, removedBy uuid.UUID) error {
	if err := s.rolePermissionRepo.RevokePermissionFromRole(ctx, roleID, permissionID, removedBy); err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}
	return nil
}

// GetRolePermissions gets all permissions for a role
func (s *RBACServiceImpl) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*Permission, error) {
	permissions, err := s.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	return permissions, nil
}

// RemoveRoleFromUser removes a role from a user
func (s *RBACServiceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID, removedBy uuid.UUID) error {
	if err := s.userRoleRepo.RemoveRoleFromUser(ctx, userID, roleID, removedBy); err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return nil
}

// GetUserRoles gets all roles for a user
func (s *RBACServiceImpl) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*Role, error) {
	userRoles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleIDs := make([]uuid.UUID, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	roles, err := s.roleRepo.GetRolesWithPermissions(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles with permissions: %w", err)
	}

	return roles, nil
}

// GetActiveUserRoles gets all active roles for a user
func (s *RBACServiceImpl) GetActiveUserRoles(ctx context.Context, userID uuid.UUID) ([]*Role, error) {
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active user roles: %w", err)
	}

	roleIDs := make([]uuid.UUID, 0)
	for _, ur := range userRoles {
		if ur.IsValid() {
			roleIDs = append(roleIDs, ur.RoleID)
		}
	}

	if len(roleIDs) == 0 {
		return []*Role{}, nil
	}

	roles, err := s.roleRepo.GetRolesWithPermissions(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles with permissions: %w", err)
	}

	return roles, nil
}

// HasAnyRole checks if user has any of the specified roles
func (s *RBACServiceImpl) HasAnyRole(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error) {
	for _, roleName := range roleNames {
		role, err := s.roleRepo.GetByName(ctx, roleName)
		if err != nil {
			continue
		}

		hasRole, err := s.userRoleRepo.HasRole(ctx, userID, role.ID)
		if err != nil {
			continue
		}

		if hasRole {
			return true, nil
		}
	}

	return false, nil
}

// HasAllRoles checks if user has all of the specified roles
func (s *RBACServiceImpl) HasAllRoles(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error) {
	for _, roleName := range roleNames {
		role, err := s.roleRepo.GetByName(ctx, roleName)
		if err != nil {
			return false, err
		}

		hasRole, err := s.userRoleRepo.HasRole(ctx, userID, role.ID)
		if err != nil {
			return false, err
		}

		if !hasRole {
			return false, nil
		}
	}

	return true, nil
}

// CanAccessResource checks if user can perform actions on a resource
func (s *RBACServiceImpl) CanAccessResource(ctx context.Context, userID uuid.UUID, resource string, actions []string) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, action := range actions {
		hasPermission, err := s.HasPermission(ctx, userID, resource, action)
		if err != nil {
			result[action] = false
		} else {
			result[action] = hasPermission
		}
	}

	return result, nil
}

// AssignMultipleRolesToUser assigns multiple roles to a user
func (s *RBACServiceImpl) AssignMultipleRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error {
	return s.userRoleRepo.AssignMultipleRolesToUser(ctx, userID, roleIDs, assignedBy)
}

// AssignMultiplePermissionsToRole assigns multiple permissions to a role
func (s *RBACServiceImpl) AssignMultiplePermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, assignedBy uuid.UUID) error {
	return s.rolePermissionRepo.GrantMultiplePermissionsToRole(ctx, roleID, permissionIDs, assignedBy)
}

// CleanupExpiredRoles removes expired role assignments
func (s *RBACServiceImpl) CleanupExpiredRoles(ctx context.Context) error {
	expiredRoles, err := s.userRoleRepo.GetExpiredRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired roles: %w", err)
	}

	for _, userRole := range expiredRoles {
		// Mark as inactive instead of deleting for audit purposes
		userRole.IsActive = false
		// Note: This would require an Update method in UserRoleRepository
	}

	return nil
}

// Additional methods would be implemented here...

package services

import (
	"context"
	"fmt"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

type authorizationService struct {
	userRepo           user.UserRepository
	roleRepo           auth.RoleRepository
	permissionRepo     auth.PermissionRepository
	userRoleRepo       auth.UserRoleRepository
	rolePermissionRepo auth.RolePermissionRepository
}

func NewAuthorizationService(
	userRepo user.UserRepository,
	roleRepo auth.RoleRepository,
	permissionRepo auth.PermissionRepository,
	userRoleRepo auth.UserRoleRepository,
	rolePermissionRepo auth.RolePermissionRepository,
) contracts.AuthorizationService {
	return &authorizationService{
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		userRoleRepo:       userRoleRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

func (s *authorizationService) UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return false, apperrors.NewInternalError("failed to get user roles", err)
	}

	if len(userRoles) == 0 {
		return false, nil
	}

	for _, userRole := range userRoles {
		hasPermission, err := s.rolePermissionRepo.RoleHasResourceAction(ctx, userRole.RoleID, resource, action)
		if err != nil {
			return false, apperrors.NewInternalError("failed to check role permission", err)
		}
		if hasPermission {
			return true, nil
		}
	}

	return false, nil
}

func (s *authorizationService) UserHasPermissionByName(ctx context.Context, userID, permissionName string) (bool, error) {
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(userRoles) == 0 {
		return false, nil
	}

	for _, userRole := range userRoles {
		hasPermission, err := s.rolePermissionRepo.RoleHasPermissionByName(ctx, userRole.RoleID, permissionName)
		if err != nil {
			return false, fmt.Errorf("failed to check role permission by name: %w", err)
		}
		if hasPermission {
			return true, nil
		}
	}

	return false, nil
}

func (s *authorizationService) UserHasRole(ctx context.Context, userID, roleName string) (bool, error) {
	return s.userRoleRepo.UserHasRoleName(ctx, userID, roleName)
}

func (s *authorizationService) UserHasAnyRole(ctx context.Context, userID string, roleNames []string) (bool, error) {
	for _, roleName := range roleNames {
		hasRole, err := s.userRoleRepo.UserHasRoleName(ctx, userID, roleName)
		if err != nil {
			return false, fmt.Errorf("failed to check role %s: %w", roleName, err)
		}
		if hasRole {
			return true, nil
		}
	}
	return false, nil
}

func (s *authorizationService) UserHasAllRoles(ctx context.Context, userID string, roleNames []string) (bool, error) {
	for _, roleName := range roleNames {
		hasRole, err := s.userRoleRepo.UserHasRoleName(ctx, userID, roleName)
		if err != nil {
			return false, fmt.Errorf("failed to check role %s: %w", roleName, err)
		}
		if !hasRole {
			return false, nil
		}
	}
	return true, nil
}

func (s *authorizationService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get user roles", err)
	}

	if len(userRoles) == 0 {
		return []string{}, nil
	}

	permissionsMap := make(map[string]bool)

	for _, userRole := range userRoles {
		rolePermissions, err := s.rolePermissionRepo.GetActiveRolePermissions(ctx, userRole.RoleID)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to get role permissions", err)
		}

		for _, rolePerm := range rolePermissions {
			permission, err := s.permissionRepo.GetByID(ctx, rolePerm.PermissionID)
			if err != nil {
				return nil, apperrors.NewInternalError("failed to get permission", err)
			}
			if permission != nil && permission.IsActive() {
				permissionsMap[permission.Name().String()] = true
			}
		}
	}

	permissions := make([]string, 0, len(permissionsMap))
	for permission := range permissionsMap {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (s *authorizationService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get user roles", err)
	}

	if len(userRoles) == 0 {
		return []string{}, nil
	}

	roles := make([]string, 0, len(userRoles))

	for _, userRole := range userRoles {
		role, err := s.roleRepo.GetByID(ctx, userRole.RoleID)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to get role", err)
		}
		if role != nil && role.IsActive() {
			roles = append(roles, role.Name().String())
		}
	}

	return roles, nil
}

func (s *authorizationService) CheckMultiplePermissions(ctx context.Context, userID string, permissions []string) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, permission := range permissions {
		hasPermission, err := s.UserHasPermissionByName(ctx, userID, permission)
		if err != nil {
			return nil, apperrors.NewInternalError(fmt.Sprintf("failed to check permission %s", permission), err)
		}
		result[permission] = hasPermission
	}

	return result, nil
}

func (s *authorizationService) CanAccessResource(ctx context.Context, userID, resource, action string) error {
	hasPermission, err := s.UserHasPermission(ctx, userID, resource, action)
	if err != nil {
		return apperrors.NewInternalError("failed to check permission", err)
	}

	if !hasPermission {
		return apperrors.NewForbiddenError(fmt.Sprintf("access denied: user does not have permission to %s on %s", action, resource))
	}

	return nil
}

func (s *authorizationService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return s.UserHasRole(ctx, userID, "admin")
}

func (s *authorizationService) IsModerator(ctx context.Context, userID string) (bool, error) {
	return s.UserHasAnyRole(ctx, userID, []string{"admin", "moderator"})
}

func (s *authorizationService) GetEffectivePermissions(ctx context.Context, userID string) ([]contracts.PermissionInfo, error) {
	userRoles, err := s.userRoleRepo.GetActiveUserRoles(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get user roles", err)
	}

	if len(userRoles) == 0 {
		return []contracts.PermissionInfo{}, nil
	}

	permissionsMap := make(map[string]contracts.PermissionInfo)

	for _, userRole := range userRoles {
		rolePermissions, err := s.rolePermissionRepo.GetActiveRolePermissions(ctx, userRole.RoleID)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to get role permissions", err)
		}

		role, err := s.roleRepo.GetByID(ctx, userRole.RoleID)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to get role", err)
		}

		roleName := "unknown"
		if role != nil {
			roleName = role.Name().String()
		}

		for _, rolePerm := range rolePermissions {
			permission, err := s.permissionRepo.GetByID(ctx, rolePerm.PermissionID)
			if err != nil {
				return nil, apperrors.NewInternalError("failed to get permission", err)
			}
			if permission != nil && permission.IsActive() {
				permissionName := permission.Name().String()
				if _, exists := permissionsMap[permissionName]; !exists {
					permissionsMap[permissionName] = contracts.PermissionInfo{
						ID:          permission.ID().String(),
						Name:        permission.Name().String(),
						Resource:    permission.Resource().String(),
						Action:      permission.Action().String(),
						Description: permission.Description(),
						GrantedBy:   roleName,
					}
				}
			}
		}
	}

	permissions := make([]contracts.PermissionInfo, 0, len(permissionsMap))
	for _, permission := range permissionsMap {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

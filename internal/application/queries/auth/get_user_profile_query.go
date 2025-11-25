package auth

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/services"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

// GetUserProfileQuery represents the get user profile query
type GetUserProfileQuery struct {
	UserID string `validate:"required"`
}

// GetUserProfileQueryHandler handles the GetUserProfileQuery
type GetUserProfileQueryHandler struct {
	userRepo             repositories.UserRepository
	roleRepo             repositories.RoleRepository
	authorizationService services.AuthorizationService
}

// NewGetUserProfileQueryHandler creates a new GetUserProfileQueryHandler
func NewGetUserProfileQueryHandler(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	authorizationService services.AuthorizationService,
) *GetUserProfileQueryHandler {
	return &GetUserProfileQueryHandler{
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		authorizationService: authorizationService,
	}
}

// Handle executes the GetUserProfileQuery
func (h *GetUserProfileQueryHandler) Handle(ctx context.Context, query GetUserProfileQuery) (*dto.UserProfileResponse, error) {
	// Get user by ID
	user, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrors.NewNotFoundError("user not found")
	}

	// Get user roles
	userRoles, err := h.roleRepo.GetRolesByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	// Get effective permissions for the user
	permissions, err := h.authorizationService.GetEffectivePermissions(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	// Convert roles to string slice
	roleNames := make([]string, len(userRoles))
	for i, role := range userRoles {
		roleNames[i] = role.Name().String()
	}

	// Convert to permission info DTOs
	permissionInfos := make([]dto.PermissionInfoDTO, len(permissions))
	for i, perm := range permissions {
		permissionInfos[i] = dto.PermissionInfoDTO{
			ID:          perm.ID,
			Name:        perm.Name,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
			GrantedBy:   perm.GrantedBy,
		}
	}

	// Create response DTO
	return &dto.UserProfileResponse{
		User:        dto.ToAuthUserDTO(user),
		Roles:       roleNames,
		Permissions: permissionInfos,
	}, nil
}

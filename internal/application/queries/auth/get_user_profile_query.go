package auth

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

type GetUserProfileQuery struct {
	UserID string `validate:"required"`
}

type GetUserProfileQueryHandler struct {
	userRepo             user.UserRepository
	roleRepo             auth.RoleRepository
	authorizationService contracts.AuthorizationService
}

func NewGetUserProfileQueryHandler(
	userRepo user.UserRepository,
	roleRepo auth.RoleRepository,
	authorizationService contracts.AuthorizationService,
) *GetUserProfileQueryHandler {
	return &GetUserProfileQueryHandler{
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		authorizationService: authorizationService,
	}
}

func (h *GetUserProfileQueryHandler) Handle(ctx context.Context, query GetUserProfileQuery) (*dto.UserProfileResponse, error) {
	user, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperrors.NewNotFoundError("user not found")
	}

	userRoles, err := h.roleRepo.GetRolesByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	permissions, err := h.authorizationService.GetEffectivePermissions(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, len(userRoles))
	for i, role := range userRoles {
		roleNames[i] = role.Name().String()
	}

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

	return &dto.UserProfileResponse{
		User:        dto.ToAuthUserDTO(user),
		Roles:       roleNames,
		Permissions: permissionInfos,
	}, nil
}

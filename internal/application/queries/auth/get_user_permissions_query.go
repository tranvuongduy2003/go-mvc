package auth

import (
	"context"

	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
)

// GetUserPermissionsQuery represents the get user permissions query
type GetUserPermissionsQuery struct {
	UserID string `validate:"required"`
}

// GetUserPermissionsQueryHandler handles the GetUserPermissionsQuery
type GetUserPermissionsQueryHandler struct {
	authorizationService contracts.AuthorizationService
}

// NewGetUserPermissionsQueryHandler creates a new GetUserPermissionsQueryHandler
func NewGetUserPermissionsQueryHandler(authorizationService contracts.AuthorizationService) *GetUserPermissionsQueryHandler {
	return &GetUserPermissionsQueryHandler{
		authorizationService: authorizationService,
	}
}

// Handle executes the GetUserPermissionsQuery
func (h *GetUserPermissionsQueryHandler) Handle(ctx context.Context, query GetUserPermissionsQuery) ([]dto.PermissionInfoDTO, error) {
	// Get user effective permissions
	permissions, err := h.authorizationService.GetEffectivePermissions(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	permissionDTOs := make([]dto.PermissionInfoDTO, len(permissions))
	for i, perm := range permissions {
		permissionDTOs[i] = dto.PermissionInfoDTO{
			ID:          perm.ID,
			Name:        perm.Name,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
			GrantedBy:   perm.GrantedBy,
		}
	}

	return permissionDTOs, nil
}

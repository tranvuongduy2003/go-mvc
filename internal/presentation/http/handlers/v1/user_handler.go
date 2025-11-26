package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	userValidators "github.com/tranvuongduy2003/go-mvc/internal/application/validators/user"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"github.com/tranvuongduy2003/go-mvc/pkg/response"
)

type UserHandler struct {
	userService   *services.UserService
	userValidator userValidators.IUserValidator
}

func NewUserHandler(userService *services.UserService, userValidator userValidators.IUserValidator) *UserHandler {
	return &UserHandler{
		userService:   userService,
		userValidator: userValidator,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userDto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if validationErrors := h.userValidator.ValidateCreateUserRequest(req); len(validationErrors) > 0 {
		response.ValidationError(c, validationErrors)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req userDto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if validationErrors := h.userValidator.ValidateUpdateUserRequest(req); len(validationErrors) > 0 {
		response.ValidationError(c, validationErrors)
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if val, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &val
		}
	}

	req := userDto.ListUsersRequest{
		Page:     page,
		Limit:    limit,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
		IsActive: isActive,
	}

	if validationErrors := h.userValidator.ValidateListUsersRequest(req); len(validationErrors) > 0 {
		response.ValidationError(c, validationErrors)
		return
	}

	result, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	pag := &pagination.Pagination{
		Page:     result.Pagination.Page,
		PageSize: result.Pagination.PageSize,
		Total:    result.Pagination.Total,
		Pages:    result.Pagination.Pages,
	}

	response.SuccessWithPagination(c, result.Users, pag)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, apperrors.NewValidationError("User ID is required", nil))
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.Error(c, apperrors.NewValidationError("Avatar file is required", err))
		return
	}
	defer file.Close()

	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxFileSize {
		response.Error(c, apperrors.NewValidationError("File size exceeds 5MB limit", nil))
		return
	}

	user, err := h.userService.UploadAvatar(c.Request.Context(), id, file, header)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

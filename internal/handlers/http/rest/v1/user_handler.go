package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	userValidators "github.com/tranvuongduy2003/go-mvc/internal/application/validators/user"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
	"github.com/tranvuongduy2003/go-mvc/pkg/response"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService   *services.UserService
	userValidator userValidators.IUserValidator
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService, userValidator userValidators.IUserValidator) *UserHandler {
	return &UserHandler{
		userService:   userService,
		userValidator: userValidator,
	}
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with email, name, phone and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User creation data"
// @Success 201 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userDto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Validate request
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

// GetUserByID retrieves a user by ID
// @Summary Get user by ID
// @Description Get a user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id} [get]
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

// UpdateUser updates an existing user
// @Summary Update user
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User update data"
// @Success 200 {object} response.APIResponse{data=dto.UserResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id} [put]
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

	// Validate request
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

// DeleteUser deletes a user
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users/{id} [delete]
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

// ListUsers retrieves a paginated list of users
// @Summary List users
// @Description Get a paginated list of users with optional filters
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_dir query string false "Sort direction (asc/desc)" default(desc)
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} response.APIResponse{data=dto.ListUsersResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
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

	// Validate request
	if validationErrors := h.userValidator.ValidateListUsersRequest(req); len(validationErrors) > 0 {
		response.ValidationError(c, validationErrors)
		return
	}

	result, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Convert pagination DTO to response pagination
	pag := &pagination.Pagination{
		Page:     result.Pagination.Page,
		PageSize: result.Pagination.PageSize,
		Total:    result.Pagination.Total,
		Pages:    result.Pagination.Pages,
	}

	response.SuccessWithPagination(c, result.Users, pag)
}

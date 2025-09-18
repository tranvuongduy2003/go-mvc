package validators

import (
	"regexp"
	"strings"

	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// UserValidator provides validation logic for user-related operations
type UserValidator struct {
	logger *logger.Logger
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// NewUserValidator creates a new user validator
func NewUserValidator(logger *logger.Logger) *UserValidator {
	return &UserValidator{
		logger: logger,
	}
}

// ValidateCreateUserRequest validates create user request
func (v *UserValidator) ValidateCreateUserRequest(req *dto.CreateUserRequest) ValidationErrors {
	var errors ValidationErrors

	// Validate email
	if err := v.validateEmail(req.Email); err != nil {
		errors = append(errors, *err)
	}

	// Validate username
	if err := v.validateUsername(req.Username); err != nil {
		errors = append(errors, *err)
	}

	// Validate password
	if err := v.validatePassword(req.Password); err != nil {
		errors = append(errors, *err)
	}

	// Validate first name
	if err := v.validateFirstName(req.FirstName); err != nil {
		errors = append(errors, *err)
	}

	// Validate last name
	if err := v.validateLastName(req.LastName); err != nil {
		errors = append(errors, *err)
	}

	return errors
}

// ValidateUpdateUserRequest validates update user request
func (v *UserValidator) ValidateUpdateUserRequest(req *dto.UpdateUserRequest) ValidationErrors {
	var errors ValidationErrors

	// Validate email if provided
	if req.Email != nil {
		if err := v.validateEmail(*req.Email); err != nil {
			errors = append(errors, *err)
		}
	}

	// Validate first name if provided
	if req.FirstName != nil {
		if err := v.validateFirstName(*req.FirstName); err != nil {
			errors = append(errors, *err)
		}
	}

	// Validate last name if provided
	if req.LastName != nil {
		if err := v.validateLastName(*req.LastName); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// ValidateChangePasswordRequest validates change password request
func (v *UserValidator) ValidateChangePasswordRequest(req *dto.ChangePasswordRequest) ValidationErrors {
	var errors ValidationErrors

	// Validate current password
	if req.CurrentPassword == "" {
		errors = append(errors, ValidationError{
			Field:   "current_password",
			Message: "Current password is required",
			Code:    "required",
		})
	}

	// Validate new password
	if err := v.validatePassword(req.NewPassword); err != nil {
		err.Field = "new_password"
		errors = append(errors, *err)
	}

	// Check if new password is different from current
	if req.CurrentPassword == req.NewPassword {
		errors = append(errors, ValidationError{
			Field:   "new_password",
			Message: "New password must be different from current password",
			Code:    "same_password",
		})
	}

	return errors
}

// ValidateLoginRequest validates login request
func (v *UserValidator) ValidateLoginRequest(req *dto.LoginRequest) ValidationErrors {
	var errors ValidationErrors

	// Validate email
	if err := v.validateEmail(req.Email); err != nil {
		errors = append(errors, *err)
	}

	// Validate password
	if req.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
			Code:    "required",
		})
	}

	return errors
}

// ValidateUserListRequest validates user list request
func (v *UserValidator) ValidateUserListRequest(req *dto.UserListRequest) ValidationErrors {
	var errors ValidationErrors

	// Validate page
	if req.Page <= 0 {
		errors = append(errors, ValidationError{
			Field:   "page",
			Message: "Page must be greater than 0",
			Code:    "min_value",
		})
	}

	// Validate limit
	if req.Limit <= 0 {
		errors = append(errors, ValidationError{
			Field:   "limit",
			Message: "Limit must be greater than 0",
			Code:    "min_value",
		})
	}

	if req.Limit > 100 {
		errors = append(errors, ValidationError{
			Field:   "limit",
			Message: "Limit cannot exceed 100",
			Code:    "max_value",
		})
	}

	// Validate sort field
	if req.Sort != "" {
		validSortFields := []string{"id", "email", "username", "first_name", "last_name", "created_at", "updated_at"}
		if !v.isValidSortField(req.Sort, validSortFields) {
			errors = append(errors, ValidationError{
				Field:   "sort",
				Message: "Invalid sort field",
				Code:    "invalid_value",
			})
		}
	}

	// Validate order
	if req.Order != "" && req.Order != "asc" && req.Order != "desc" {
		errors = append(errors, ValidationError{
			Field:   "order",
			Message: "Order must be 'asc' or 'desc'",
			Code:    "invalid_value",
		})
	}

	return errors
}

// ValidateChangeRoleRequest validates change role request
func (v *UserValidator) ValidateChangeRoleRequest(req *dto.ChangeRoleRequest) ValidationErrors {
	var errors ValidationErrors

	// Validate role
	validRoles := []string{"admin", "user", "moderator"}
	if !v.contains(validRoles, req.Role) {
		errors = append(errors, ValidationError{
			Field:   "role",
			Message: "Invalid role. Must be one of: admin, user, moderator",
			Code:    "invalid_value",
		})
	}

	return errors
}

// Private validation methods

func (v *UserValidator) validateEmail(email string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "Email is required",
			Code:    "required",
		}
	}

	// Basic email validation using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   "email",
			Message: "Invalid email format",
			Code:    "invalid_format",
		}
	}

	return nil
}

func (v *UserValidator) validateUsername(username string) *ValidationError {
	if username == "" {
		return &ValidationError{
			Field:   "username",
			Message: "Username is required",
			Code:    "required",
		}
	}

	if len(username) < 3 {
		return &ValidationError{
			Field:   "username",
			Message: "Username must be at least 3 characters long",
			Code:    "min_length",
		}
	}

	if len(username) > 50 {
		return &ValidationError{
			Field:   "username",
			Message: "Username cannot exceed 50 characters",
			Code:    "max_length",
		}
	}

	// Username should contain only alphanumeric characters and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return &ValidationError{
			Field:   "username",
			Message: "Username can only contain alphanumeric characters and underscores",
			Code:    "invalid_format",
		}
	}

	return nil
}

func (v *UserValidator) validatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   "password",
			Message: "Password is required",
			Code:    "required",
		}
	}

	if len(password) < 8 {
		return &ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
			Code:    "min_length",
		}
	}

	if len(password) > 128 {
		return &ValidationError{
			Field:   "password",
			Message: "Password cannot exceed 128 characters",
			Code:    "max_length",
		}
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one uppercase letter",
			Code:    "missing_uppercase",
		}
	}

	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one lowercase letter",
			Code:    "missing_lowercase",
		}
	}

	// Check for at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one digit",
			Code:    "missing_digit",
		}
	}

	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return &ValidationError{
			Field:   "password",
			Message: "Password must contain at least one special character",
			Code:    "missing_special",
		}
	}

	return nil
}

func (v *UserValidator) validateFirstName(firstName string) *ValidationError {
	if firstName == "" {
		return &ValidationError{
			Field:   "first_name",
			Message: "First name is required",
			Code:    "required",
		}
	}

	if len(firstName) < 2 {
		return &ValidationError{
			Field:   "first_name",
			Message: "First name must be at least 2 characters long",
			Code:    "min_length",
		}
	}

	if len(firstName) > 50 {
		return &ValidationError{
			Field:   "first_name",
			Message: "First name cannot exceed 50 characters",
			Code:    "max_length",
		}
	}

	// First name should contain only letters and spaces
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(firstName) {
		return &ValidationError{
			Field:   "first_name",
			Message: "First name can only contain letters and spaces",
			Code:    "invalid_format",
		}
	}

	return nil
}

func (v *UserValidator) validateLastName(lastName string) *ValidationError {
	if lastName == "" {
		return &ValidationError{
			Field:   "last_name",
			Message: "Last name is required",
			Code:    "required",
		}
	}

	if len(lastName) < 2 {
		return &ValidationError{
			Field:   "last_name",
			Message: "Last name must be at least 2 characters long",
			Code:    "min_length",
		}
	}

	if len(lastName) > 50 {
		return &ValidationError{
			Field:   "last_name",
			Message: "Last name cannot exceed 50 characters",
			Code:    "max_length",
		}
	}

	// Last name should contain only letters and spaces
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(lastName) {
		return &ValidationError{
			Field:   "last_name",
			Message: "Last name can only contain letters and spaces",
			Code:    "invalid_format",
		}
	}

	return nil
}

// Helper methods

func (v *UserValidator) isValidSortField(field string, validFields []string) bool {
	return v.contains(validFields, field)
}

func (v *UserValidator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

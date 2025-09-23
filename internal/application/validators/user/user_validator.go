package user

import (
	"regexp"
	"strings"

	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
)

// IUserValidator defines the interface for user validation operations
type IUserValidator interface {
	ValidateCreateUserRequest(req userDto.CreateUserRequest) map[string]string
	ValidateUpdateUserRequest(req userDto.UpdateUserRequest) map[string]string
	ValidateListUsersRequest(req userDto.ListUsersRequest) map[string]string
}

// UserValidator handles validation for user-related operations
type UserValidator struct{}

// NewUserValidator creates a new UserValidator instance
func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateCreateUserRequest validates the create user request
func (v *UserValidator) ValidateCreateUserRequest(req userDto.CreateUserRequest) map[string]string {
	errors := make(map[string]string)

	// Validate email
	if req.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(req.Email) {
		errors["email"] = "Email format is invalid"
	}

	// Validate name
	if req.Name == "" {
		errors["name"] = "Name is required"
	} else if len(strings.TrimSpace(req.Name)) < 2 {
		errors["name"] = "Name must be at least 2 characters long"
	} else if len(req.Name) > 100 {
		errors["name"] = "Name must not exceed 100 characters"
	}

	// Validate password
	if req.Password == "" {
		errors["password"] = "Password is required"
	} else if len(req.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	} else if !isValidPassword(req.Password) {
		errors["password"] = "Password must contain at least one uppercase letter, one lowercase letter, and one number"
	}

	// Validate phone (optional)
	if req.Phone != "" && !isValidPhone(req.Phone) {
		errors["phone"] = "Phone number format is invalid"
	}

	return errors
}

// ValidateUpdateUserRequest validates the update user request
func (v *UserValidator) ValidateUpdateUserRequest(req userDto.UpdateUserRequest) map[string]string {
	errors := make(map[string]string)

	// Validate name
	if req.Name == "" {
		errors["name"] = "Name is required"
	} else if len(strings.TrimSpace(req.Name)) < 2 {
		errors["name"] = "Name must be at least 2 characters long"
	} else if len(req.Name) > 100 {
		errors["name"] = "Name must not exceed 100 characters"
	}

	// Validate phone (optional)
	if req.Phone != "" && !isValidPhone(req.Phone) {
		errors["phone"] = "Phone number format is invalid"
	}

	return errors
}

// ValidateListUsersRequest validates the list users request
func (v *UserValidator) ValidateListUsersRequest(req userDto.ListUsersRequest) map[string]string {
	errors := make(map[string]string)

	// Validate page
	if req.Page < 1 {
		errors["page"] = "Page must be greater than 0"
	}

	// Validate limit
	if req.Limit < 1 {
		errors["limit"] = "Limit must be greater than 0"
	} else if req.Limit > 100 {
		errors["limit"] = "Limit must not exceed 100"
	}

	// Validate sort direction
	if req.SortDir != "" && req.SortDir != "asc" && req.SortDir != "desc" {
		errors["sort_dir"] = "Sort direction must be 'asc' or 'desc'"
	}

	// Validate sort by field
	if req.SortBy != "" && !isValidSortField(req.SortBy) {
		errors["sort_by"] = "Invalid sort field"
	}

	return errors
}

// Helper functions

// isValidEmail validates email format
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// isValidPassword validates password strength
func isValidPassword(password string) bool {
	// At least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// At least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// At least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

// isValidPhone validates phone number format
func isValidPhone(phone string) bool {
	// Vietnamese phone number pattern
	phoneRegex := `^(\+84|84|0)[1-9][0-9]{8,9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

// isValidSortField validates sort field
func isValidSortField(field string) bool {
	validFields := []string{"id", "email", "name", "created_at", "updated_at", "is_active"}
	for _, validField := range validFields {
		if field == validField {
			return true
		}
	}
	return false
}

package user

import (
	"regexp"
	"strings"

	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
)

type IUserValidator interface {
	ValidateCreateUserRequest(req userDto.CreateUserRequest) map[string]string
	ValidateUpdateUserRequest(req userDto.UpdateUserRequest) map[string]string
	ValidateListUsersRequest(req userDto.ListUsersRequest) map[string]string
}

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (v *UserValidator) ValidateCreateUserRequest(req userDto.CreateUserRequest) map[string]string {
	errors := make(map[string]string)

	if req.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(req.Email) {
		errors["email"] = "Email format is invalid"
	}

	if req.Name == "" {
		errors["name"] = "Name is required"
	} else if len(strings.TrimSpace(req.Name)) < 2 {
		errors["name"] = "Name must be at least 2 characters long"
	} else if len(req.Name) > 100 {
		errors["name"] = "Name must not exceed 100 characters"
	}

	if req.Password == "" {
		errors["password"] = "Password is required"
	} else if len(req.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	} else if !isValidPassword(req.Password) {
		errors["password"] = "Password must contain at least one uppercase letter, one lowercase letter, and one number"
	}

	if req.Phone != "" && !isValidPhone(req.Phone) {
		errors["phone"] = "Phone number format is invalid"
	}

	return errors
}

func (v *UserValidator) ValidateUpdateUserRequest(req userDto.UpdateUserRequest) map[string]string {
	errors := make(map[string]string)

	if req.Name == "" {
		errors["name"] = "Name is required"
	} else if len(strings.TrimSpace(req.Name)) < 2 {
		errors["name"] = "Name must be at least 2 characters long"
	} else if len(req.Name) > 100 {
		errors["name"] = "Name must not exceed 100 characters"
	}

	if req.Phone != "" && !isValidPhone(req.Phone) {
		errors["phone"] = "Phone number format is invalid"
	}

	return errors
}

func (v *UserValidator) ValidateListUsersRequest(req userDto.ListUsersRequest) map[string]string {
	errors := make(map[string]string)

	if req.Page < 1 {
		errors["page"] = "Page must be greater than 0"
	}

	if req.Limit < 1 {
		errors["limit"] = "Limit must be greater than 0"
	} else if req.Limit > 100 {
		errors["limit"] = "Limit must not exceed 100"
	}

	if req.SortDir != "" && req.SortDir != "asc" && req.SortDir != "desc" {
		errors["sort_dir"] = "Sort direction must be 'asc' or 'desc'"
	}

	if req.SortBy != "" && !isValidSortField(req.SortBy) {
		errors["sort_by"] = "Invalid sort field"
	}

	return errors
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit
}

func isValidPhone(phone string) bool {
	phoneRegex := `^(\+84|84|0)[1-9][0-9]{8,9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

func isValidSortField(field string) bool {
	validFields := []string{"id", "email", "name", "created_at", "updated_at", "is_active"}
	for _, validField := range validFields {
		if field == validField {
			return true
		}
	}
	return false
}

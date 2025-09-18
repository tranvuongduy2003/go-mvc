package user

import "errors"

// Domain errors for User
var (
	// Validation errors
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInvalidFirstName = errors.New("invalid first name")
	ErrInvalidLastName  = errors.New("invalid last name")
	ErrInvalidRole      = errors.New("invalid role")
	ErrInvalidPhone     = errors.New("invalid phone number")

	// Business logic errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserDeleted        = errors.New("user is deleted")
	ErrUserAlreadyDeleted = errors.New("user is already deleted")
	ErrUserInactive       = errors.New("user is inactive")
	ErrUserAlreadyActive  = errors.New("user is already active")
	ErrUserNotActivated   = errors.New("user is not activated")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordMismatch   = errors.New("password mismatch")
	ErrWeakPassword       = errors.New("password is too weak")
	ErrPasswordExpired    = errors.New("password has expired")

	// Authorization errors
	ErrUnauthorized           = errors.New("unauthorized")
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrAccessDenied           = errors.New("access denied")

	// Duplicate errors
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// Profile errors
	ErrProfileNotFound      = errors.New("profile not found")
	ErrInvalidProfileData   = errors.New("invalid profile data")
	ErrProfileAlreadyExists = errors.New("profile already exists")

	// Session errors
	ErrSessionExpired  = errors.New("session expired")
	ErrInvalidSession  = errors.New("invalid session")
	ErrSessionNotFound = errors.New("session not found")

	// Token errors
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenNotFound = errors.New("token not found")
)

package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/shared/valueobject"
)

// User represents a user entity
type User struct {
	id        uuid.UUID
	email     valueobject.Email
	username  string
	password  string
	firstName string
	lastName  string
	isActive  bool
	role      Role
	profile   *Profile
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

// Role represents user roles
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

// Profile represents user profile information
type Profile struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Avatar      string
	Bio         string
	DateOfBirth *time.Time
	Phone       string
	Address     string
	City        string
	Country     string
	Website     string
	SocialLinks map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewUser creates a new user
func NewUser(email valueobject.Email, username, hashedPassword, firstName, lastName string) (*User, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	if hashedPassword == "" {
		return nil, ErrInvalidPassword
	}
	if firstName == "" {
		return nil, ErrInvalidFirstName
	}

	now := time.Now().UTC()

	return &User{
		id:        uuid.New(),
		email:     email,
		username:  username,
		password:  hashedPassword,
		firstName: firstName,
		lastName:  lastName,
		isActive:  true,
		role:      RoleUser,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstruct creates a user from stored data (for repository use)
func Reconstruct(id uuid.UUID, email valueobject.Email, username, password, firstName, lastName string,
	isActive bool, role Role, createdAt, updatedAt time.Time, deletedAt *time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		username:  username,
		password:  password,
		firstName: firstName,
		lastName:  lastName,
		isActive:  isActive,
		role:      role,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

// Getters
func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) FullName() string {
	if u.lastName == "" {
		return u.firstName
	}
	return u.firstName + " " + u.lastName
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) Profile() *Profile {
	return u.profile
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

// Business methods
func (u *User) UpdateEmail(email valueobject.Email) error {
	if u.deletedAt != nil {
		return ErrUserDeleted
	}

	u.email = email
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) UpdatePassword(hashedPassword string) error {
	if u.deletedAt != nil {
		return ErrUserDeleted
	}
	if hashedPassword == "" {
		return ErrInvalidPassword
	}

	u.password = hashedPassword
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) UpdateProfile(firstName, lastName string) error {
	if u.deletedAt != nil {
		return ErrUserDeleted
	}
	if firstName == "" {
		return ErrInvalidFirstName
	}

	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ChangeRole(role Role) error {
	if u.deletedAt != nil {
		return ErrUserDeleted
	}
	if !isValidRole(role) {
		return ErrInvalidRole
	}

	u.role = role
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) Activate() error {
	if u.deletedAt != nil {
		return ErrUserDeleted
	}

	u.isActive = true
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) Deactivate() error {
	if u.deletedAt != nil {
		return ErrUserDeleted
	}

	u.isActive = false
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) Delete() error {
	if u.deletedAt != nil {
		return ErrUserAlreadyDeleted
	}

	now := time.Now().UTC()
	u.deletedAt = &now
	u.isActive = false
	u.updatedAt = now
	return nil
}

func (u *User) IsDeleted() bool {
	return u.deletedAt != nil
}

func (u *User) HasRole(role Role) bool {
	return u.role == role
}

func (u *User) IsAdmin() bool {
	return u.role == RoleAdmin
}

func (u *User) IsModerator() bool {
	return u.role == RoleModerator || u.role == RoleAdmin
}

func (u *User) CanManageUsers() bool {
	return u.role == RoleAdmin
}

func (u *User) CanModerateContent() bool {
	return u.role == RoleModerator || u.role == RoleAdmin
}

func (u *User) SetProfile(profile *Profile) {
	u.profile = profile
	u.updatedAt = time.Now().UTC()
}

// Helper functions
func isValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleModerator:
		return true
	default:
		return false
	}
}

// String returns string representation of role
func (r Role) String() string {
	return string(r)
}

// IsValid checks if role is valid
func (r Role) IsValid() bool {
	return isValidRole(r)
}

// NewProfile creates a new user profile
func NewProfile(userID uuid.UUID) *Profile {
	now := time.Now().UTC()
	return &Profile{
		ID:          uuid.New(),
		UserID:      userID,
		SocialLinks: make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateAvatar updates user avatar
func (p *Profile) UpdateAvatar(avatar string) {
	p.Avatar = avatar
	p.UpdatedAt = time.Now().UTC()
}

// UpdateBio updates user bio
func (p *Profile) UpdateBio(bio string) {
	p.Bio = bio
	p.UpdatedAt = time.Now().UTC()
}

// UpdatePersonalInfo updates personal information
func (p *Profile) UpdatePersonalInfo(phone, address, city, country, website string, dateOfBirth *time.Time) {
	p.Phone = phone
	p.Address = address
	p.City = city
	p.Country = country
	p.Website = website
	p.DateOfBirth = dateOfBirth
	p.UpdatedAt = time.Now().UTC()
}

// AddSocialLink adds a social media link
func (p *Profile) AddSocialLink(platform, url string) {
	if p.SocialLinks == nil {
		p.SocialLinks = make(map[string]string)
	}
	p.SocialLinks[platform] = url
	p.UpdatedAt = time.Now().UTC()
}

// RemoveSocialLink removes a social media link
func (p *Profile) RemoveSocialLink(platform string) {
	if p.SocialLinks != nil {
		delete(p.SocialLinks, platform)
		p.UpdatedAt = time.Now().UTC()
	}
}

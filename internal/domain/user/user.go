package user

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user aggregate root
type User struct {
	id        UserID
	email     Email
	name      Name
	phone     Phone
	password  Password
	avatar    Avatar
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
	version   int64
	events    []events.DomainEvent
}

// UserID value object
type UserID struct {
	value string
}

// NewUserID creates a new user ID
func NewUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// NewUserIDFromString creates a user ID from string
func NewUserIDFromString(id string) (UserID, error) {
	if id == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}
	if _, err := uuid.Parse(id); err != nil {
		return UserID{}, errors.New("invalid user ID format")
	}
	return UserID{value: id}, nil
}

// String returns the string representation of user ID
func (id UserID) String() string {
	return id.value
}

// Email value object
type Email struct {
	value string
}

// NewEmail creates a new email value object
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: email}, nil
}

// String returns the string representation of email
func (e Email) String() string {
	return e.value
}

// Name value object
type Name struct {
	value string
}

// NewName creates a new name value object
func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Name{}, errors.New("name cannot be empty")
	}
	if len(name) < 2 {
		return Name{}, errors.New("name must be at least 2 characters")
	}
	if len(name) > 100 {
		return Name{}, errors.New("name cannot exceed 100 characters")
	}
	return Name{value: name}, nil
}

// String returns the string representation of name
func (n Name) String() string {
	return n.value
}

// Phone value object
type Phone struct {
	value string
}

// NewPhone creates a new phone value object
func NewPhone(phone string) (Phone, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		// Phone is optional
		return Phone{}, nil
	}

	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-\(\)]{10,15}$`)
	if !phoneRegex.MatchString(phone) {
		return Phone{}, errors.New("invalid phone format")
	}

	return Phone{value: phone}, nil
}

// String returns the string representation of phone
func (p Phone) String() string {
	return p.value
}

// IsEmpty checks if phone is empty
func (p Phone) IsEmpty() bool {
	return p.value == ""
}

// Avatar value object
type Avatar struct {
	fileKey string
	cdnUrl  string
}

// NewAvatar creates a new avatar value object
func NewAvatar(fileKey, cdnUrl string) Avatar {
	return Avatar{
		fileKey: strings.TrimSpace(fileKey),
		cdnUrl:  strings.TrimSpace(cdnUrl),
	}
}

// FileKey returns the file key
func (a Avatar) FileKey() string {
	return a.fileKey
}

// CDNUrl returns the CDN URL
func (a Avatar) CDNUrl() string {
	return a.cdnUrl
}

// IsEmpty checks if avatar is empty
func (a Avatar) IsEmpty() bool {
	return a.fileKey == "" || a.cdnUrl == ""
}

// Password value object
type Password struct {
	hashedValue string
}

// NewPassword creates a new password value object
func NewPassword(plainPassword string) (Password, error) {
	if plainPassword == "" {
		return Password{}, errors.New("password cannot be empty")
	}
	if len(plainPassword) < 8 {
		return Password{}, errors.New("password must be at least 8 characters")
	}
	if len(plainPassword) > 72 {
		return Password{}, errors.New("password cannot exceed 72 characters")
	}

	// Hash the password using bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, errors.New("failed to hash password")
	}

	return Password{hashedValue: string(hashedBytes)}, nil
}

// NewHashedPassword creates a password from already hashed value
func NewHashedPassword(hashedPassword string) Password {
	return Password{hashedValue: hashedPassword}
}

// Hash returns the hashed password
func (p Password) Hash() string {
	return p.hashedValue
}

// VerifyPassword verifies if the provided password matches
func (p Password) VerifyPassword(plainPassword string) bool {
	// Use bcrypt to compare the password
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}

// User Domain Events
type UserCreated struct {
	*events.BaseDomainEvent
	UserID    string
	Email     string
	Name      string
	CreatedAt time.Time
}

type UserUpdated struct {
	*events.BaseDomainEvent
	UserID    string
	Email     string
	Name      string
	UpdatedAt time.Time
}

type UserDeleted struct {
	*events.BaseDomainEvent
	UserID    string
	DeletedAt time.Time
}

// NewUser creates a new user aggregate
func NewUser(email, name, phone, password string) (*User, error) {
	userID := NewUserID()

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	phoneVO, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}

	passwordVO, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &User{
		id:        userID,
		email:     emailVO,
		name:      nameVO,
		phone:     phoneVO,
		password:  passwordVO,
		avatar:    NewAvatar("", ""), // Empty avatar initially
		isActive:  true,
		createdAt: now,
		updatedAt: now,
		version:   1,
		events:    make([]events.DomainEvent, 0),
	}

	// Add domain event
	userUUID, _ := uuid.Parse(userID.String())
	user.addEvent(&UserCreated{
		BaseDomainEvent: events.NewBaseDomainEvent("UserCreated", userUUID, "User", map[string]interface{}{
			"user_id": userID.String(),
			"email":   emailVO.String(),
			"name":    nameVO.String(),
		}),
		UserID:    userID.String(),
		Email:     emailVO.String(),
		Name:      nameVO.String(),
		CreatedAt: now,
	})

	return user, nil
}

// ReconstructUser reconstructs a user from persistence
func ReconstructUser(id, email, name, phone, hashedPassword, avatarFileKey, avatarCDNUrl string, isActive bool, createdAt, updatedAt time.Time, version int64) (*User, error) {
	userID, err := NewUserIDFromString(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	phoneVO, err := NewPhone(phone)
	if err != nil {
		return nil, err
	}

	avatarVO := NewAvatar(avatarFileKey, avatarCDNUrl)

	return &User{
		id:        userID,
		email:     emailVO,
		name:      nameVO,
		phone:     phoneVO,
		password:  NewHashedPassword(hashedPassword),
		avatar:    avatarVO,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
		version:   version,
		events:    make([]events.DomainEvent, 0),
	}, nil
}

// Getters
func (u *User) ID() string {
	return u.id.String()
}

func (u *User) Email() string {
	return u.email.String()
}

func (u *User) Name() string {
	return u.name.String()
}

func (u *User) Phone() string {
	return u.phone.String()
}

func (u *User) HashedPassword() string {
	return u.password.Hash()
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) Version() int64 {
	return u.version
}

func (u *User) Avatar() Avatar {
	return u.avatar
}

// Business Methods
func (u *User) UpdateProfile(name, phone string) error {
	nameVO, err := NewName(name)
	if err != nil {
		return err
	}

	phoneVO, err := NewPhone(phone)
	if err != nil {
		return err
	}

	u.name = nameVO
	u.phone = phoneVO
	u.updatedAt = time.Now()
	u.version++

	// Add domain event
	userUUID, _ := uuid.Parse(u.id.String())
	u.addEvent(&UserUpdated{
		BaseDomainEvent: events.NewBaseDomainEvent("UserUpdated", userUUID, "User", map[string]interface{}{
			"user_id": u.id.String(),
			"email":   u.email.String(),
			"name":    u.name.String(),
		}),
		UserID:    u.id.String(),
		Email:     u.email.String(),
		Name:      u.name.String(),
		UpdatedAt: u.updatedAt,
	})

	return nil
}

func (u *User) ChangePassword(newPassword string) error {
	passwordVO, err := NewPassword(newPassword)
	if err != nil {
		return err
	}

	u.password = passwordVO
	u.updatedAt = time.Now()
	u.version++

	return nil
}

func (u *User) UpdateAvatar(fileKey, cdnUrl string) {
	u.avatar = NewAvatar(fileKey, cdnUrl)
	u.updatedAt = time.Now()
	u.version++
}

func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = time.Now()
	u.version++

	// Add domain event
	userUUID, _ := uuid.Parse(u.id.String())
	u.addEvent(&UserDeleted{
		BaseDomainEvent: events.NewBaseDomainEvent("UserDeleted", userUUID, "User", map[string]interface{}{
			"user_id": u.id.String(),
		}),
		UserID:    u.id.String(),
		DeletedAt: u.updatedAt,
	})
}

func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = time.Now()
	u.version++
}

func (u *User) VerifyPassword(password string) bool {
	return u.password.VerifyPassword(password)
}

// Domain Events Management
func (u *User) addEvent(event events.DomainEvent) {
	u.events = append(u.events, event)
}

func (u *User) GetEvents() []events.DomainEvent {
	return u.events
}

func (u *User) ClearEvents() {
	u.events = make([]events.DomainEvent, 0)
}

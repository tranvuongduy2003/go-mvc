# Chapter 3: Domain Layer

## Understanding the Domain Layer

The Domain Layer is the heart of your application. It contains the core business logic, entities, and rules that define your system's behavior. This layer is completely independent of external concerns like databases, frameworks, or UI.

## Domain Layer Principles

### Independence
The domain layer should not depend on any other layer. It contains:
- Pure business logic
- Domain entities and aggregates
- Value objects
- Domain events
- Repository interfaces (but not implementations)

### Rich Domain Model
Domain entities are not just data containers—they encapsulate business logic and enforce invariants.

**❌ Anemic Domain Model (Avoid)**
```go
type User struct {
    ID       uuid.UUID
    Email    string
    Password string
}

// Business logic in service layer
func (s *UserService) ChangeEmail(userID uuid.UUID, newEmail string) error {
    user := s.repo.GetByID(userID)
    if !isValidEmail(newEmail) {
        return errors.New("invalid email")
    }
    user.Email = newEmail
    return s.repo.Save(user)
}
```

**✅ Rich Domain Model (Preferred)**
```go
type User struct {
    id       UserID
    email    Email // Value object
    password HashedPassword
    events   []DomainEvent
}

// Business logic in entity
func (u *User) ChangeEmail(newEmail Email) error {
    if u.email.Equals(newEmail) {
        return ErrEmailNotChanged
    }
    
    oldEmail := u.email
    u.email = newEmail
    
    u.addEvent(UserEmailChanged{
        UserID:   u.id,
        OldEmail: oldEmail,
        NewEmail: newEmail,
    })
    
    return nil
}
```

## Directory Structure

```
internal/domain/
├── user/                      # User bounded context
│   ├── user.go               # User aggregate root
│   ├── user_id.go           # User ID value object
│   ├── email.go             # Email value object
│   ├── password.go          # Password value object
│   ├── events.go            # Domain events
│   └── errors.go            # Domain-specific errors
├── auth/                      # Auth bounded context
│   ├── role.go
│   ├── permission.go
│   └── role_permission.go
├── messaging/                 # Messaging bounded context
│   ├── inbox.go
│   ├── outbox.go
│   └── message.go
├── job/                       # Job bounded context
│   ├── job.go
│   └── job_status.go
├── repositories/              # Repository interfaces
│   ├── user_repository.go
│   ├── role_repository.go
│   └── permission_repository.go
├── contracts/                 # Service interfaces
│   ├── event_bus.go
│   └── message_broker.go
└── shared/                    # Shared domain concepts
    ├── events/
    │   └── domain_event.go
    └── specifications/
        └── specification.go
```

## Entities and Aggregates

### Entity Definition

An entity is an object that has a unique identity that runs through time and different representations.

```go
package user

import (
    "time"
    "github.com/google/uuid"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
)

// User is the aggregate root for user management
type User struct {
    id        UserID
    email     Email
    password  HashedPassword
    name      string
    isActive  bool
    roles     []Role
    createdAt time.Time
    updatedAt time.Time
    version   int64          // For optimistic locking
    events    []events.DomainEvent
}

// NewUser creates a new user (factory method)
func NewUser(email Email, password Password, name string) (*User, error) {
    if name == "" {
        return nil, ErrNameRequired
    }
    
    hashedPass, err := HashPassword(password)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        id:        NewUserID(),
        email:     email,
        password:  hashedPass,
        name:      name,
        isActive:  true,
        roles:     []Role{},
        createdAt: time.Now(),
        updatedAt: time.Now(),
        version:   1,
        events:    []events.DomainEvent{},
    }
    
    user.addEvent(UserCreated{
        UserID:    user.id,
        Email:     user.email,
        Name:      user.name,
        CreatedAt: user.createdAt,
    })
    
    return user, nil
}

// Business methods
func (u *User) Activate() error {
    if u.isActive {
        return ErrUserAlreadyActive
    }
    
    u.isActive = true
    u.updatedAt = time.Now()
    u.version++
    
    u.addEvent(UserActivated{
        UserID:      u.id,
        ActivatedAt: time.Now(),
    })
    
    return nil
}

func (u *User) Deactivate() error {
    if !u.isActive {
        return ErrUserAlreadyInactive
    }
    
    u.isActive = false
    u.updatedAt = time.Now()
    u.version++
    
    u.addEvent(UserDeactivated{
        UserID:        u.id,
        DeactivatedAt: time.Now(),
    })
    
    return nil
}

func (u *User) AssignRole(role Role) error {
    if u.HasRole(role) {
        return ErrRoleAlreadyAssigned
    }
    
    u.roles = append(u.roles, role)
    u.updatedAt = time.Now()
    u.version++
    
    u.addEvent(RoleAssigned{
        UserID:     u.id,
        Role:       role,
        AssignedAt: time.Now(),
    })
    
    return nil
}

func (u *User) HasRole(role Role) bool {
    for _, r := range u.roles {
        if r.ID().Equals(role.ID()) {
            return true
        }
    }
    return false
}

// Getters (no setters to enforce encapsulation)
func (u *User) ID() UserID                    { return u.id }
func (u *User) Email() Email                  { return u.email }
func (u *User) Name() string                  { return u.name }
func (u *User) IsActive() bool                { return u.isActive }
func (u *User) Roles() []Role                 { return u.roles }
func (u *User) CreatedAt() time.Time          { return u.createdAt }
func (u *User) UpdatedAt() time.Time          { return u.updatedAt }
func (u *User) Version() int64                { return u.version }

// Event management
func (u *User) Events() []events.DomainEvent  { return u.events }
func (u *User) ClearEvents()                  { u.events = []events.DomainEvent{} }

func (u *User) addEvent(event events.DomainEvent) {
    u.events = append(u.events, event)
}
```

### Aggregate Rules

1. **Single Aggregate Root**: User is the aggregate root
2. **Consistency Boundary**: All changes go through the aggregate root
3. **Transactional**: Changes to aggregate are atomic
4. **Reference by ID**: Other aggregates reference by UserID, not direct reference

## Value Objects

Value objects represent concepts that are defined by their attributes rather than identity.

### Characteristics
- Immutable
- No identity
- Compared by value
- Self-validating

### Email Value Object

```go
package user

import (
    "errors"
    "regexp"
    "strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    email = strings.TrimSpace(strings.ToLower(email))
    
    if email == "" {
        return Email{}, errors.New("email cannot be empty")
    }
    
    if len(email) > 255 {
        return Email{}, errors.New("email too long")
    }
    
    if !emailRegex.MatchString(email) {
        return Email{}, errors.New("invalid email format")
    }
    
    return Email{value: email}, nil
}

func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

// Domain-specific methods
func (e Email) Domain() string {
    parts := strings.Split(e.value, "@")
    if len(parts) != 2 {
        return ""
    }
    return parts[1]
}

func (e Email) IsCompanyEmail(domain string) bool {
    return e.Domain() == domain
}
```

### Money Value Object

```go
package shared

import (
    "errors"
    "fmt"
)

type Money struct {
    amount   int64  // Store as cents to avoid floating point issues
    currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
    if amount < 0 {
        return Money{}, errors.New("amount cannot be negative")
    }
    
    if currency == "" {
        return Money{}, errors.New("currency required")
    }
    
    return Money{
        amount:   int64(amount * 100), // Convert to cents
        currency: currency,
    }, nil
}

func (m Money) Amount() float64 {
    return float64(m.amount) / 100.0
}

func (m Money) Currency() string {
    return m.currency
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    
    return Money{
        amount:   m.amount + other.amount,
        currency: m.currency,
    }, nil
}

func (m Money) Subtract(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot subtract different currencies")
    }
    
    if m.amount < other.amount {
        return Money{}, errors.New("insufficient funds")
    }
    
    return Money{
        amount:   m.amount - other.amount,
        currency: m.currency,
    }, nil
}

func (m Money) Multiply(factor float64) (Money, error) {
    if factor < 0 {
        return Money{}, errors.New("factor cannot be negative")
    }
    
    return Money{
        amount:   int64(float64(m.amount) * factor),
        currency: m.currency,
    }, nil
}

func (m Money) String() string {
    return fmt.Sprintf("%.2f %s", m.Amount(), m.currency)
}

func (m Money) Equals(other Money) bool {
    return m.amount == other.amount && m.currency == other.currency
}
```

## Domain Events

Domain events represent something that happened in the domain that domain experts care about.

### Event Interface

```go
package events

import "time"

type DomainEvent interface {
    EventType() string
    AggregateID() string
    OccurredAt() time.Time
}
```

### User Events

```go
package user

import "time"

// UserCreated event
type UserCreated struct {
    UserID    UserID
    Email     Email
    Name      string
    CreatedAt time.Time
}

func (e UserCreated) EventType() string    { return "user.created" }
func (e UserCreated) AggregateID() string  { return e.UserID.String() }
func (e UserCreated) OccurredAt() time.Time { return e.CreatedAt }

// UserEmailChanged event
type UserEmailChanged struct {
    UserID    UserID
    OldEmail  Email
    NewEmail  Email
    ChangedAt time.Time
}

func (e UserEmailChanged) EventType() string    { return "user.email.changed" }
func (e UserEmailChanged) AggregateID() string  { return e.UserID.String() }
func (e UserEmailChanged) OccurredAt() time.Time { return e.ChangedAt }

// UserActivated event
type UserActivated struct {
    UserID      UserID
    ActivatedAt time.Time
}

func (e UserActivated) EventType() string    { return "user.activated" }
func (e UserActivated) AggregateID() string  { return e.UserID.String() }
func (e UserActivated) OccurredAt() time.Time { return e.ActivatedAt }
```

## Repository Interfaces

Repository interfaces belong to the domain layer, but implementations are in infrastructure.

```go
package repositories

import (
    "context"
    "github.com/tranvuongduy2003/go-mvc/internal/domain/user"
    "github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

type UserRepository interface {
    // Basic CRUD
    Create(ctx context.Context, user *user.User) error
    GetByID(ctx context.Context, id user.UserID) (*user.User, error)
    Update(ctx context.Context, user *user.User) error
    Delete(ctx context.Context, id user.UserID) error
    
    // Queries
    List(ctx context.Context, pagination pagination.Pagination) ([]*user.User, int64, error)
    
    // Domain-specific queries
    GetByEmail(ctx context.Context, email user.Email) (*user.User, error)
    ExistsByEmail(ctx context.Context, email user.Email) (bool, error)
    
    // Complex queries
    FindActiveUsersWithRole(ctx context.Context, roleID string) ([]*user.User, error)
    FindByEmailDomain(ctx context.Context, domain string) ([]*user.User, error)
}
```

## Domain Services

Sometimes business logic doesn't naturally fit in an entity. Use domain services for operations involving multiple entities or external concerns.

```go
package user

type UserDomainService struct {
    userRepo UserRepository
}

func NewUserDomainService(userRepo UserRepository) *UserDomainService {
    return &UserDomainService{userRepo: userRepo}
}

// CheckEmailUniqueness checks if email is unique
func (s *UserDomainService) CheckEmailUniqueness(ctx context.Context, email Email) error {
    exists, err := s.userRepo.ExistsByEmail(ctx, email)
    if err != nil {
        return err
    }
    
    if exists {
        return ErrEmailAlreadyExists
    }
    
    return nil
}

// TransferOwnership transfers ownership from one user to another
func (s *UserDomainService) TransferOwnership(
    ctx context.Context,
    fromUser *User,
    toUser *User,
    resource Resource,
) error {
    if !fromUser.OwnsResource(resource) {
        return ErrNotResourceOwner
    }
    
    if !toUser.IsActive() {
        return ErrTargetUserInactive
    }
    
    // Domain logic for transfer
    fromUser.RemoveResource(resource)
    toUser.AddResource(resource)
    
    return nil
}
```

## Specifications Pattern

Specifications encapsulate business rules and query criteria.

```go
package specifications

import "github.com/tranvuongduy2003/go-mvc/internal/domain/user"

type Specification interface {
    IsSatisfiedBy(user *user.User) bool
}

// ActiveUserSpecification checks if user is active
type ActiveUserSpecification struct{}

func (s ActiveUserSpecification) IsSatisfiedBy(user *user.User) bool {
    return user.IsActive()
}

// UserWithRoleSpecification checks if user has specific role
type UserWithRoleSpecification struct {
    roleName string
}

func (s UserWithRoleSpecification) IsSatisfiedBy(user *user.User) bool {
    for _, role := range user.Roles() {
        if role.Name() == s.roleName {
            return true
        }
    }
    return false
}

// CompositeSpecification allows combining specifications
type AndSpecification struct {
    left  Specification
    right Specification
}

func (s AndSpecification) IsSatisfiedBy(user *user.User) bool {
    return s.left.IsSatisfiedBy(user) && s.right.IsSatisfiedBy(user)
}
```

## Best Practices

### DO ✅

1. **Keep Domain Pure**: No external dependencies (no database, HTTP, etc.)
2. **Rich Behavior**: Entities contain business logic
3. **Immutable Value Objects**: Always create new instances
4. **Validate in Constructors**: Fail fast with clear errors
5. **Use Domain Events**: Communicate changes to other parts of the system
6. **Encapsulate**: Use private fields with public methods
7. **Self-Documenting**: Use ubiquitous language from domain experts

### DON'T ❌

1. **Anemic Domain**: Don't make entities just data bags
2. **Break Encapsulation**: Don't expose internal state via setters
3. **External Dependencies**: Don't reference infrastructure in domain
4. **Mutable Value Objects**: Don't allow value objects to change
5. **Skip Validation**: Don't create invalid domain objects
6. **Public Fields**: Don't expose fields directly
7. **Technical Terms**: Don't use technical jargon in domain models

## Testing Domain Layer

### Unit Testing Entities

```go
func TestUser_ChangeEmail(t *testing.T) {
    // Arrange
    oldEmail, _ := NewEmail("old@example.com")
    newEmail, _ := NewEmail("new@example.com")
    password, _ := NewPassword("Password123!")
    user, _ := NewUser(oldEmail, password, "John Doe")
    
    // Act
    err := user.ChangeEmail(newEmail)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, newEmail, user.Email())
    assert.Len(t, user.Events(), 2) // Created + EmailChanged
    
    event := user.Events()[1].(UserEmailChanged)
    assert.Equal(t, oldEmail, event.OldEmail)
    assert.Equal(t, newEmail, event.NewEmail)
}

func TestUser_ChangeEmail_SameEmail(t *testing.T) {
    // Arrange
    email, _ := NewEmail("test@example.com")
    password, _ := NewPassword("Password123!")
    user, _ := NewUser(email, password, "John Doe")
    
    // Act
    err := user.ChangeEmail(email)
    
    // Assert
    assert.ErrorIs(t, err, ErrEmailNotChanged)
}
```

### Testing Value Objects

```go
func TestEmail_Validation(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"empty email", "", true},
        {"invalid format", "not-an-email", true},
        {"too long", strings.Repeat("a", 250) + "@example.com", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Summary

The Domain Layer is the most important layer in Clean Architecture:

- **Contains business logic**: All core business rules and invariants
- **Independent**: No dependencies on external concerns
- **Rich models**: Entities encapsulate behavior, not just data
- **Value objects**: Immutable, self-validating types
- **Domain events**: Communicate state changes
- **Repository interfaces**: Define data access contracts
- **Specifications**: Encapsulate business rules

**Next**: [Chapter 4: Application Layer](04-application-layer.md) - Learn how to orchestrate domain logic with use cases.

---

**Key Takeaway**: A well-designed domain layer makes your business logic clear, testable, and maintainable. It's the foundation of your application.

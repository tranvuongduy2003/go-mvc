# Clean Architecture & SOLID Principles - Practical Examples

This document provides concrete examples of how the refactoring improved the codebase's adherence to Clean Architecture and SOLID principles.

---

## 1. Dependency Inversion Principle (DIP)

### ❌ BEFORE: High-level depends on low-level
```go
// internal/application/commands/user/upload_avatar_command.go
package commands

import (
    "github.com/tranvuongduy2003/go-mvc/internal/adapters/external" // ❌ Infrastructure layer
)

type UploadAvatarCommandHandler struct {
    fileStorageService *external.FileStorageService // ❌ Concrete type
}
```

**Problems:**
- Application layer knows about MinIO implementation details
- Cannot test without actual MinIO connection
- Cannot switch to S3 without changing application code
- Violates Clean Architecture's Dependency Rule

### ✅ AFTER: Both depend on abstraction
```go
// internal/core/ports/services/file_storage_service.go (Domain/Application Layer)
package services

type FileStorageService interface {
    Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (fileKey string, cdnURL string, err error)
    Delete(ctx context.Context, fileKey string) error
}

// internal/application/commands/user/upload_avatar_command.go
package commands

import (
    "github.com/tranvuongduy2003/go-mvc/internal/core/ports/services" // ✅ Port interface
)

type UploadAvatarCommandHandler struct {
    fileStorageService services.FileStorageService // ✅ Interface type
}

// internal/adapters/external/file_storage_service.go (Infrastructure Layer)
package external

var _ services.FileStorageService = (*FileStorageService)(nil) // ✅ Compile-time check

type FileStorageService struct {
    client *minio.Client // Implementation detail
}

func (s *FileStorageService) Upload(...) (string, string, error) {
    // MinIO-specific implementation
}
```

**Benefits:**
- ✅ Application layer is infrastructure-agnostic
- ✅ Can easily create test doubles
- ✅ Can swap MinIO for S3, Azure, GCS without touching application
- ✅ Dependencies point inward (toward domain)

**Testing Example:**
```go
// Easy to create mock
type mockFileStorage struct {
    UploadFunc func(ctx context.Context, ...) (string, string, error)
}

func (m *mockFileStorage) Upload(ctx context.Context, ...) (string, string, error) {
    return m.UploadFunc(ctx, ...)
}

// Test with mock
mockStorage := &mockFileStorage{
    UploadFunc: func(...) (string, string, error) {
        return "test-key", "https://cdn.test.com/test-key", nil
    },
}

handler := NewUploadAvatarCommandHandler(userRepo, mockStorage, eventBus)
// Test handler without actual storage
```

---

## 2. Interface Segregation Principle (ISP)

### ❌ BEFORE: Fat interface
```go
type AuthService interface {
    // Authentication
    Register(ctx context.Context, req *RegisterRequest) (*AuthenticatedUser, error)
    Login(ctx context.Context, credentials *LoginCredentials) (*AuthenticatedUser, error)
    RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
    
    // Token Management
    Logout(ctx context.Context, userID string) error
    LogoutAll(ctx context.Context, userID string) error
    IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
    BlacklistToken(ctx context.Context, token string) error
    
    // Password Management
    ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
    ResetPassword(ctx context.Context, email string) error
    ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
    
    // Email Verification
    VerifyEmail(ctx context.Context, token string) error
    ResendVerificationEmail(ctx context.Context, email string) error
    
    // ... and more
}
```

**Problems:**
- Handler that only needs login must depend on 17 methods
- Testing requires mocking all methods even if unused
- Single interface has multiple reasons to change
- Violates Single Responsibility Principle

**Example of the problem:**
```go
// LoginHandler only needs Login, but depends on everything
type LoginHandler struct {
    authService services.AuthService // ❌ Depends on 17 methods, uses 1
}

// Mock must implement all 17 methods
type mockAuthService struct {
    mock.Mock
}
func (m *mockAuthService) Register(...) {...}
func (m *mockAuthService) Login(...) {...}
func (m *mockAuthService) RefreshToken(...) {...}
func (m *mockAuthService) Logout(...) {...}
// ... 13 more methods to implement! ❌
```

### ✅ AFTER: Focused interfaces
```go
// Core authentication only
type AuthService interface {
    Register(ctx context.Context, req *RegisterRequest) (*AuthenticatedUser, error)
    Login(ctx context.Context, credentials *LoginCredentials) (*AuthenticatedUser, error)
    RefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
    ValidateToken(ctx context.Context, accessToken string) (*user.User, error)
}

// Token lifecycle management
type TokenManagementService interface {
    Logout(ctx context.Context, userID string) error
    LogoutAll(ctx context.Context, userID string) error
    IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
    BlacklistToken(ctx context.Context, token string) error
}

// Password operations
type PasswordManagementService interface {
    ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
    ResetPassword(ctx context.Context, email string) error
    ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
}
```

**Benefits:**
```go
// LoginHandler depends only on what it needs
type LoginHandler struct {
    authService services.AuthService // ✅ Only 5 methods
}

// LogoutHandler depends only on token management
type LogoutHandler struct {
    tokenService services.TokenManagementService // ✅ Only 4 methods
}

// ChangePasswordHandler depends only on password management
type ChangePasswordHandler struct {
    passwordService services.PasswordManagementService // ✅ Only 3 methods
}

// Easy to mock - only implement what you need
type mockAuthService struct {
    LoginFunc func(ctx context.Context, credentials *LoginCredentials) (*AuthenticatedUser, error)
}

func (m *mockAuthService) Login(ctx context.Context, credentials *LoginCredentials) (*AuthenticatedUser, error) {
    return m.LoginFunc(ctx, credentials)
}

// Only 1 method to mock! ✅
```

---

## 3. Single Responsibility Principle (SRP)

### ❌ BEFORE: Mixed responsibilities
```go
type UserService struct {
    // Write operations
    createUserHandler   *userCommands.CreateUserCommandHandler
    updateUserHandler   *userCommands.UpdateUserCommandHandler
    deleteUserHandler   *userCommands.DeleteUserCommandHandler
    
    // Read operations
    getUserByIDHandler  *userQueries.GetUserByIDQueryHandler
    listUsersHandler    *userQueries.ListUsersQueryHandler
}

// Single service handles both reads and writes
func (s *UserService) CreateUser(...) {...}  // Write
func (s *UserService) GetUserByID(...) {...} // Read
```

**Problems:**
- Service has two reasons to change: read optimization and write logic
- Cannot optimize read and write paths independently
- Violates CQRS principle
- Difficult to scale reads and writes separately

### ✅ AFTER: Separated responsibilities
```go
// Handles only write operations
type UserCommandService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
    UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (UserResponse, error)
    DeleteUser(ctx context.Context, userID string) error
    UploadAvatar(ctx context.Context, userID string, req UploadAvatarRequest) (UserResponse, error)
}

// Handles only read operations
type UserQueryService interface {
    GetUserByID(ctx context.Context, userID string) (UserResponse, error)
    GetUserByEmail(ctx context.Context, email string) (UserResponse, error)
    ListUsers(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error)
}
```

**Benefits:**
```go
// Different implementations can optimize differently
type userCommandService struct {
    repo repositories.UserRepository
    // Write-optimized dependencies
}

type userQueryService struct {
    repo repositories.UserRepository
    cache cache.Service // ✅ Can add caching for queries only
    // Read-optimized dependencies
}

// Can scale independently
// - Write service: Strong consistency, ACID transactions
// - Read service: Eventually consistent, cached, read replicas
```

---

## 4. Open/Closed Principle (OCP)

### ✅ Example: Adding new storage provider

**Without changing existing code:**

```go
// Add S3 implementation
package external

type S3FileStorageService struct {
    client *s3.Client
    bucket string
}

// Implement the same interface
var _ services.FileStorageService = (*S3FileStorageService)(nil)

func (s *S3FileStorageService) Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, string, error) {
    // S3-specific implementation
}

func (s *S3FileStorageService) Delete(ctx context.Context, fileKey string) error {
    // S3-specific implementation
}

// Register in DI
func NewS3FileStorageService(cfg *config.AppConfig) services.FileStorageService {
    return &S3FileStorageService{...}
}
```

**No changes needed in:**
- ✅ Application layer (commands/queries)
- ✅ Domain layer
- ✅ Handlers
- ✅ Business logic

---

## 5. Liskov Substitution Principle (LSP)

### ✅ Example: All storage implementations are substitutable

```go
// Any implementation can be used interchangeably
var storage services.FileStorageService

storage = external.NewFileStorageService(...)    // MinIO
storage = external.NewS3FileStorageService(...)  // S3
storage = external.NewAzureStorageService(...)   // Azure

// Handler works with any implementation
handler := NewUploadAvatarCommandHandler(userRepo, storage, eventBus)
```

---

## 6. Command/Query Validation (CQRS Enhancement)

### ✅ BEFORE: Implicit validation
```go
type CreateUserCommand struct {
    Email    string
    Password string
}

// Validation scattered or missing
```

### ✅ AFTER: Built-in validation
```go
type CreateUserCommand struct {
    Email    string
    Name     string
    Password string
}

// Command implements Command interface
func (c CreateUserCommand) Validate() error {
    if c.Email == "" {
        return errors.New("email is required")
    }
    if len(c.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}

// Can validate before handling
cmd := CreateUserCommand{Email: "test@test.com", Password: "short"}
if err := cmd.Validate(); err != nil {
    // Handle validation error
}
```

---

## 7. Error Handling Consistency

### ❌ BEFORE: Inconsistent errors
```go
func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
    existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
    if err != nil {
        return nil, err // ❌ Raw error
    }
    
    if existingUser != nil {
        return nil, errors.New("user exists") // ❌ Plain error
    }
}
```

### ✅ AFTER: Typed errors
```go
func (h *CreateUserCommandHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
    existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
    if err != nil {
        return nil, apperrors.NewInternalError("Failed to check existing user", err) // ✅ Typed, with context
    }
    
    if existingUser != nil {
        return nil, apperrors.NewConflictError("User already exists", nil) // ✅ Correct HTTP status
    }
    
    newUser, err := user.NewUser(...)
    if err != nil {
        return nil, apperrors.NewValidationError(err.Error(), err) // ✅ Validation error
    }
}
```

**Benefits:**
- ✅ Consistent error types throughout codebase
- ✅ Automatic HTTP status code mapping
- ✅ Proper error context and wrapping
- ✅ Can check error types programmatically

---

## Summary: Architecture Quality Metrics

| Principle | Before | After |
|-----------|--------|-------|
| **Clean Architecture** | ⚠️ Some violations | ✅ Fully compliant |
| **DDD** | ✅ Good | ✅ Excellent |
| **CQRS** | ⚠️ Implicit | ✅ Explicit |
| **DIP** | ❌ Violated | ✅ Compliant |
| **ISP** | ❌ Fat interfaces | ✅ Focused interfaces |
| **SRP** | ⚠️ Mixed | ✅ Separated |
| **OCP** | ⚠️ Hard to extend | ✅ Easy to extend |
| **LSP** | ✅ Good | ✅ Good |
| **Testability** | ⚠️ Difficult | ✅ Easy |
| **Maintainability** | ⚠️ Medium | ✅ High |

---

## Practical Testing Example

### Before: Hard to test
```go
// ❌ Requires actual MinIO connection
handler := NewUploadAvatarCommandHandler(
    userRepo,
    external.NewFileStorageService(realConfig, logger), // ❌ Concrete type
    eventBus,
)
```

### After: Easy to test
```go
// ✅ Use mock
mockStorage := &mockFileStorage{
    UploadFunc: func(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, string, error) {
        return "test-key", "https://test.cdn/test-key", nil
    },
}

handler := NewUploadAvatarCommandHandler(
    mockUserRepo,
    mockStorage,   // ✅ Interface
    mockEventBus,
)

// Test without external dependencies
result, err := handler.Handle(ctx, cmd)
assert.NoError(t, err)
assert.Equal(t, "https://test.cdn/test-key", result.AvatarURL)
```

---

**Conclusion**: The refactoring transformed the codebase from good to excellent, making it fully compliant with all architectural best practices while maintaining the existing strengths of the domain model.

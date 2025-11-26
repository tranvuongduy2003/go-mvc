# RBAC Authentication & Authorization System

This document explains how to use the implemented Role-Based Access Control (RBAC) system for authentication and authorization.

## Overview

The RBAC system provides:
- **Authentication**: JWT-based token validation
- **Authorization**: Role and permission-based access control
- **Middleware**: Ready-to-use Gin middleware for protecting routes
- **Domain Model**: Complete RBAC entities and services

## Components

### 1. Domain Entities
- **Role**: Represents user roles (admin, moderator, user, etc.)
- **Permission**: Represents specific permissions (users:read, posts:write, etc.)
- **UserRole**: Links users to their roles
- **RolePermission**: Links roles to their permissions

### 2. Services
- **RBACService**: Core business logic for role and permission management
- **AuthMiddleware**: JWT authentication middleware (`auth.go`)
- **AuthzMiddleware**: RBAC authorization middleware (`authorization.go`)

### 3. Repositories
- **RoleRepository**: CRUD operations for roles
- **PermissionRepository**: CRUD operations for permissions
- **UserRoleRepository**: User-role assignment operations
- **RolePermissionRepository**: Role-permission assignment operations

## Usage Examples

### 1. Protecting Routes with Authentication

```go
// routes/user_routes.go
func SetupUserRoutes(router *gin.Engine, middlewareManager *middleware.MiddlewareManager, userHandler *handlers.UserHandler) {
    api := router.Group("/api/v1")
    
    // Public routes (no authentication required)
    api.POST("/auth/login", userHandler.Login)
    api.POST("/auth/register", userHandler.Register)
    
    // Protected routes (authentication required)
    authenticated := api.Group("/users")
    authenticated.Use(middlewareManager.AuthRequired())
    {
        authenticated.GET("/profile", userHandler.GetProfile)
        authenticated.PUT("/profile", userHandler.UpdateProfile)
    }
}
```

### 2. Role-Based Route Protection

```go
// routes/admin_routes.go
func SetupAdminRoutes(router *gin.Engine, middlewareManager *middleware.MiddlewareManager, adminHandler *handlers.AdminHandler) {
    api := router.Group("/api/v1/admin")
    
    // Admin only routes
    api.Use(middlewareManager.AuthRequired())
    api.Use(middlewareManager.AdminOnly())
    {
        api.GET("/users", adminHandler.ListUsers)
        api.DELETE("/users/:id", adminHandler.DeleteUser)
        api.POST("/roles", adminHandler.CreateRole)
    }
    
    // Moderator or Admin routes
    moderator := api.Group("/moderate")
    moderator.Use(middlewareManager.ModeratorOrAdmin())
    {
        moderator.PUT("/posts/:id/approve", adminHandler.ApprovePost)
        moderator.DELETE("/comments/:id", adminHandler.DeleteComment)
    }
}
```

### 3. Permission-Based Route Protection

```go
// routes/api_routes.go
func SetupAPIRoutes(router *gin.Engine, middlewareManager *middleware.MiddlewareManager, apiHandler *handlers.APIHandler) {
    api := router.Group("/api/v1")
    api.Use(middlewareManager.AuthRequired())
    
    // Require specific permissions
    users := api.Group("/users")
    {
        users.GET("", middlewareManager.RequirePermission("users", "read"), apiHandler.ListUsers)
        users.POST("", middlewareManager.RequirePermission("users", "create"), apiHandler.CreateUser)
        users.PUT("/:id", middlewareManager.RequirePermission("users", "update"), apiHandler.UpdateUser)
        users.DELETE("/:id", middlewireManager.RequirePermission("users", "delete"), apiHandler.DeleteUser)
    }
    
    posts := api.Group("/posts")
    {
        posts.GET("", middlewareManager.RequirePermission("posts", "read"), apiHandler.ListPosts)
        posts.POST("", middlewareManager.RequirePermission("posts", "create"), apiHandler.CreatePost)
        posts.PUT("/:id", middlewareManager.OwnerOrAdmin("id"), apiHandler.UpdatePost) // Owner or admin can update
    }
}
```

### 4. Multiple Role Requirements

```go
// routes/content_routes.go
func SetupContentRoutes(router *gin.Engine, middlewareManager *middleware.MiddlewareManager, contentHandler *handlers.ContentHandler) {
    api := router.Group("/api/v1/content")
    api.Use(middlewareManager.AuthRequired())
    
    // Require any of these roles
    api.Use(middlewareManager.RequireAnyRole("editor", "moderator", "admin"))
    {
        api.POST("/articles", contentHandler.CreateArticle)
        api.PUT("/articles/:id", contentHandler.UpdateArticle)
        api.DELETE("/articles/:id", contentHandler.DeleteArticle)
    }
}
```

### 5. Optional Authentication

```go
// routes/public_routes.go
func SetupPublicRoutes(router *gin.Engine, middlewareManager *middleware.MiddlewareManager, publicHandler *handlers.PublicHandler) {
    api := router.Group("/api/v1")
    
    // Optional authentication (user info available if logged in)
    api.Use(middlewareManager.OptionalAuth())
    {
        api.GET("/posts", publicHandler.ListPosts) // Public, but may show different data for authenticated users
        api.GET("/posts/:id", publicHandler.GetPost)
    }
}
```

## Available Middleware Methods

### Authentication Middleware
- `AuthRequired()`: Requires valid JWT token
- `OptionalAuth()`: Makes user info available if authenticated, but doesn't require it

### Role-Based Middleware
- `RequireRole(roleName string)`: Requires specific role
- `RequireAnyRole(roleNames ...string)`: Requires any of the specified roles
- `AdminOnly()`: Shortcut for admin role requirement
- `ModeratorOrAdmin()`: Requires moderator or admin role

### Permission-Based Middleware
- `RequirePermission(resource, action string)`: Requires specific permission
- `OwnerOrAdmin(resourceIDParam string)`: Allows resource owner or admin

## Setting Up Default Roles and Permissions

```go
// In your application startup (e.g., main.go or migration)
func setupDefaultRBAC(rbacService rbac.RBACService) error {
    ctx := context.Background()
    
    // Initialize default roles and permissions
    if err := rbacService.InitializeDefaultRolesAndPermissions(ctx); err != nil {
        return fmt.Errorf("failed to setup default RBAC: %w", err)
    }
    
    return nil
}
```

## Accessing User Information in Handlers

```go
// handlers/user_handler.go
func (h *UserHandler) GetProfile(c *gin.Context) {
    // Get authenticated user from context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "User not authenticated")
        return
    }
    
    userClaims, exists := c.Get("user_claims")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "User claims not found")
        return
    }
    
    claims := userClaims.(*jwt.Claims)
    
    // Use userID and claims for business logic
    profile, err := h.userService.GetProfile(c.Request.Context(), userID.(uuid.UUID))
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to get profile")
        return
    }
    
    response.Success(c, profile)
}
```

## JWT Token Format

The JWT tokens contain the following claims:
```go
type Claims struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    Roles     []string  `json:"roles"`
    jwt.RegisteredClaims
}
```

## Error Handling

The middleware will automatically return appropriate HTTP status codes:
- **401 Unauthorized**: Invalid or missing token
- **403 Forbidden**: Valid token but insufficient permissions/roles

## Database Tables

The system uses the following database tables:
- `roles`: Store role definitions
- `permissions`: Store permission definitions  
- `user_roles`: Link users to roles
- `role_permissions`: Link roles to permissions

## Example Role/Permission Structure

```
Roles:
- admin: Full system access
- moderator: Content moderation access
- user: Basic user access

Permissions:
- users:read, users:create, users:update, users:delete
- posts:read, posts:create, posts:update, posts:delete
- comments:read, comments:create, comments:update, comments:delete

Role-Permission Mapping:
- admin: All permissions
- moderator: posts:*, comments:*, users:read
- user: posts:read, comments:read, comments:create (own)
```

This RBAC system provides a flexible and secure way to handle authentication and authorization in your Go application.
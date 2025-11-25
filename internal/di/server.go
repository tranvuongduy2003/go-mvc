package di

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	portServices "github.com/tranvuongduy2003/go-mvc/internal/domain/ports/services"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/interfaces/http/middleware"
	v1 "github.com/tranvuongduy2003/go-mvc/internal/interfaces/http/rest/v1"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// ServerModule provides HTTP server dependencies
var ServerModule = fx.Module("server",
	fx.Provide(
		NewHTTPServer,
		NewGinRouter,
		NewMiddlewareManager,
	),
	// Note: Routes registration moved to main.go after middleware setup
)

// ServerParams holds parameters for server providers
type ServerParams struct {
	fx.In
	Config *config.AppConfig
	Router *gin.Engine
}

// RouterParams holds parameters for router
type RouterParams struct {
	fx.In
	Config *config.AppConfig
	Logger *logger.Logger
}

// RouteParams holds parameters for route registration
type RouteParams struct {
	fx.In
	Router      *gin.Engine
	UserHandler *v1.UserHandler
	UserService *services.UserService
	AuthHandler *v1.AuthHandler
	AuthService portServices.AuthService
}

// MiddlewareParams holds parameters for middleware setup
type MiddlewareParams struct {
	fx.In
	Router            *gin.Engine
	MiddlewareManager *middleware.MiddlewareManager
	Config            *config.AppConfig
}

// NewGinRouter creates and configures Gin router
func NewGinRouter(params RouterParams) *gin.Engine {
	if params.Config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router without default middleware
	router := gin.New()

	return router
}

// NewMiddlewareManager creates middleware manager
func NewMiddlewareManager(
	config *config.AppConfig,
	logger *logger.Logger,
	jwtService jwt.JWTService,
) *middleware.MiddlewareManager {
	middlewareConfig := middleware.DefaultMiddlewareConfig()

	// Customize config based on environment
	if config.App.Environment == "production" {
		middlewareConfig.Logging.LogRequestBody = false
		middlewareConfig.Logging.LogResponseBody = false
		middlewareConfig.RateLimit.RPS = 100
		middlewareConfig.RateLimit.Burst = 200
	} else {
		middlewareConfig.Logging.LogRequestBody = true
		middlewareConfig.Logging.LogResponseBody = true
		middlewareConfig.RateLimit.RPS = 1000
		middlewareConfig.RateLimit.Burst = 2000
	}

	return middleware.NewMiddlewareManager(logger, middlewareConfig, jwtService)
}

// SetupMiddleware configures all middleware
func SetupMiddleware(params MiddlewareParams) {
	if params.Config.App.Environment == "production" {
		allowedOrigins := []string{
			"https://yourdomain.com",
			"https://www.yourdomain.com",
		}
		params.MiddlewareManager.SetupProductionMiddleware(params.Router, allowedOrigins)
	} else {
		params.MiddlewareManager.SetupDevelopmentMiddleware(params.Router)
	}
}

// RegisterRoutes registers all application routes
func RegisterRoutes(params RouteParams) {
	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(params.AuthService)

	v1API := params.Router.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1API.Group("/auth")
		{
			auth.POST("/register", params.AuthHandler.Register)
			auth.POST("/login", params.AuthHandler.Login)
			auth.POST("/refresh", params.AuthHandler.RefreshToken)
			auth.POST("/verify-email", params.AuthHandler.VerifyEmail)
			auth.POST("/reset-password", params.AuthHandler.ResetPassword)
			auth.POST("/confirm-reset", params.AuthHandler.ConfirmPasswordReset)
			auth.POST("/resend-verification", params.AuthHandler.ResendVerificationEmail)
		}

		// Protected auth routes (authentication required)
		protectedAuth := v1API.Group("/auth")
		protectedAuth.Use(authMiddleware.RequireAuth())
		{
			protectedAuth.POST("/logout", params.AuthHandler.Logout)
			protectedAuth.POST("/logout-all", params.AuthHandler.LogoutAllDevices)
			protectedAuth.GET("/profile", params.AuthHandler.GetProfile)
			protectedAuth.GET("/permissions", params.AuthHandler.GetPermissions)
			protectedAuth.PUT("/change-password", params.AuthHandler.ChangePassword)
		}

		// User routes (protected)
		users := v1API.Group("/users")
		{
			users.POST("", params.UserHandler.CreateUser)
			users.GET("", params.UserHandler.ListUsers)
			users.GET("/:id", params.UserHandler.GetUserByID)
			users.PUT("/:id", params.UserHandler.UpdateUser)
			users.DELETE("/:id", params.UserHandler.DeleteUser)
			users.POST("/:id/avatar", params.UserHandler.UploadAvatar) // Avatar upload endpoint
		}

		// Test route
		v1API.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Test API endpoint",
				"data":    "Hello from Go MVC!",
			})
		})
	}
}

// NewHTTPServer creates HTTP server
func NewHTTPServer(params ServerParams) *http.Server {
	addr := fmt.Sprintf(":%d", params.Config.Server.HTTP.Port)
	return &http.Server{
		Addr:    addr,
		Handler: params.Router,
	}
}

// HTTPServerLifecycle handles HTTP server lifecycle
func HTTPServerLifecycle(
	lc fx.Lifecycle,
	server *http.Server,
	config *config.AppConfig,
	logger *logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server",
				zap.String("addr", server.Addr),
				zap.String("environment", config.App.Environment),
			)

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}

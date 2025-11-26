package presentation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	appservices "github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	v1 "github.com/tranvuongduy2003/go-mvc/internal/presentation/http/handlers/v1"
	"github.com/tranvuongduy2003/go-mvc/internal/presentation/http/middleware"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

var ServerModule = fx.Module("server",
	fx.Provide(
		NewHTTPServer,
		NewGinRouter,
		NewMiddlewareManager,
	),
)

type ServerParams struct {
	fx.In
	Config *config.AppConfig
	Router *gin.Engine
}

type RouterParams struct {
	fx.In
	Config *config.AppConfig
	Logger *logger.Logger
}

type RouteParams struct {
	fx.In
	Router      *gin.Engine
	UserHandler *v1.UserHandler
	UserService *appservices.UserService
	AuthHandler *v1.AuthHandler
	AuthService appservices.AuthService
}

type MiddlewareParams struct {
	fx.In
	Router            *gin.Engine
	MiddlewareManager *middleware.MiddlewareManager
	Config            *config.AppConfig
}

func NewGinRouter(params RouterParams) *gin.Engine {
	if params.Config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	return router
}

func NewMiddlewareManager(
	config *config.AppConfig,
	logger *logger.Logger,
	jwtService jwt.JWTService,
) *middleware.MiddlewareManager {
	middlewareConfig := middleware.DefaultMiddlewareConfig()

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

func RegisterRoutes(params RouteParams) {
	authMiddleware := middleware.NewAuthMiddleware(&params.AuthService)

	v1API := params.Router.Group("/api/v1")
	{
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

		protectedAuth := v1API.Group("/auth")
		protectedAuth.Use(authMiddleware.RequireAuth())
		{
			protectedAuth.POST("/logout", params.AuthHandler.Logout)
			protectedAuth.POST("/logout-all", params.AuthHandler.LogoutAllDevices)
			protectedAuth.GET("/profile", params.AuthHandler.GetProfile)
			protectedAuth.GET("/permissions", params.AuthHandler.GetPermissions)
			protectedAuth.PUT("/change-password", params.AuthHandler.ChangePassword)
		}

		users := v1API.Group("/users")
		{
			users.POST("", params.UserHandler.CreateUser)
			users.GET("", params.UserHandler.ListUsers)
			users.GET("/:id", params.UserHandler.GetUserByID)
			users.PUT("/:id", params.UserHandler.UpdateUser)
			users.DELETE("/:id", params.UserHandler.DeleteUser)
			users.POST("/:id/avatar", params.UserHandler.UploadAvatar) // Avatar upload endpoint
		}

		v1API.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Test API endpoint",
				"data":    "Hello from Go MVC!",
			})
		})
	}
}

func NewHTTPServer(params ServerParams) *http.Server {
	addr := fmt.Sprintf(":%d", params.Config.Server.HTTP.Port)
	return &http.Server{
		Addr:    addr,
		Handler: params.Router,
	}
}

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

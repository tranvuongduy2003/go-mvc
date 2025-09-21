package fxmodules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/handlers"
	"github.com/tranvuongduy2003/go-mvc/internal/handlers/http/middleware"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// ServerModule provides HTTP server dependencies
var ServerModule = fx.Module("server",
	fx.Provide(
		NewHTTPServer,
		NewGinRouter,
		NewMiddlewareManager,
	),
	fx.Invoke(RegisterRoutes),
	fx.Invoke(SetupMiddleware),
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
	Logger *zap.Logger
}

// RouteParams holds parameters for route registration
type RouteParams struct {
	fx.In
	Router      *gin.Engine
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
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
	logger *zap.Logger,
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
	v1 := params.Router.Group("/api/v1")
	{
		// Test route
		v1.GET("/test", func(c *gin.Context) {
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
	zapLogger *zap.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			zapLogger.Info("Starting HTTP server",
				zap.String("addr", server.Addr),
				zap.String("environment", config.App.Environment),
			)

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					zapLogger.Fatal("Failed to start HTTP server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			zapLogger.Info("Shutting down HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}

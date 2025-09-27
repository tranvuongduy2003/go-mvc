package di

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/di/modules"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
	"github.com/tranvuongduy2003/go-mvc/pkg/validator"
)

// ApplicationModule provides application layer dependencies
var ApplicationModule = fx.Module("application",
	// Include domain modules
	modules.UserModule,
	modules.AuthModule,
	modules.JobsModule,

	// Core application services
	fx.Provide(
		NewJWTService,
	),

	// Initialize validator first before other services
	fx.Invoke(InitializeValidator),
)

// ApplicationParams holds parameters for application service providers
type ApplicationParams struct {
	fx.In
	JWTService jwt.JWTService
	Logger     *logger.Logger
	Tracing    *tracing.TracingService
}

// JWTParams holds parameters for JWT service
type JWTParams struct {
	fx.In
	Config *config.AppConfig
}

// ValidatorParams holds parameters for validator
type ValidatorParams struct {
	fx.In
	Logger *logger.Logger
}

// NewJWTService provides JWT service
func NewJWTService(params JWTParams) jwt.JWTService {
	return jwt.NewService(params.Config.JWT)
}

// InitializeValidator initializes the validator
func InitializeValidator(logger *logger.Logger) error {
	if err := validator.InitValidator(); err != nil {
		logger.Error("Failed to initialize validator: " + err.Error())
		return err
	}
	logger.Info("Validator initialized successfully")
	return nil
}

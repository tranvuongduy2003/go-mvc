package fxmodules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/application/validators"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// ApplicationModule provides application layer dependencies
var ApplicationModule = fx.Module("application",
	fx.Provide(
		NewJWTService,
		NewUserApplicationService,
		NewUserValidator,
	),
)

// ApplicationParams holds parameters for application service providers
type ApplicationParams struct {
	fx.In
	UserService *user.Service
	Repository  user.Repository
	JWTService  *jwt.Service
	Logger      *logger.Logger
	Tracing     *tracing.TracingService
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
func NewJWTService(params JWTParams) *jwt.Service {
	return jwt.NewService(params.Config.JWT)
}

// NewUserApplicationService provides user application service
func NewUserApplicationService(params ApplicationParams) *services.UserApplicationService {
	return services.NewUserApplicationService(
		params.UserService,
		params.Repository,
		params.JWTService,
		params.Logger,
		params.Tracing,
	)
}

// NewUserValidator provides user validator
func NewUserValidator(params ValidatorParams) *validators.UserValidator {
	return validators.NewUserValidator(params.Logger)
}

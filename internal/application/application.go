package application

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/tracing"
	"github.com/tranvuongduy2003/go-mvc/internal/modules"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
	"github.com/tranvuongduy2003/go-mvc/pkg/validator"
)

var ApplicationModule = fx.Module("application",
	modules.UserModule,
	modules.AuthModule,
	modules.JobModule,
	modules.MessagingModule,

	fx.Provide(
		NewJWTService,
	),

	fx.Invoke(InitializeValidator),
)

type ApplicationParams struct {
	fx.In
	JWTService jwt.JWTService
	Logger     *logger.Logger
	Tracing    *tracing.TracingService
}

type JWTParams struct {
	fx.In
	Config *config.AppConfig
}

type ValidatorParams struct {
	fx.In
	Logger *logger.Logger
}

func NewJWTService(params JWTParams) jwt.JWTService {
	return jwt.NewService(params.Config.JWT)
}

func InitializeValidator(logger *logger.Logger) error {
	if err := validator.InitValidator(); err != nil {
		logger.Error("Failed to initialize validator: " + err.Error())
		return err
	}
	logger.Info("Validator initialized successfully")
	return nil
}

package fxmodules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/security"
)

// DomainModule provides domain layer dependencies
var DomainModule = fx.Module("domain",
	fx.Provide(),
)

// DomainParams holds parameters for domain service providers
type DomainParams struct {
	fx.In
	PasswordHasher *security.PasswordHasher
	Logger         *logger.Logger
}

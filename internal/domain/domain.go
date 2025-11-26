package domain

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/security"
)

var DomainModule = fx.Module("domain",
	fx.Provide(),
)

type DomainParams struct {
	fx.In
	PasswordHasher *security.PasswordHasher
	Logger         *logger.Logger
}

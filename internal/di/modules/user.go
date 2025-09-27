package modules

import (
	"go.uber.org/fx"

	"github.com/tranvuongduy2003/go-mvc/internal/adapters/external"
	userCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/user"
	"github.com/tranvuongduy2003/go-mvc/internal/application/events/handlers"
	userQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/user"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	userValidators "github.com/tranvuongduy2003/go-mvc/internal/application/validators/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// UserModule provides user domain dependencies
var UserModule = fx.Module("user",
	fx.Provide(
		NewCreateUserCommandHandler,
		NewUpdateUserCommandHandler,
		NewDeleteUserCommandHandler,
		NewUploadAvatarCommandHandler,
		NewGetUserByIDQueryHandler,
		NewListUsersQueryHandler,
		NewUserService,
		NewUserValidator,
		NewUserEventHandler,
	),
	fx.Invoke(SetupUserEventSubscriptions),
)

// NewCreateUserCommandHandler provides CreateUserCommandHandler
func NewCreateUserCommandHandler(userRepo repositories.UserRepository) *userCommands.CreateUserCommandHandler {
	return userCommands.NewCreateUserCommandHandler(userRepo)
}

// NewUpdateUserCommandHandler provides UpdateUserCommandHandler
func NewUpdateUserCommandHandler(userRepo repositories.UserRepository) *userCommands.UpdateUserCommandHandler {
	return userCommands.NewUpdateUserCommandHandler(userRepo)
}

// NewDeleteUserCommandHandler provides DeleteUserCommandHandler
func NewDeleteUserCommandHandler(userRepo repositories.UserRepository) *userCommands.DeleteUserCommandHandler {
	return userCommands.NewDeleteUserCommandHandler(userRepo)
}

// NewGetUserByIDQueryHandler provides GetUserByIDQueryHandler
func NewGetUserByIDQueryHandler(userRepo repositories.UserRepository) *userQueries.GetUserByIDQueryHandler {
	return userQueries.NewGetUserByIDQueryHandler(userRepo)
}

// NewListUsersQueryHandler provides ListUsersQueryHandler
func NewListUsersQueryHandler(userRepo repositories.UserRepository) *userQueries.ListUsersQueryHandler {
	return userQueries.NewListUsersQueryHandler(userRepo)
}

// NewUploadAvatarCommandHandler provides UploadAvatarCommandHandler
func NewUploadAvatarCommandHandler(
	userRepo repositories.UserRepository,
	fileStorageService *external.FileStorageService,
	eventBus messaging.EventBus,
) *userCommands.UploadAvatarCommandHandler {
	return userCommands.NewUploadAvatarCommandHandler(userRepo, fileStorageService, eventBus)
}

// UserServiceParams holds parameters for UserService
type UserServiceParams struct {
	fx.In
	CreateUserHandler   *userCommands.CreateUserCommandHandler
	UpdateUserHandler   *userCommands.UpdateUserCommandHandler
	DeleteUserHandler   *userCommands.DeleteUserCommandHandler
	UploadAvatarHandler *userCommands.UploadAvatarCommandHandler
	GetUserByIDHandler  *userQueries.GetUserByIDQueryHandler
	ListUsersHandler    *userQueries.ListUsersQueryHandler
}

// NewUserService provides UserService
func NewUserService(params UserServiceParams) *services.UserService {
	return services.NewUserService(
		params.CreateUserHandler,
		params.UpdateUserHandler,
		params.DeleteUserHandler,
		params.UploadAvatarHandler,
		params.GetUserByIDHandler,
		params.ListUsersHandler,
	)
}

// NewUserValidator provides UserValidator
func NewUserValidator() userValidators.IUserValidator {
	return userValidators.NewUserValidator()
}

// NewUserEventHandler provides UserEventHandler
func NewUserEventHandler(logger *logger.Logger) *handlers.UserEventHandler {
	return handlers.NewUserEventHandler(logger.Logger)
}

// SetupUserEventSubscriptions sets up event subscriptions for user events
func SetupUserEventSubscriptions(eventHandler *handlers.UserEventHandler, eventBus messaging.EventBus) error {
	return eventHandler.SetupEventSubscriptions(eventBus)
}

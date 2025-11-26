package modules

import (
	"go.uber.org/fx"

	userCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/user"
	eventHandlers "github.com/tranvuongduy2003/go-mvc/internal/application/event_handlers"
	userQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/user"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	userValidators "github.com/tranvuongduy2003/go-mvc/internal/application/validators/user"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/external"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
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
func NewCreateUserCommandHandler(userRepo user.UserRepository) *userCommands.CreateUserCommandHandler {
	return userCommands.NewCreateUserCommandHandler(userRepo)
}

// NewUpdateUserCommandHandler provides UpdateUserCommandHandler
func NewUpdateUserCommandHandler(userRepo user.UserRepository) *userCommands.UpdateUserCommandHandler {
	return userCommands.NewUpdateUserCommandHandler(userRepo)
}

// NewDeleteUserCommandHandler provides DeleteUserCommandHandler
func NewDeleteUserCommandHandler(userRepo user.UserRepository) *userCommands.DeleteUserCommandHandler {
	return userCommands.NewDeleteUserCommandHandler(userRepo)
}

// NewGetUserByIDQueryHandler provides GetUserByIDQueryHandler
func NewGetUserByIDQueryHandler(userRepo user.UserRepository) *userQueries.GetUserByIDQueryHandler {
	return userQueries.NewGetUserByIDQueryHandler(userRepo)
}

// NewListUsersQueryHandler provides ListUsersQueryHandler
func NewListUsersQueryHandler(userRepo user.UserRepository) *userQueries.ListUsersQueryHandler {
	return userQueries.NewListUsersQueryHandler(userRepo)
}

// NewUploadAvatarCommandHandler provides UploadAvatarCommandHandler
func NewUploadAvatarCommandHandler(
	userRepo user.UserRepository,
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
func NewUserEventHandler(logger *logger.Logger) *eventHandlers.UserEventHandler {
	return eventHandlers.NewUserEventHandler(logger.Logger)
}

// SetupUserEventSubscriptions sets up event subscriptions for user events
func SetupUserEventSubscriptions(eventHandler *eventHandlers.UserEventHandler, eventBus messaging.EventBus) error {
	return eventHandler.SetupEventSubscriptions(eventBus)
}

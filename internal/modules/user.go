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

func NewCreateUserCommandHandler(userRepo user.UserRepository) *userCommands.CreateUserCommandHandler {
	return userCommands.NewCreateUserCommandHandler(userRepo)
}

func NewUpdateUserCommandHandler(userRepo user.UserRepository) *userCommands.UpdateUserCommandHandler {
	return userCommands.NewUpdateUserCommandHandler(userRepo)
}

func NewDeleteUserCommandHandler(userRepo user.UserRepository) *userCommands.DeleteUserCommandHandler {
	return userCommands.NewDeleteUserCommandHandler(userRepo)
}

func NewGetUserByIDQueryHandler(userRepo user.UserRepository) *userQueries.GetUserByIDQueryHandler {
	return userQueries.NewGetUserByIDQueryHandler(userRepo)
}

func NewListUsersQueryHandler(userRepo user.UserRepository) *userQueries.ListUsersQueryHandler {
	return userQueries.NewListUsersQueryHandler(userRepo)
}

func NewUploadAvatarCommandHandler(
	userRepo user.UserRepository,
	fileStorageService *external.FileStorageService,
	eventBus messaging.EventBus,
) *userCommands.UploadAvatarCommandHandler {
	return userCommands.NewUploadAvatarCommandHandler(userRepo, fileStorageService, eventBus)
}

type UserServiceParams struct {
	fx.In
	CreateUserHandler   *userCommands.CreateUserCommandHandler
	UpdateUserHandler   *userCommands.UpdateUserCommandHandler
	DeleteUserHandler   *userCommands.DeleteUserCommandHandler
	UploadAvatarHandler *userCommands.UploadAvatarCommandHandler
	GetUserByIDHandler  *userQueries.GetUserByIDQueryHandler
	ListUsersHandler    *userQueries.ListUsersQueryHandler
}

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

func NewUserValidator() userValidators.IUserValidator {
	return userValidators.NewUserValidator()
}

func NewUserEventHandler(logger *logger.Logger) *eventHandlers.UserEventHandler {
	return eventHandlers.NewUserEventHandler(logger.Logger)
}

func SetupUserEventSubscriptions(eventHandler *eventHandlers.UserEventHandler, eventBus messaging.EventBus) error {
	return eventHandler.SetupEventSubscriptions(eventBus)
}

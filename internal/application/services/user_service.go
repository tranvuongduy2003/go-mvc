package services

import (
	"context"
	"mime/multipart"

	userCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/user"
	userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"
	userQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/user"
)

type UserService struct {
	createUserHandler   *userCommands.CreateUserCommandHandler
	updateUserHandler   *userCommands.UpdateUserCommandHandler
	deleteUserHandler   *userCommands.DeleteUserCommandHandler
	uploadAvatarHandler *userCommands.UploadAvatarCommandHandler
	getUserByIDHandler  *userQueries.GetUserByIDQueryHandler
	listUsersHandler    *userQueries.ListUsersQueryHandler
}

func NewUserService(
	createUserHandler *userCommands.CreateUserCommandHandler,
	updateUserHandler *userCommands.UpdateUserCommandHandler,
	deleteUserHandler *userCommands.DeleteUserCommandHandler,
	uploadAvatarHandler *userCommands.UploadAvatarCommandHandler,
	getUserByIDHandler *userQueries.GetUserByIDQueryHandler,
	listUsersHandler *userQueries.ListUsersQueryHandler,
) *UserService {
	return &UserService{
		createUserHandler:   createUserHandler,
		updateUserHandler:   updateUserHandler,
		deleteUserHandler:   deleteUserHandler,
		uploadAvatarHandler: uploadAvatarHandler,
		getUserByIDHandler:  getUserByIDHandler,
		listUsersHandler:    listUsersHandler,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req userDto.CreateUserRequest) (userDto.UserResponse, error) {
	cmd := userCommands.CreateUserCommand{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}

	user, err := s.createUserHandler.Handle(ctx, cmd)
	if err != nil {
		return userDto.UserResponse{}, err
	}

	return userDto.UserResponseFromDomain(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (userDto.UserResponse, error) {
	query := userQueries.GetUserByIDQuery{
		ID: id,
	}

	user, err := s.getUserByIDHandler.Handle(ctx, query)
	if err != nil {
		return userDto.UserResponse{}, err
	}

	return userDto.UserResponseFromDomain(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req userDto.UpdateUserRequest) (userDto.UserResponse, error) {
	cmd := userCommands.UpdateUserCommand{
		ID:    id,
		Name:  req.Name,
		Phone: req.Phone,
	}

	user, err := s.updateUserHandler.Handle(ctx, cmd)
	if err != nil {
		return userDto.UserResponse{}, err
	}

	return userDto.UserResponseFromDomain(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	cmd := userCommands.DeleteUserCommand{
		ID: id,
	}

	return s.deleteUserHandler.Handle(ctx, cmd)
}

func (s *UserService) ListUsers(ctx context.Context, req userDto.ListUsersRequest) (userDto.ListUsersResponse, error) {
	query := userQueries.ListUsersQuery{
		Page:     req.Page,
		Limit:    req.Limit,
		Search:   req.Search,
		SortBy:   req.SortBy,
		SortDir:  req.SortDir,
		IsActive: req.IsActive,
	}

	users, pag, err := s.listUsersHandler.Handle(ctx, query)
	if err != nil {
		return userDto.ListUsersResponse{}, err
	}

	return userDto.ListUsersResponse{
		Users: userDto.UserResponseListFromDomain(users),
		Pagination: userDto.PaginationDTO{
			Page:     pag.Page,
			PageSize: pag.PageSize,
			Total:    pag.Total,
			Pages:    pag.Pages,
		},
	}, nil
}

func (s *UserService) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (userDto.UserResponse, error) {
	cmd := userCommands.UploadAvatarCommand{
		UserID: userID,
		File:   file,
		Header: header,
	}

	return s.uploadAvatarHandler.Handle(ctx, cmd)
}

package user

import userDto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/user"

// IUserValidator defines the interface for user validation operations
type IUserValidator interface {
	ValidateCreateUserRequest(req userDto.CreateUserRequest) map[string]string
	ValidateUpdateUserRequest(req userDto.UpdateUserRequest) map[string]string
	ValidateListUsersRequest(req userDto.ListUsersRequest) map[string]string
}

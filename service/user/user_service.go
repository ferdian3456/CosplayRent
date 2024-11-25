package user

import (
	"context"
	"cosplayrent/model/web/user"
)

type UserService interface {
	Create(ctx context.Context, request user.UserCreateRequest) string
	Login(ctx context.Context, request user.UserLoginRequest) (string, error)
	FindByUUID(ctx context.Context, uuid string) user.UserResponse
	FindAll(ctx context.Context, uuid string) []user.UserResponse
	Update(ctx context.Context, request user.UserUpdateRequest, uuid string)
	Delete(ctx context.Context, uuid string)
	VerifyAndRetrieve(ctx context.Context, token string) (user.UserResponse, error)
}

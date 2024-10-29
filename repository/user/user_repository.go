package user

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/user"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user domain.User)
	Login(ctx context.Context, tx *sql.Tx, name string) (domain.User, error)
	FindByUUID(ctx context.Context, tx *sql.Tx, uuid string) (user.UserResponse, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]user.UserResponse, error)
	Update(ctx context.Context, tx *sql.Tx, user user.UserUpdateRequest)
	Delete(ctx context.Context, tx *sql.Tx, uuid string)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (user.UserResponse, error)
}

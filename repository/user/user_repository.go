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
	FindAll(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserResponse, error)
	Update(ctx context.Context, tx *sql.Tx, user user.UserUpdateRequest, uuid string)
	Delete(ctx context.Context, tx *sql.Tx, uuid string)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (user.UserResponse, error)
	AddOrUpdateIdentityCard(ctx context.Context, tx *sql.Tx, uuid string, IdentityCardImage string)
	GetIdentityCard(ctx context.Context, tx *sql.Tx, uuid string) (string, error)
	GetEMoneyAmount(ctx context.Context, tx *sql.Tx, uuid string) (float64, error)
	TopUp(ctx context.Context, tx *sql.Tx, emoney user.TopUpEmoney, uuid string)
}

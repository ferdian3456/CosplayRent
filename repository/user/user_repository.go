package user

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/user"
	"database/sql"
	"time"
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
	GetEMoneyAmount(ctx context.Context, tx *sql.Tx, uuid string) (user.UserEmoneyResponse, error)
	TopUp(ctx context.Context, tx *sql.Tx, emoney float64, uuid string, time *time.Time)
	AfterBuy(ctx context.Context, tx *sql.Tx, orderamount float64, buyeruuid string, selleruuid string, timeNow *time.Time)
	//CreateTopUpOrder(ctx context.Context, tx *sql.Tx, orderid string, uuid string, emoney user.TopUpEmoney, time *time.Time)
}

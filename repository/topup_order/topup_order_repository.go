package topup_order

import (
	"context"
	"cosplayrent/model/web/user"
	"database/sql"
	"time"
)

type TopUpOrderRepository interface {
	CreateTopUpOrder(ctx context.Context, tx *sql.Tx, orderid string, uuid string, emoney user.TopUpEmoney, time *time.Time)
	FindUserIdByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (userId string, err error)
	UpdateTopUpOrder(ctx context.Context, tx *sql.Tx, topuporderid string, time *time.Time)
	FindTopUpOrderHistoryByUserId(ctx context.Context, tx *sql.Tx, userid string) ([]user.UserEmoneyResponse, error)
}

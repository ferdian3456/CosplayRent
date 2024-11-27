package order

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/order"
	"database/sql"
)

type OrderRepository interface {
	Create(ctx context.Context, tx *sql.Tx, costume domain.Order)
	FindByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.OrderResponse, error)
	DirectlyOrderToMidtrans(ctx context.Context, tx *sql.Tx, uuid string, sendOrderToDatabase order.DirectlyOrderToMidtrans)
}

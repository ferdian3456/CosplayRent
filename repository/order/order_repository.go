package order

import (
	"context"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/order"
	"cosplayrent/model/web/user"
	"database/sql"
)

type OrderRepository interface {
	Create(ctx context.Context, tx *sql.Tx, costume domain.Order)
	FindByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.OrderResponse, error)
	DirectlyOrderToMidtrans(ctx context.Context, tx *sql.Tx, uuid string, sendOrderToDatabase order.DirectlyOrderToMidtrans)
	FindBuyerIdByOrderId(ctx context.Context, tx *sql.Tx, uuid string) (string, error)
	FindSellerIdByOrderId(ctx context.Context, tx *sql.Tx, uuid string) (string, error)
	FindOrderDetailByOrderId(ctx context.Context, tx *sql.Tx, orderid string) (order.OrderResponse, error)
	FindOrderHistoryByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserEmoneyResponse, error)
	FindOrderHistoryBySellerId(ctx context.Context, tx *sql.Tx, uuid string) ([]user.UserEmoneyResponse, error)
	FindOrderBySellerId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.AllSellerOrderResponse, error)
	UpdateSellerOrder(ctx context.Context, tx *sql.Tx, updateRequest order.OrderUpdateRequest, sellerid string, orderid string)
	FindOrderByUserId(ctx context.Context, tx *sql.Tx, uuid string) ([]order.AllUserOrderResponse, error)
}

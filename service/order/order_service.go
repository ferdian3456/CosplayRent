package order

import (
	"context"
	"cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/order"
)

type OrderService interface {
	Create(ctx context.Context, request order.OrderCreateRequest)
	FindByUserId(ctx context.Context, uuid string) []order.OrderResponse
	DirectlyOrderToMidtrans(ctx context.Context, uuid string, directOrderToMidtrans order.DirectlyOrderToMidtrans) midtrans.MidtransResponse
}

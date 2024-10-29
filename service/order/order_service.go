package order

import (
	"context"
	"cosplayrent/model/web/order"
)

type OrderService interface {
	Create(ctx context.Context, request order.OrderCreateRequest)
	FindByUserId(ctx context.Context, uuid string) []order.OrderResponse
}

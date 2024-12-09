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
	FindOrderDetailByOrderId(ctx context.Context, orderid string) order.OrderResponse
	GetAllSellerOrder(ctx context.Context, sellerid string) []order.AllSellerOrderResponse
	UpdateSellerOrder(ctx context.Context, updateRequest order.OrderUpdateRequest, sellerid string, orderid string)
	GetDetailOrderByOrderId(ctx context.Context, sellerid string, orderid string) order.OrderDetailByOrderIdResponse
	GetUserDetailOrder(ctx context.Context, userid string, orderid string) order.GetUserOrderDetailResponse
	GetAllUserOrder(ctx context.Context, userid string) []order.AllUserOrderResponse
	CheckBalanceWithOrderAmount(ctx context.Context, checkbalance order.CheckBalanceWithOrderAmount, uuid string) order.CheckBalanceWithOrderAmountReponse
}

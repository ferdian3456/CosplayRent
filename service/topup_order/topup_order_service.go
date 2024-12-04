package topup_order

import (
	"context"
	midtranss "cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/topup_order"
	"cosplayrent/model/web/user"
)

type TopupOrderService interface {
	CreateTopUpOrder(ctx context.Context, topUpEMoneyRequest user.TopUpEmoney, uuid string) midtranss.MidtransResponse
	CheckTopUpOrderByOrderId(ctx context.Context, orderID string) topup_order.TopupOrderResponse
}

package topup_order

import (
	"context"
	midtranss "cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/user"
)

type TopupOrderService interface {
	CreateTopUpOrder(ctx context.Context, topUpEMoneyRequest user.TopUpEmoney, uuid string) midtranss.MidtransResponse
	FindUserIdByOrderId(ctx context.Context)
}

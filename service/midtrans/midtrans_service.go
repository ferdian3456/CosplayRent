package midtrans

import (
	"context"
	"cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/order"
)

type MidtransService interface {
	CreateTransaction(ctx context.Context, request order.DirectlyOrderToMidtrans, uuid string) midtrans.MidtransResponse
	MidtransCallBack(ctx context.Context, orderid string)
}

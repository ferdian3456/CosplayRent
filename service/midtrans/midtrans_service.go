package midtrans

import (
	"context"
	"cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/order"
)

type MidtransService interface {
	CreateTransaction(ctx context.Context, request order.DirectlyOrderToMidtrans, uuid string) midtrans.MidtransResponse
	MidtransCallBack(ctx context.Context, orderid string, orderamount string)
	CreateOrderTopUp(ctx context.Context, orderid string, username string, useremail string, uuid string, emoneyamount float64) midtrans.MidtransResponse
}

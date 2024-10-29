package midtrans

import (
	"context"
	"cosplayrent/model/web/midtrans"
)

type MidtransService interface {
	CreateTransaction(ctx context.Context, request midtrans.MidtransRequest, orderid string) midtrans.MidtransResponse
	MidtransCallBack(ctx context.Context, orderid string)
}

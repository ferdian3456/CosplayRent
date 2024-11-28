package midtrans

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	midtransWeb "cosplayrent/model/web/midtrans"
	orders "cosplayrent/model/web/order"
	midtranss "cosplayrent/repository/midtrans"
	"cosplayrent/repository/order"
	"cosplayrent/repository/user"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"log"
	"os"
	"strconv"
	"time"
)

type MidtransServiceImpl struct {
	UserRepository     user.UserRepository
	OrderRepository    order.OrderRepository
	MidtransRepository midtranss.MidtransRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewMidtransService(midtransRepository midtranss.MidtransRepository, userRepository user.UserRepository, orderRepository order.OrderRepository, DB *sql.DB, validate *validator.Validate) MidtransService {
	return &MidtransServiceImpl{
		OrderRepository:    orderRepository,
		MidtransRepository: midtransRepository,
		UserRepository:     userRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *MidtransServiceImpl) CreateTransaction(ctx context.Context, request orders.DirectlyOrderToMidtrans, userUUID string) midtransWeb.MidtransResponse {
	log.Printf("User with uuid: %s enter Midtrans Service: CreateTransaction", userUUID)
	//fmt.Println("test")
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	server_key := os.Getenv("MIDTRANS_SERVER_KEY")

	//fmt.Println("request.Id =", request.Id)
	//fmt.Println("request.Costumer_id =", request.Costumer_id)
	//fmt.Println("request.Costumer_name =", request.Costumer_name)
	//fmt.Println("request.Costumer_email =", request.Costumer_email)
	//fmt.Println("request.Costume_id =", request.Costume_id)
	//fmt.Println("request.Costume_name =", request.Costume_name)
	//fmt.Println("request.Costume_category =", request.Costume_category)
	//fmt.Println("request.Costume_price =", request.Costume_price)
	//fmt.Println("request.Merchant_name =", request.Merchant_name)
	//fmt.Println("request.TotalAmount =", request.TotalAmount)

	var snapClient = snap.Client{}
	snapClient.New(server_key, midtrans.Sandbox)

	//log.Println("token", request.Id)
	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.Costumer_name,
			Email: request.Costumer_email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           strconv.Itoa(request.Costume_id),
				Name:         request.Costume_name,
				Price:        int64(request.Costume_price),
				Qty:          1,
				Category:     request.Costume_category,
				MerchantName: request.Merchant_name,
			},
			{
				Name:         "Tax From CosplayRent",
				Price:        int64(request.TotalAmount) - int64(request.Costume_price),
				Category:     "Tax",
				Qty:          1,
				MerchantName: "CosplayRent",
			},
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.Id.String(),
			GrossAmt: int64(request.TotalAmount),
		},
	}

	response, errSnap := snapClient.CreateTransaction(req)
	if errSnap != nil {
		panic(err)
	}

	midtransResponse := midtransWeb.MidtransResponse{
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse
}

func (service *MidtransServiceImpl) MidtransCallBack(ctx context.Context, orderid string) {
	log.Printf("Midtrans Callback with orderid:%s enter MidtransService: MidtransCallback", orderid)
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	service.MidtransRepository.Update(ctx, tx, orderid, &now)

	//// Find buyer id by orderid
	//id, err := service.MidtransRepository.FindBuyerId(ctx, tx, orderid)
	//if err != nil {
	//	panic(exception.NewNotFoundError(err.Error()))
	//}
	//
	//// Find seller id by order_id
	//id, err := serv
}

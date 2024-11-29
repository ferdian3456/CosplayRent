package midtrans

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	midtransWeb "cosplayrent/model/web/midtrans"
	orders "cosplayrent/model/web/order"
	midtranss "cosplayrent/repository/midtrans"
	"cosplayrent/repository/order"
	"cosplayrent/repository/topup_order"
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
	UserRepository       user.UserRepository
	OrderRepository      order.OrderRepository
	MidtransRepository   midtranss.MidtransRepository
	TopUpOrderRepository topup_order.TopUpOrderRepository
	DB                   *sql.DB
	Validate             *validator.Validate
}

func NewMidtransService(midtransRepository midtranss.MidtransRepository, topuporderRepository topup_order.TopUpOrderRepository, userRepository user.UserRepository, orderRepository order.OrderRepository, DB *sql.DB, validate *validator.Validate) MidtransService {
	return &MidtransServiceImpl{
		OrderRepository:      orderRepository,
		MidtransRepository:   midtransRepository,
		TopUpOrderRepository: topuporderRepository,
		UserRepository:       userRepository,
		DB:                   DB,
		Validate:             validate,
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
				ID:           "COSPLAYRENT-1-TAX",
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

func (service *MidtransServiceImpl) MidtransCallBack(ctx context.Context, orderid string, orderamount string) {
	log.Printf("Midtrans Callback with orderid:%s enter MidtransService: MidtransCallback", orderid)
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	orderAmount, err := strconv.ParseFloat(orderamount, 64)
	helper.PanicIfError(err)

	//log.Println("order id: ", orderid)
	userResult, errTopUpOrder := service.TopUpOrderRepository.FindUserIdByOrderId(ctx, tx, orderid)
	//log.Println("userResult: ", userResult)
	//log.Println("Error pertama: ", errTopUpOrder)
	buyerResult, errOrder := service.OrderRepository.FindBuyerIdByOrderId(ctx, tx, orderid)
	//log.Println("buyerResult: ", buyerResult)
	//log.Println("Error kedua: ", errOrder)
	//log.Println("masuk sini luar if")
	if errTopUpOrder != nil { // gk ketemu
		//log.Println("masuk sini dalam if")
		if errOrder != nil { // gk ketemu lagi
			//log.Println("masuk sini if - if")
			panic(exception.NewNotFoundError(errTopUpOrder.Error()))
		} else {
			//log.Println("masuk sini if - else")
			// lakuin update order
			service.MidtransRepository.Update(ctx, tx, orderid, &now)

			//// Find buyer id by orderid
			//buyerid, err := service.OrderRepository.FindBuyerIdByOrderId(ctx, tx, orderid)
			//if err != nil {
			//	panic(exception.NewNotFoundError(err.Error()))
			//}

			// Find seller id by order_id
			var sellerid string
			sellerid, err = service.OrderRepository.FindSellerIdByOrderId(ctx, tx, orderid)
			if err != nil {
				panic(exception.NewNotFoundError(err.Error()))
			}

			service.UserRepository.AfterBuy(ctx, tx, orderAmount, buyerResult, sellerid, &now)
			//log.Println("Success to buy")
		}
	} else {
		//log.Println("masuk sini dalam else")
		if errOrder == nil {
			//log.Println("masuk sini dalam else - if")
			panic(exception.NewNotFoundError(errOrder.Error()))
		} else {
			//log.Println("masuk sini dalam else - else")
			//lakuin update topup
			service.TopUpOrderRepository.UpdateTopUpOrder(ctx, tx, orderid, &now)

			service.UserRepository.TopUp(ctx, tx, orderAmount, userResult, &now)
			//log.Println("Success to top tup")

		}
	}
}

func (service *MidtransServiceImpl) CreateOrderTopUp(ctx context.Context, orderid string, username string, useremail string, uuid string, emoneyamount float64) midtransWeb.MidtransResponse {
	log.Printf("User with uuid: %s enter Midtrans Service: CreateOrderTopUp", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	server_key := os.Getenv("MIDTRANS_SERVER_KEY")

	var snapClient = snap.Client{}
	snapClient.New(server_key, midtrans.Sandbox)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			FName: username,
			Email: useremail,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           "COSPLAYRENT-2-TOPUP",
				Name:         "Top Up Emoney",
				Price:        int64(emoneyamount),
				Qty:          1,
				Category:     "Top-Up",
				MerchantName: "CosplayRent",
			},
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderid,
			GrossAmt: int64(emoneyamount),
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

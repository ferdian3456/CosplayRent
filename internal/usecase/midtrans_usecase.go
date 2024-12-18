package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	midtransWeb "cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/repository"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	googleuuid "github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog"
)

type MidtransUsecase struct {
	MidtransRepository *repository.MidtransRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewMidtransUsecase(midtransRepository *repository.MidtransRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *MidtransUsecase {
	return &MidtransUsecase{
		MidtransRepository: midtransRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (Usecase *MidtransUsecase) CreateTransaction(ctx context.Context, request order.DirectlyOrderToMidtrans, userUUID string) midtransWeb.MidtransResponse {
	log.Printf("User with uuid: %s enter Midtrans Usecase: CreateTransaction", userUUID)
	//fmt.Println("test")
	err := Usecase.Validate.Struct(request)
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
		Orderid:     request.Id,
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse
}

func (Usecase *MidtransUsecase) MidtransCallBack(ctx context.Context, orderid string, orderamount string) {
	log.Printf("Midtrans Callback with orderid:%s enter MidtransUsecase: MidtransCallback", orderid)
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	orderAmount, err := strconv.ParseFloat(orderamount, 64)
	helper.PanicIfError(err)

	//log.Println("order id: ", orderid)
	userResult, errTopUpOrder := Usecase.TopUpOrderRepository.FindUserIdByOrderId(ctx, tx, orderid)
	//log.Println("userResult: ", userResult)
	//log.Println("Error pertama: ", errTopUpOrder)
	buyerResult, errOrder := Usecase.OrderRepository.FindBuyerIdByOrderId(ctx, tx, orderid)
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
			Usecase.MidtransRepository.Update(ctx, tx, orderid, &now)

			//// Find buyer id by orderid
			//buyerid, err := Usecase.OrderRepository.FindBuyerIdByOrderId(ctx, tx, orderid)
			//if err != nil {
			//	panic(exception.NewNotFoundError(err.Error()))
			//}

			// Find seller id by order_id
			var sellerid string
			sellerid, err = Usecase.OrderRepository.FindSellerIdByOrderId(ctx, tx, orderid)
			if err != nil {
				panic(exception.NewNotFoundError(err.Error()))
			}

			Usecase.UserRepository.AfterBuy(ctx, tx, orderAmount, buyerResult, sellerid, &now)
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
			Usecase.TopUpOrderRepository.UpdateTopUpOrder(ctx, tx, orderid, &now)

			Usecase.UserRepository.TopUp(ctx, tx, orderAmount, userResult, &now)
			//log.Println("Success to top tup")

		}
	}
}

func (Usecase *MidtransUsecase) CreateOrderTopUp(ctx context.Context, orderid string, username string, useremail string, uuid string, emoneyamount float64) midtransWeb.MidtransResponse {
	log.Printf("User with uuid: %s enter Midtrans Usecase: CreateOrderTopUp", uuid)

	tx, err := Usecase.DB.Begin()
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

	finalOrderId, err := googleuuid.Parse(orderid)
	helper.PanicIfError(err)

	midtransResponse := midtransWeb.MidtransResponse{
		Orderid:     finalOrderId,
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse
}

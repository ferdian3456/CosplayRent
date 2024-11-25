package midtrans

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	midtransWeb "cosplayrent/model/web/midtrans"
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

func (service *MidtransServiceImpl) CreateTransaction(ctx context.Context, request midtransWeb.MidtransRequest, userUUID string) midtransWeb.MidtransResponse {
	log.Printf("User with uuid: %s enter Review Controller: FindUserReview", userUUID)
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	server_key := os.Getenv("MIDTRANS_SERVER_KEY")

	var snapClient = snap.Client{}
	snapClient.New(server_key, midtrans.Sandbox)

	log.Println("token", request.OrderId)
	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.CustomerName,
			Email: request.CustomerEmail,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           strconv.Itoa(request.CostumeId),
				Name:         request.CostumeName,
				Price:        request.Price,
				Category:     request.CostumeCategory,
				MerchantName: request.MerchantName,
			},
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderId,
			GrossAmt: request.FinalAmount,
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
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)
	service.MidtransRepository.Update(ctx, tx, orderid)
}

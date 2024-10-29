package midtrans

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	midtransWeb "cosplayrent/model/web/midtrans"
	midtranss "cosplayrent/repository/midtrans"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"log"
	"os"
)

type MidtransServiceImpl struct {
	MidtransRepository midtranss.MidtransRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewMidtransService(midtransRepository midtranss.MidtransRepository, DB *sql.DB, validate *validator.Validate) MidtransService {
	return &MidtransServiceImpl{
		MidtransRepository: midtransRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *MidtransServiceImpl) CreateTransaction(ctx context.Context, request midtransWeb.MidtransRequest, orderid string) midtransWeb.MidtransResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	server_key := os.Getenv("MIDTRANS_SERVER_KEY")

	var snapClient = snap.Client{}
	snapClient.New(server_key, midtrans.Sandbox)

	log.Println("token", orderid)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderid,
			GrossAmt: request.Amount,
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

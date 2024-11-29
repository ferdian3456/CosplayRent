package topup_order

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	midtranss "cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/user"
	"cosplayrent/repository/topup_order"
	users "cosplayrent/repository/user"
	"cosplayrent/service/midtrans"
	"database/sql"
	"github.com/go-playground/validator"
	googleuuid "github.com/google/uuid"
	"log"
	"time"
)

type TopUpOrderServiceImpl struct {
	TopUpOrderRepository topup_order.TopUpOrderRepository
	UserRepository       users.UserRepository
	MidtransService      midtrans.MidtransService
	DB                   *sql.DB
	Validate             *validator.Validate
}

func NewTopUpOrderService(topuporderRepository topup_order.TopUpOrderRepository, userRepository users.UserRepository, midtransService midtrans.MidtransService, DB *sql.DB, validate *validator.Validate) *TopUpOrderServiceImpl {
	return &TopUpOrderServiceImpl{
		TopUpOrderRepository: topuporderRepository,
		UserRepository:       userRepository,
		MidtransService:      midtransService,
		DB:                   DB,
		Validate:             validate,
	}
}

func (service *TopUpOrderServiceImpl) CreateTopUpOrder(ctx context.Context, topUpEMoneyRequest user.TopUpEmoney, uuid string) midtranss.MidtransResponse {
	log.Printf("User with uuid: %s enter User Service: TopUp", uuid)

	err := service.Validate.Struct(topUpEMoneyRequest)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	userResult, err := service.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()

	orderid := googleuuid.New()

	service.TopUpOrderRepository.CreateTopUpOrder(ctx, tx, orderid.String(), uuid, topUpEMoneyRequest, &now)

	midtransResult := service.MidtransService.CreateOrderTopUp(ctx, orderid.String(), userResult.Name, userResult.Email, uuid, topUpEMoneyRequest.Emoney_amont)

	return midtransResult
}

func (TopUpOrderServiceImpl) FindUserIdByOrderId(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

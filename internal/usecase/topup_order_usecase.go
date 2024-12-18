package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	midtransWeb "cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/model/web/topup_order"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/repository"
	"database/sql"
	"log"
	"time"

	"github.com/go-playground/validator"
	googleuuid "github.com/google/uuid"
	"github.com/rs/zerolog"
)

type TopUpOrderUsecase struct {
	TopUpOrderRepository *repository.TopUpOrderRepository
	DB                   *sql.DB
	Validate             *validator.Validate
	Log                  *zerolog.Logger
}

func NewTopUpOrderUsecase(topuporderRepository *repository.TopUpOrderRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *TopUpOrderUsecase {
	return &TopUpOrderUsecase{
		TopUpOrderRepository: topuporderRepository,
		DB:                   DB,
		Validate:             validate,
		Log:                  zerolog,
	}
}

func (Usecase *TopUpOrderUsecase) CreateTopUpOrder(ctx context.Context, topUpEMoneyRequest user.TopUpEmoney, uuid string) midtransWeb.MidtransResponse {
	log.Printf("User with uuid: %s enter User Usecase: TopUp", uuid)

	err := Usecase.Validate.Struct(topUpEMoneyRequest)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()

	orderid := googleuuid.New()

	Usecase.TopUpOrderRepository.CreateTopUpOrder(ctx, tx, orderid.String(), uuid, topUpEMoneyRequest, &now)

	midtransResult := Usecase.MidtransUsecase.CreateOrderTopUp(ctx, orderid.String(), userResult.Name, userResult.Email, uuid, topUpEMoneyRequest.Emoney_amont)

	return midtransResult
}

func (Usecase *TopUpOrderUsecase) CheckTopUpOrderByOrderId(ctx context.Context, orderID string) topup_order.TopupOrderResponse {
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	topuporderResult, err := Usecase.TopUpOrderRepository.CheckTopUpOrderByOrderId(ctx, tx, orderID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return topuporderResult
}

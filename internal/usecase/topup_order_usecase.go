package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/model/web/topup_order"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator"
	googleuuid "github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type TopUpOrderUsecase struct {
	UserRepository       *repository.UserRepository
	TopUpOrderRepository *repository.TopUpOrderRepository
	MidtransUsecase      *MidtransUsecase
	DB                   *sql.DB
	Validate             *validator.Validate
	Log                  *zerolog.Logger
	Config               *koanf.Koanf
}

func NewTopUpOrderUsecase(userRepository *repository.UserRepository, topUpOrderRepository *repository.TopUpOrderRepository, midtransUsecase *MidtransUsecase, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *TopUpOrderUsecase {
	return &TopUpOrderUsecase{
		UserRepository:       userRepository,
		TopUpOrderRepository: topUpOrderRepository,
		MidtransUsecase:      midtransUsecase,
		DB:                   DB,
		Validate:             validate,
		Log:                  zerolog,
		Config:               koanf,
	}
}

func (usecase *TopUpOrderUsecase) CreateTopUpOrder(ctx context.Context, userRequest user.TopUpEmoney, uuid string) (midtrans.MidtransResponse, error) {
	err := usecase.Validate.Struct(userRequest)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return midtrans.MidtransResponse{}, respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	user, err := usecase.UserRepository.FindNameAndEmailById(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return midtrans.MidtransResponse{}, err
	}

	now := time.Now()

	orderid := googleuuid.New()

	topuporder := domain.TopUpOrder{
		Id:           orderid.String(),
		User_id:      uuid,
		TopUp_amount: userRequest.Emoney_amount,
		Created_at:   &now,
		Updated_at:   &now,
	}

	usecase.TopUpOrderRepository.CreateTopUpOrder(ctx, tx, topuporder)

	midtransResult := usecase.MidtransUsecase.CreateOrderTopUp(ctx, topuporder, user)
	expiredTime := now.Add(24 * time.Hour)
	midtransResult.MidtransCreated_at = now.Format("2006-01-02 15:04:05")
	midtransResult.MidtransExpired = expiredTime.Format("2006-01-02 15:04:05")

	return midtransResult, nil
}

func (usecase *TopUpOrderUsecase) CheckTopUpOrderByOrderId(ctx context.Context, orderID string) (topup_order.TopupOrderResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	topuporderResult, err := usecase.TopUpOrderRepository.CheckTopUpOrderByOrderId(ctx, tx, orderID)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return topuporderResult, err
	}

	return topuporderResult, nil
}

package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator"
	googleuuid "github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type OrderUsecase struct {
	OrderRepository *repository.OrderRepository
	DB              *sql.DB
	Validate        *validator.Validate
	Log             *zerolog.Logger
	Config          *koanf.Koanf
}

func NewOrderUsecase(orderRepository *repository.OrderRepository, db *sql.DB, validator *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *OrderUsecase {
	return &OrderUsecase{
		OrderRepository: orderRepository,
		DB:              db,
		Validate:        validator,
		Log:             zerolog,
		Config:          koanf,
	}
}

func (usecase *OrderUsecase) Create(ctx context.Context, request order.OrderCreateRequest, useruuid string) error {
	err := usecase.Validate.Struct(request)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	uuid := googleuuid.New()
	orderDomain := domain.Order{
		Id:         uuid.String(),
		User_id:    useruuid,
		Costume_id: request.Costume_id,
		Total:      request.Total,
		Created_at: &now,
	}

	usecase.OrderRepository.Create(ctx, tx, orderDomain)
	return nil
}

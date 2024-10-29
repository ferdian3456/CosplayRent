package order

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/order"
	orders "cosplayrent/repository/order"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"log"
	"time"
)

type OrderServiceImpl struct {
	OrderRepository orders.OrderRepository
	DB              *sql.DB
	Validate        *validator.Validate
}

func NewOrderService(orderRepository orders.OrderRepository, DB *sql.DB, validate *validator.Validate) OrderService {
	return &OrderServiceImpl{
		OrderRepository: orderRepository,
		DB:              DB,
		Validate:        validate,
	}
}

func (service *OrderServiceImpl) Create(ctx context.Context, request order.OrderCreateRequest) {

	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)
	now := time.Now()
	log.Println("service")
	uuid := uuid.New()
	orderDomain := domain.Order{
		Id:         uuid.String(),
		User_id:    request.User_id,
		Costume_id: request.Costume_id,
		Total:      request.Total,
		Created_at: &now,
	}
	service.OrderRepository.Create(ctx, tx, orderDomain)
}

func (service *OrderServiceImpl) FindByUserId(ctx context.Context, uuid string) []order.OrderResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	order := []order.OrderResponse{}

	order, err = service.OrderRepository.FindByUserId(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return order
}

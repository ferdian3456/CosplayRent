package order

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/midtrans"
	"cosplayrent/model/web/order"
	orders "cosplayrent/repository/order"
	"cosplayrent/repository/user"
	midtranss "cosplayrent/service/midtrans"
	"database/sql"
	"github.com/go-playground/validator"
	googleUUID "github.com/google/uuid"
	"log"
	"time"
)

type OrderServiceImpl struct {
	OrderRepository orders.OrderRepository
	UserRepository  user.UserRepository
	MidtransService midtranss.MidtransService
	DB              *sql.DB
	Validate        *validator.Validate
}

func NewOrderService(orderRepository orders.OrderRepository, userRepository user.UserRepository, midtransService midtranss.MidtransService, DB *sql.DB, validate *validator.Validate) OrderService {
	return &OrderServiceImpl{
		OrderRepository: orderRepository,
		UserRepository:  userRepository,
		MidtransService: midtransService,
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
	uuid := googleUUID.New()
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

func (service *OrderServiceImpl) DirectlyOrderToMidtrans(ctx context.Context, uuid string, directOrderToMidtrans order.DirectlyOrderToMidtrans) midtrans.MidtransResponse {
	log.Printf("User with uuid: %s enter Order Service: DirectlyOrderToMidtrans", uuid)

	err := service.Validate.Struct(directOrderToMidtrans)
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

	SendOrderToDatabase := order.DirectlyOrderToMidtrans{
		Id:          googleUUID.New(),
		Costumer_id: directOrderToMidtrans.Costumer_id,
		Costume_id:  directOrderToMidtrans.Costume_id,
		TotalAmount: directOrderToMidtrans.TotalAmount,
		Created_at:  &now,
	}

	service.OrderRepository.DirectlyOrderToMidtrans(ctx, tx, uuid, SendOrderToDatabase)

	log.Println("hi")

	SendOrderToMidtrans := order.DirectlyOrderToMidtrans{
		Id:               SendOrderToDatabase.Id,
		Costumer_id:      userResult.Id,
		Costumer_name:    userResult.Name,
		Costumer_email:   userResult.Email,
		Costume_id:       directOrderToMidtrans.Costume_id,
		Costume_name:     directOrderToMidtrans.Costume_name,
		Costume_category: directOrderToMidtrans.Costume_category,
		Costume_price:    directOrderToMidtrans.Costume_price,
		Merchant_name:    directOrderToMidtrans.Merchant_name,
		TotalAmount:      directOrderToMidtrans.TotalAmount,
	}

	//log.Println("Send Order: ", SendOrderToDatabase)
	//log.Println("User Result: ", userResult)
	//
	//log.Println("Direct order: ", directOrderToMidtrans)

	//var SendOrderToMidtrans order.DirectlyOrderToMidtrans
	//
	//SendOrderToMidtrans.Id = SendOrderToDatabase.Id
	//fmt.Println("SendOrderToMidtrans.Id =", SendOrderToMidtrans.Id)
	//
	//SendOrderToMidtrans.Costumer_id = userResult.Id
	//fmt.Println("SendOrderToMidtrans.Costumer_id =", SendOrderToMidtrans.Costumer_id)
	//
	//SendOrderToMidtrans.Costumer_name = userResult.Name
	//fmt.Println("SendOrderToMidtrans.Costumer_name =", SendOrderToMidtrans.Costumer_name)
	//
	//SendOrderToMidtrans.Costumer_email = userResult.Email
	//fmt.Println("SendOrderToMidtrans.Costumer_email =", SendOrderToMidtrans.Costumer_email)
	//
	//SendOrderToMidtrans.Costume_id = directOrderToMidtrans.Costume_id
	//fmt.Println("SendOrderToMidtrans.Costume_id =", SendOrderToMidtrans.Costume_id)
	//
	//SendOrderToMidtrans.Costume_name = directOrderToMidtrans.Costume_name
	//fmt.Println("SendOrderToMidtrans.Costume_name =", SendOrderToMidtrans.Costume_name)
	//
	//SendOrderToMidtrans.Costume_category = directOrderToMidtrans.Costume_category
	//fmt.Println("SendOrderToMidtrans.Costume_category =", SendOrderToMidtrans.Costume_category)
	//
	//SendOrderToMidtrans.Costume_price = directOrderToMidtrans.Costume_price
	//fmt.Println("SendOrderToMidtrans.Costume_price =", SendOrderToMidtrans.Costume_price)
	//
	//SendOrderToMidtrans.Merchant_name = directOrderToMidtrans.Merchant_name
	//fmt.Println("SendOrderToMidtrans.Merchant_name =", SendOrderToMidtrans.Merchant_name)
	//
	//SendOrderToMidtrans.TotalAmount = directOrderToMidtrans.TotalAmount
	//fmt.Println("SendOrderToMidtrans.TotalAmount =", SendOrderToDatabase.TotalAmount)
	//fmt.Println("Uuid:", uuid)

	//log.Println("isi :%v", SendOrderToMidtrans)

	result := service.MidtransService.CreateTransaction(ctx, SendOrderToMidtrans, uuid)

	return result

	//return midtrans.MidtransResponse{}
}

package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/repository"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"
	googleUUID "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type OrderUsecase struct {
	OrderRepository *repository.OrderRepository
	DB              *sql.DB
	Validate        *validator.Validate
	Log             *zerolog.Logger
}

func NewOrderUsecase(orderRepository *repository.OrderRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *OrderUsecase {
	return &OrderUsecase{
		OrderRepository: orderRepository,
		DB:              DB,
		Validate:        validate,
		Log:             zerolog,
	}
}

func (Usecase *OrderUsecase) Create(ctx context.Context, request order.OrderCreateRequest) {

	err := Usecase.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
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
	Usecase.OrderRepository.Create(ctx, tx, orderDomain)
}

func (Usecase *OrderUsecase) FindByUserId(ctx context.Context, uuid string) []order.OrderResponse {
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	order := []order.OrderResponse{}

	order, err = Usecase.OrderRepository.FindByUserId(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return order
}

func (Usecase *OrderUsecase) DirectlyOrderToMidtrans(ctx context.Context, uuid string, directOrderToMidtrans order.DirectlyOrderToMidtrans) midtrans.MidtransResponse {
	log.Printf("User with uuid: %s enter Order Usecase: DirectlyOrderToMidtrans", uuid)

	err := Usecase.Validate.Struct(directOrderToMidtrans)
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

	SendOrderToDatabase := order.DirectlyOrderToMidtrans{
		Id:          googleUUID.New(),
		Costumer_id: directOrderToMidtrans.Costumer_id,
		Seller_id:   directOrderToMidtrans.Seller_id,
		Costume_id:  directOrderToMidtrans.Costume_id,
		TotalAmount: directOrderToMidtrans.TotalAmount,
		Created_at:  &now,
	}

	log.Println(SendOrderToDatabase.Id)

	Usecase.OrderRepository.DirectlyOrderToMidtrans(ctx, tx, uuid, SendOrderToDatabase)

	SendOrderToMidtrans := order.DirectlyOrderToMidtrans{
		Id:               SendOrderToDatabase.Id,
		Seller_id:        directOrderToMidtrans.Seller_id,
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

	result := Usecase.MidtransUsecase.CreateTransaction(ctx, SendOrderToMidtrans, uuid)

	return result

	//return midtrans.MidtransResponse{}
}

func (Usecase *OrderUsecase) FindOrderDetailByOrderId(ctx context.Context, orderid string) order.OrderResponse {
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	orderResult, err := Usecase.OrderRepository.FindOrderDetailByOrderId(ctx, tx, orderid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return orderResult
}

func (Usecase *OrderUsecase) GetAllSellerOrder(ctx context.Context, sellerid string) []order.AllSellerOrderResponse {
	log.Printf("User with uuid: %s enter Order Usecase: GetAllSellerOrder", sellerid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, sellerid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	sellerOrderResult, err := Usecase.OrderRepository.FindOrderBySellerId(ctx, tx, sellerid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range sellerOrderResult {
		costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, sellerOrderResult[i].Costume_id)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		sellerOrderResult[i].Costume_name = costumeResult.Name
		sellerOrderResult[i].Costume_id = costumeResult.Id
		sellerOrderResult[i].Costume_price = costumeResult.Price
		sellerOrderResult[i].Costume_size = costumeResult.Ukuran
		if costumeResult.Picture != nil {
			value := imageEnv + *costumeResult.Picture
			sellerOrderResult[i].Costume_picture = &value
		}
	}

	return sellerOrderResult
}

func (Usecase *OrderUsecase) GetAllUserOrder(ctx context.Context, userid string) []order.AllUserOrderResponse {
	log.Printf("User with uuid: %s enter Order Usecase: GetAllUserOrder", userid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userOrderResult, err := Usecase.OrderRepository.FindOrderByUserId(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range userOrderResult {
		costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, userOrderResult[i].Costume_id)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		userOrderResult[i].Costume_name = costumeResult.Name
		userOrderResult[i].Costume_id = costumeResult.Id
		userOrderResult[i].Costume_price = costumeResult.Price
		userOrderResult[i].Costume_size = costumeResult.Ukuran
		if costumeResult.Picture != nil {
			value := imageEnv + *costumeResult.Picture
			userOrderResult[i].Costume_picture = &value
		}
	}

	return userOrderResult
}

func (Usecase *OrderUsecase) UpdateSellerOrder(ctx context.Context, updateRequest order.OrderUpdateRequest, sellerid string, orderid string) {
	log.Printf("User with uuid: %s enter Order Usecase: UpdateSellerOrder", sellerid)

	err := Usecase.Validate.Struct(updateRequest)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, sellerid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()
	updateRequest.Updated_at = &now

	Usecase.OrderRepository.UpdateSellerOrder(ctx, tx, updateRequest, sellerid, orderid)
}

func (Usecase *OrderUsecase) GetDetailOrderByOrderId(ctx context.Context, sellerid string, orderid string) order.OrderDetailByOrderIdResponse {
	log.Printf("User with uuid: %s enter Order Usecase: GetDetailOrderByOrderId", sellerid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, sellerid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	orderResult, err := Usecase.OrderRepository.FindOrderDetailByOrderId(ctx, tx, orderid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, orderResult.User_id.String())
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, orderResult.Costume_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userIdentityCardPicture, err := Usecase.UserRepository.GetIdentityCard(ctx, tx, userResult.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	if costumeResult.Picture != nil {
		value := imageEnv + *costumeResult.Picture
		costumeResult.Picture = &value
	}

	if userIdentityCardPicture != "" {
		value := imageEnv + userIdentityCardPicture
		userIdentityCardPicture = value
	}

	orderResponse := order.OrderDetailByOrderIdResponse{
		Costume_name:             costumeResult.Name,
		Costume_price:            costumeResult.Price,
		Costume_size:             costumeResult.Ukuran,
		Costume_picture:          costumeResult.Picture,
		Costumer_name:            userResult.Name,
		Costumer_address:         userResult.Address,
		Costumer_origin_province: userResult.Origin_province_name,
		Costumer_origin_city:     userResult.Origin_city_name,
		Costumer_identity_card:   userIdentityCardPicture,
	}

	return orderResponse
}

func (Usecase *OrderUsecase) GetUserDetailOrder(ctx context.Context, userid string, orderid string) order.GetUserOrderDetailResponse {
	log.Printf("User with uuid: %s enter Order Usecase: GetUserOrder", userid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	orderResult, err := Usecase.OrderRepository.FindOrderDetailByOrderId(ctx, tx, orderid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, orderResult.Seller_id.String())
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, orderResult.Costume_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	if costumeResult.Picture != nil {
		value := imageEnv + *costumeResult.Picture
		costumeResult.Picture = &value
	}

	if orderResult.Description != nil {
		value := *orderResult.Description
		orderResult.Description = &value
	}

	orderResponse := order.GetUserOrderDetailResponse{
		Costume_name:           costumeResult.Name,
		Costume_price:          costumeResult.Price,
		Costume_size:           costumeResult.Ukuran,
		Costume_picture:        costumeResult.Picture,
		Seller_name:            userResult.Name,
		Seller_address:         userResult.Address,
		Seller_origin_province: userResult.Origin_province_name,
		Seller_origin_city:     userResult.Origin_city_name,
		Seller_response:        orderResult.Description,
	}

	return orderResponse
}

func (Usecase *OrderUsecase) CheckBalanceWithOrderAmount(ctx context.Context, checkbalance order.CheckBalanceWithOrderAmount, uuid string) order.CheckBalanceWithOrderAmountReponse {
	log.Printf("User with uuid: %s enter Order Usecase: CheckBalanceWithOrderAmount", uuid)

	err := Usecase.Validate.Struct(checkbalance)
	helper.PanicIfError(err)

	log.Println(checkbalance.Order_amount)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	emoneyResult, err := Usecase.UserRepository.GetEMoneyAmount(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if emoneyResult.Emoney_amont <= checkbalance.Order_amount {
		panic(exception.NewNotFoundError("your money is not sufficient"))
	}

	CheckBalanceResult := order.CheckBalanceWithOrderAmountReponse{
		Status_to_order: "true",
	}

	return CheckBalanceResult
}

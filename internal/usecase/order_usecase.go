package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator"
	googleuuid "github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type OrderUsecase struct {
	UserRepository    *repository.UserRepository
	CostumeRepository *repository.CostumeRepository
	OrderRepository   *repository.OrderRepository
	MidtransUsecase   *MidtransUsecase
	DB                *sql.DB
	Validate          *validator.Validate
	Log               *zerolog.Logger
	Config            *koanf.Koanf
}

func NewOrderUsecase(userRepository *repository.UserRepository, costumeRepository *repository.CostumeRepository, orderRepository *repository.OrderRepository, midtransUsecase *MidtransUsecase, db *sql.DB, validator *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *OrderUsecase {
	return &OrderUsecase{
		UserRepository:    userRepository,
		CostumeRepository: costumeRepository,
		OrderRepository:   orderRepository,
		MidtransUsecase:   midtransUsecase,
		DB:                db,
		Validate:          validator,
		Log:               zerolog,
		Config:            koanf,
	}
}

func (usecase *OrderUsecase) Create(ctx context.Context, uuid string, userRequest order.DirectlyOrderToMidtrans) (midtrans.MidtransResponse, error) {
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

	userResult, err := usecase.UserRepository.FindBasicInfo(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return midtrans.MidtransResponse{}, err
	}

	now := time.Now()

	SendOrderToDatabase := order.DirectlyOrderToMidtrans{
		Id:          googleuuid.New().String(),
		Costumer_id: userRequest.Costumer_id,
		Seller_id:   userRequest.Seller_id,
		Costume_id:  userRequest.Costume_id,
		TotalAmount: userRequest.TotalAmount,
		Created_at:  &now,
		Updated_at:  &now,
	}

	fmt.Println("data:", SendOrderToDatabase)

	usecase.OrderRepository.Create(ctx, tx, SendOrderToDatabase)

	SendOrderToMidtrans := order.DirectlyOrderToMidtrans{
		Id:               SendOrderToDatabase.Id,
		Seller_id:        userRequest.Seller_id,
		Costumer_id:      userResult.Id,
		Costumer_name:    userResult.Name,
		Costumer_email:   userResult.Email,
		Costume_id:       userRequest.Costume_id,
		Costume_name:     userRequest.Costume_name,
		Costume_category: userRequest.Costume_category,
		Costume_price:    userRequest.Costume_price,
		Merchant_name:    userRequest.Merchant_name,
		TotalAmount:      userRequest.TotalAmount,
	}

	result := usecase.MidtransUsecase.CreateTransaction(ctx, SendOrderToMidtrans)

	return result, nil
}

func (usecase *OrderUsecase) CheckStatusPayment(ctx context.Context, orderid string) (bool, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	orderResult, err := usecase.OrderRepository.CheckStatusPayment(ctx, tx, orderid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return *orderResult, err
	}

	return *orderResult, nil
}

func (usecase *OrderUsecase) GetAllSellerOrder(ctx context.Context, sellerid string) ([]order.AllSellerOrderResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	sellerOrderResult, err := usecase.OrderRepository.FindOrderBySellerId(ctx, tx, sellerid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return sellerOrderResult, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range sellerOrderResult {
		costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, sellerOrderResult[i].Costume_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return sellerOrderResult, err
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

	return sellerOrderResult, nil
}

func (usecase *OrderUsecase) GetDetailOrderByOrderId(ctx context.Context, sellerid string, orderid string) (order.OrderDetailByOrderIdResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	orderResult, err := usecase.OrderRepository.FindUserAndCostumeById(ctx, tx, orderid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.OrderDetailByOrderIdResponse{}, err
	}

	userResult, err := usecase.UserRepository.FindByUUID(ctx, tx, orderResult.User_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.OrderDetailByOrderIdResponse{}, err
	}

	costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, orderResult.Costume_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.OrderDetailByOrderIdResponse{}, err
	}

	userIdentityCardPicture, err := usecase.UserRepository.GetIdentityCard(ctx, tx, userResult.Id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.OrderDetailByOrderIdResponse{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

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

	return orderResponse, nil
}

func (usecase *OrderUsecase) GetUserDetailOrder(ctx context.Context, userid string, orderid string) (order.GetUserOrderDetailResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	orderResult, err := usecase.OrderRepository.FindSellerAndCostumeById(ctx, tx, orderid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.GetUserOrderDetailResponse{}, err
	}

	userResult, err := usecase.UserRepository.FindByUUID(ctx, tx, orderResult.Seller_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.GetUserOrderDetailResponse{}, err
	}

	costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, orderResult.Costume_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.GetUserOrderDetailResponse{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	if costumeResult.Picture != nil {
		value := imageEnv + *costumeResult.Picture
		costumeResult.Picture = &value
	}

	if orderResult.Description != "" {
		value := orderResult.Description
		orderResult.Description = value
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
		Seller_response:        &orderResult.Description,
	}

	return orderResponse, nil
}

func (usecase *OrderUsecase) GetAllUserOrder(ctx context.Context, userid string) ([]order.AllUserOrderResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	userOrderResult, err := usecase.OrderRepository.FindOrderByUserId(ctx, tx, userid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return userOrderResult, err
	}

	imageEnv := os.Getenv("application.image_env")

	for i := range userOrderResult {
		costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, userOrderResult[i].Costume_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return userOrderResult, err
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

	return userOrderResult, nil
}

func (usecase *OrderUsecase) CheckBalanceWithOrderAmount(ctx context.Context, userRequest order.CheckBalanceWithOrderAmount, uuid string) (order.CheckBalanceWithOrderAmountReponse, error) {
	err := usecase.Validate.Struct(userRequest)
	if err != nil {
		respErr := errors.New("invalid request body")
		usecase.Log.Warn().Err(respErr).Msg(err.Error())
		return order.CheckBalanceWithOrderAmountReponse{}, respErr
	}

	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	emoneyResult := usecase.UserRepository.GetEMoneyAmount(ctx, tx, uuid)

	if emoneyResult.Emoney_amont <= userRequest.Order_amount {
		respErr := errors.New("your money is not sufficient")
		usecase.Log.Warn().Err(respErr).Msg(respErr.Error())
		return order.CheckBalanceWithOrderAmountReponse{}, respErr
	}

	CheckBalanceResult := order.CheckBalanceWithOrderAmountReponse{
		Status_to_order: "true",
	}

	return CheckBalanceResult, nil
}

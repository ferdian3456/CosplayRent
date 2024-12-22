package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	midtransWeb "cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog"
)

type MidtransUsecase struct {
	UserRepository       *repository.UserRepository
	OrderRepository      *repository.OrderRepository
	TopUpOrderRepository *repository.TopUpOrderRepository
	DB                   *sql.DB
	Validate             *validator.Validate
	Log                  *zerolog.Logger
	Config               *koanf.Koanf
}

func NewMidtransUsecase(userRepository *repository.UserRepository, orderRepository *repository.OrderRepository, topUpOrderRepository *repository.TopUpOrderRepository, db *sql.DB, validator *validator.Validate, zerolog *zerolog.Logger, config *koanf.Koanf) *MidtransUsecase {
	return &MidtransUsecase{
		UserRepository:       userRepository,
		OrderRepository:      orderRepository,
		TopUpOrderRepository: topUpOrderRepository,
		DB:                   db,
		Validate:             validator,
		Log:                  zerolog,
		Config:               config,
	}
}

func (usecase *MidtransUsecase) CreateTransaction(ctx context.Context, userRequest order.DirectlyOrderToMidtrans) midtransWeb.MidtransResponse {
	server_key := usecase.Config.String("payment_gateway.midtrans.server_key")

	//fmt.Println("request.Id =", request.Id)
	//fmt.Println("request.Costumer_id =", request.Costumer_id)
	//fmt.Println("request.Costumer_name =", request.Costumer_name)
	//fmt.Println("request.Costumer_email =", request.Costumer_email)
	//fmt.Println("request.Costume_id =", request.Costume_id)
	//fmt.Println("request.Costume_name =", request.Costume_name)
	//fmt.Println("request.Costume_category =", request.Costume_category)
	//fmt.Println("request.Costume_price =", request.Costume_price)
	//fmt.Println("request.Merchant_name =", request.Merchant_name)
	//fmt.Println("request.TotalAmount =", request.TotalAmount)

	//midtrans.DefaultLoggerLevel = &midtrans.LoggerImplementation{LogLevel: midtrans.NoLogging}

	var snapClient = snap.Client{}
	snapClient.New(server_key, midtrans.Sandbox)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			FName: userRequest.Costumer_name,
			Email: userRequest.Costumer_email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           strconv.Itoa(userRequest.Costume_id),
				Name:         userRequest.Costume_name,
				Price:        int64(userRequest.Costume_price),
				Qty:          1,
				Category:     userRequest.Costume_category,
				MerchantName: userRequest.Merchant_name,
			},
			{
				ID:           "COSPLAYRENT-1-TAX",
				Name:         "Tax From CosplayRent",
				Price:        int64(userRequest.TotalAmount) - int64(userRequest.Costume_price),
				Category:     "Tax",
				Qty:          1,
				MerchantName: "CosplayRent",
			},
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  userRequest.Id,
			GrossAmt: int64(userRequest.TotalAmount),
		},
	}

	response, err := snapClient.CreateTransaction(req)
	if err != nil {
		respErr := errors.New("failed to create midtrans transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	midtransResponse := midtransWeb.MidtransResponse{
		Orderid:     userRequest.Id,
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse
}

func (usecase *MidtransUsecase) MidtransCallBack(ctx context.Context, midtransWeb midtransWeb.MidtransCallback) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	floatOrderAmount, err := strconv.ParseFloat(midtransWeb.GrossAmount, 64)
	if err != nil {
		respErr := errors.New("failed to parse midtrans's gross amount string into float64")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	userResult, errTopUpOrder := usecase.TopUpOrderRepository.FindUserIdByOrderId(ctx, tx, midtransWeb.OrderID)
	buyerResult, errOrder := usecase.OrderRepository.FindBuyerIdByOrderId(ctx, tx, midtransWeb.OrderID)

	midtransDomain := domain.Midtrans{
		Order_id:      midtransWeb.OrderID,
		Order_amount:  floatOrderAmount,
		TopUpUser_id:  userResult,
		OrderBuyer_id: buyerResult,
		Updated_at:    &now,
	}

	if errTopUpOrder != nil {
		if errOrder != nil {
			return
		} else {
			usecase.OrderRepository.Update(ctx, tx, midtransDomain)

			var sellerid string

			midtransDomain.OrderSeller_id = sellerid
			sellerid, err = usecase.OrderRepository.FindSellerIdByOrderId(ctx, tx, midtransDomain.Order_id)
			if err != nil {
				return
			}

			usecase.UserRepository.AfterBuy(ctx, tx, midtransDomain)
		}
	} else {
		if errOrder == nil {
			return
		} else {
			usecase.TopUpOrderRepository.Update(ctx, tx, midtransDomain)

			usecase.UserRepository.TopUp(ctx, tx, midtransDomain)
		}
	}
}

func (usecase *MidtransUsecase) CreateOrderTopUp(ctx context.Context, topuporder domain.TopUpOrder, user domain.User) midtransWeb.MidtransResponse {
	server_key := usecase.Config.String("payment_gateway.midtrans.server_key")

	midtrans.DefaultLoggerLevel = &midtrans.LoggerImplementation{LogLevel: midtrans.NoLogging}

	var snapClient = snap.Client{}
	snapClient.New(server_key, midtrans.Sandbox)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           "COSPLAYRENT-2-TOPUP",
				Name:         "Top Up Emoney",
				Price:        int64(topuporder.TopUp_amount),
				Qty:          1,
				Category:     "Top-Up",
				MerchantName: "CosplayRent",
			},
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  topuporder.Id,
			GrossAmt: int64(topuporder.TopUp_amount),
		},
	}

	response, err := snapClient.CreateTransaction(req)
	if err != nil {
		respErr := errors.New("failed to create midtrans transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	midtransResponse := midtransWeb.MidtransResponse{
		Orderid:     topuporder.Id,
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse
}

package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/midtrans"
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
	UserRepository     *repository.UserRepository
	CostumeRepository  *repository.CostumeRepository
	CategoryRepository *repository.CategoryRepository
	OrderRepository    *repository.OrderRepository
	MidtransUsecase    *MidtransUsecase
	DB                 *sql.DB
	Validate           *validator.Validate
	Log                *zerolog.Logger
	Config             *koanf.Koanf
}

func NewOrderUsecase(userRepository *repository.UserRepository, costumeRepository *repository.CostumeRepository, categoryRepository *repository.CategoryRepository, orderRepository *repository.OrderRepository, midtransUsecase *MidtransUsecase, db *sql.DB, validator *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *OrderUsecase {
	return &OrderUsecase{
		UserRepository:     userRepository,
		CostumeRepository:  costumeRepository,
		CategoryRepository: categoryRepository,
		OrderRepository:    orderRepository,
		MidtransUsecase:    midtransUsecase,
		DB:                 db,
		Validate:           validator,
		Log:                zerolog,
		Config:             koanf,
	}
}

func (usecase *OrderUsecase) Create(ctx context.Context, uuid string, userRequest order.OrderRequest) (midtrans.MidtransResponse, error) {
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

	orderid := googleuuid.New().String()

	orderToDatabase := domain.Order{
		Id:                   orderid,
		Costumer_id:          uuid,
		Seller_id:            userRequest.Seller_id,
		Costume_id:           userRequest.Costume_id,
		Total_amount:         userRequest.TotalAmount - 3000,
		Shipment_origin:      userRequest.Shipment_origin,
		Shipment_destination: userRequest.Shippment_destination,
		Created_at:           &now,
		Updated_at:           &now,
	}

	SendOrderToMidtrans := domain.OrderToMidtrans{
		Id:               orderid,
		Seller_id:        userRequest.Seller_id,
		Seller_name:      userRequest.Seller_name,
		Costumer_id:      userResult.Id,
		Costumer_name:    userResult.Name,
		Costumer_email:   userResult.Email,
		Costume_id:       userRequest.Costume_id,
		Costume_name:     userRequest.Costume_name,
		Costume_category: userRequest.Costume_category,
		Costume_price:    userRequest.Costume_price,
		Total_amount:     userRequest.TotalAmount,
	}

	payment := domain.Payments{
		Order_id:       orderid,
		Customer_id:    uuid,
		Seller_id:      userRequest.Seller_id,
		Status:         "Pending",
		Amount:         userRequest.TotalAmount - 3000,
		Payment_method: userRequest.Payment_method,
		Created_at:     &now,
		Updated_at:     &now,
	}

	event := domain.OrderEvents{
		User_id:    uuid,
		Order_id:   orderid,
		Status:     "Paid",
		Created_at: &now,
	}

	err = usecase.UserRepository.CheckUserExistance(ctx, tx, userRequest.Seller_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return midtrans.MidtransResponse{}, err
	}

	err = usecase.CostumeRepository.CheckCostume(ctx, tx, userRequest.Seller_id, userRequest.Costume_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return midtrans.MidtransResponse{}, err
	}

	usecase.OrderRepository.Create(ctx, tx, orderToDatabase)

	if userRequest.Payment_method == "Emoney" {
		usecase.UserRepository.AfterBuy(ctx, tx, userRequest.TotalAmount, &now, uuid, userRequest.Seller_id)
		payment.Status = "Paid"
		usecase.OrderRepository.CreatePayment(ctx, tx, payment)
		usecase.OrderRepository.CreateOrderEvents(ctx, tx, event)
		return midtrans.MidtransResponse{}, nil
	}

	result := usecase.MidtransUsecase.CreateTransaction(ctx, SendOrderToMidtrans)
	payment.Midtrans_redirect_url = result.RedirectUrl
	payment.Midtrans_url_expired_time = now.Add(24 * time.Hour)
	usecase.OrderRepository.CreatePayment(ctx, tx, payment)
	expiredTime := now.Add(24 * time.Hour)

	result.MidtransExpired = expiredTime.Format("2006-01-02 15:04:05")

	return result, nil
}

func (usecase *OrderUsecase) CheckStatusPayment(ctx context.Context, orderid string) (string, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	orderResult, err := usecase.OrderRepository.CheckStatusPayment(ctx, tx, orderid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return "", err
	}

	return orderResult, nil
}

func (usecase *OrderUsecase) FindPaymentInfoByPaymentId(ctx context.Context, uuid string, paymentid int) (order.PaymentInfo, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	paymentResult, err := usecase.OrderRepository.FindPaymentInfoByPaymentId(ctx, tx, paymentid, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.PaymentInfo{}, err
	}

	formattedDate := paymentResult.Midtrans_url_expired_time.Format("2006-01-02 15:04:05")

	paymentFormatted := order.PaymentInfo{
		Payment_amount:                     paymentResult.Amount,
		Midtrans_redirect_url:              paymentResult.Midtrans_redirect_url,
		Midtrans_redirect_url_expired_time: formattedDate,
		Status:                             paymentResult.Status,
	}

	return paymentFormatted, nil
}

func (usecase *OrderUsecase) FindListPaymentTransaction(ctx context.Context, uuid string) ([]order.PaymentTransationForOrderResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	orders, err := usecase.OrderRepository.FindListOrderByCostumeId(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return []order.PaymentTransationForOrderResponse{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range orders {
		costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, orders[i].Costume_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return []order.PaymentTransationForOrderResponse{}, err
		}

		orders[i].Costume_name = costumeResult.Name
		orders[i].Costume_price = costumeResult.Price
		orders[i].Costume_size = *costumeResult.Ukuran
		if costumeResult.Picture != nil {
			value := imageEnv + *costumeResult.Picture
			orders[i].Costume_picture = &value
		}

		payment, err := usecase.OrderRepository.FindPaymentInfoByOrderId(ctx, tx, orders[i].Order_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return []order.PaymentTransationForOrderResponse{}, err
		}

		orders[i].Payment_id = payment.Id
		orders[i].Payment_status = payment.Status

		orders[i].Midtrans_redirect_url_expired_time = payment.Midtrans_url_expired_time.Format("2006-01-02 15:04:05")
	}

	return orders, nil
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

	userResult, err := usecase.UserRepository.FindByUUID(ctx, tx, orderResult.Costumer_id)
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

	eventResult, err := usecase.OrderRepository.FindEventInfoById(ctx, tx, orderid)
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
		Shipment_destination:     orderResult.Shipment_destination,
		Costumer_identity_card:   userIdentityCardPicture,
		Shipment_receipt_user_id: eventResult.Shipment_receipt_user_id,
		Shipment_notes:           eventResult.Notes,
	}

	return orderResponse, nil
}

func (usecase *OrderUsecase) CreateOrderEvent(ctx context.Context, uuid string, orderRequest order.OrderEventRequest, orderId string) error {
	err := usecase.Validate.Struct(orderRequest)
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

	orderEvent := domain.OrderEvents{
		User_id:                  uuid,
		Order_id:                 orderId,
		Status:                   orderRequest.OrderEventStatus,
		Notes:                    orderRequest.OrderEventNotes,
		Shipment_receipt_user_id: orderRequest.Shipment_receipt_user_id,
		Created_at:               &now,
	}

	usecase.OrderRepository.CreateOrderEvents(ctx, tx, orderEvent)
	return nil
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

	eventResult, err := usecase.OrderRepository.FindEventInfoById(ctx, tx, orderResult.Id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return order.GetUserOrderDetailResponse{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	if costumeResult.Picture != nil {
		value := imageEnv + *costumeResult.Picture
		costumeResult.Picture = &value
	}

	orderResponse := order.GetUserOrderDetailResponse{
		Costume_name:             costumeResult.Name,
		Costume_price:            costumeResult.Price,
		Costume_size:             costumeResult.Ukuran,
		Costume_picture:          costumeResult.Picture,
		Seller_name:              userResult.Name,
		Seller_address:           userResult.Address,
		Shipment_receipt_user_id: eventResult.Shipment_receipt_user_id,
		Shipment_notes:           eventResult.Notes,
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

	imageEnv := usecase.Config.String("IMAGE_ENV")

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

func (usecase *OrderUsecase) UpdateOrder(ctx context.Context, userRequest order.OrderUpdateRequest, uuid string, orderid string) error {
	err := usecase.Validate.Struct(userRequest)
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

	order := domain.Order{
		Id:        orderid,
		Seller_id: uuid,
	}

	err = usecase.OrderRepository.CheckIfUserOrSeller(ctx, tx, order.Seller_id, order.Id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	usecase.OrderRepository.UpdateOrder(ctx, tx, order)

	return nil
}

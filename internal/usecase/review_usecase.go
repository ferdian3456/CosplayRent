package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/review"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"time"
)

type ReviewUsecase struct {
	UserRepository    *repository.UserRepository
	CostumeRepository *repository.CostumeRepository
	ReviewRepository  *repository.ReviewRepository
	OrderRepository   *repository.OrderRepository
	DB                *sql.DB
	Validate          *validator.Validate
	Log               *zerolog.Logger
	Config            *koanf.Koanf
}

func NewReviewUsecase(userRepository *repository.UserRepository, costumeRepository *repository.CostumeRepository, reviewRepository *repository.ReviewRepository, db *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *ReviewUsecase {
	return &ReviewUsecase{
		UserRepository:    userRepository,
		CostumeRepository: costumeRepository,
		ReviewRepository:  reviewRepository,
		DB:                db,
		Validate:          validate,
		Log:               zerolog,
		Config:            koanf,
	}
}

func (usecase *ReviewUsecase) Create(ctx context.Context, request review.ReviewCreateRequest, uuid string) error {
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

	err = usecase.OrderRepository.CheckOrderAndCostumeId(ctx, tx, request.Order_id, request.Costume_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	now := time.Now()

	review := domain.Review{
		Order_id:       request.Order_id,
		User_id:        request.Customer_id,
		Costume_id:     request.Costume_id,
		Description:    request.Description,
		Review_picture: *request.Review_picture,
		Rating:         request.Rating,
		Created_at:     &now,
		Updated_at:     &now,
	}

	usecase.ReviewRepository.Create(ctx, tx, review)

	return nil
}

func (usecase *ReviewUsecase) Update(ctx context.Context, request review.ReviewUpdateRequest, uuid string, reviewid int) error {
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

	err = usecase.ReviewRepository.CheckReviewId(ctx, tx, reviewid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	now := time.Now()

	review := domain.Review{
		Id:             reviewid,
		Review_picture: request.Review_picture,
		Description:    request.Description,
		Rating:         request.Rating,
		Updated_at:     &now,
	}

	usecase.ReviewRepository.Update(ctx, tx, review)

	return nil
}

func (usecase *ReviewUsecase) FindUserReview(ctx context.Context, uuid string) ([]review.UserReviewResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	orderResult, err := usecase.OrderRepository.FindOrderInfoByUserId(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return orderResult, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range orderResult {
		costumeResult, err := usecase.CostumeRepository.FindCostumeInfoById(ctx, tx, orderResult[i].Custome_Id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return orderResult, err
		}

		reviewPicture := usecase.ReviewRepository.FindReviewPictureById(ctx, tx, orderResult[i].Order_id)

		sellerNameResult, err := usecase.UserRepository.FindNameById(ctx, tx, orderResult[i].Seller_id)
		orderResult[i].Seller_name = sellerNameResult
		orderResult[i].Costume_name = costumeResult.Name
		orderResult[i].Costume_picture = imageEnv + costumeResult.Costume_picture
		orderResult[i].Costume_size = costumeResult.Costume_size
		orderResult[i].Costume_weight = costumeResult.Costume_weight
		orderResult[i].Review_picture = reviewPicture
	}

	return orderResult, nil
}

func (usecase *ReviewUsecase) FindByCostumeId(ctx context.Context, id int) ([]review.ReviewResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	review := []review.ReviewResponse{}
	review, err = usecase.ReviewRepository.FindByCostumeId(ctx, tx, id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return review, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range review {
		userResult, err := usecase.UserRepository.FindNameAndProfile(ctx, tx, review[i].User_id)
		helper.PanicIfError(err)
		review[i].Name = userResult.Name
		if userResult.Profile_picture != nil {
			value := imageEnv + *userResult.Profile_picture
			review[i].Profile_picture = &value
		}
		value := imageEnv + *review[i].Review_picture
		review[i].Review_picture = &value
	}

	return review, nil
}

func (usecase *ReviewUsecase) FindAllUserReview(ctx context.Context, userid string) ([]review.ListOfNonReviewedOrder, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	userOrderResult, err := usecase.OrderRepository.FindAllUserOrder(ctx, tx, userid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return []review.ListOfNonReviewedOrder{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	var fixedOrderResult []domain.Review

	for i := range userOrderResult {
		nonReviewedOrders := usecase.OrderRepository.FindAllNonReviewedOrder(ctx, tx, userOrderResult[i].Order_id)
		fixedOrderResult = append(fixedOrderResult, nonReviewedOrders...)
	}

	var reviewResponse []review.ListOfNonReviewedOrder

	if len(fixedOrderResult) == 0 {
		errors := errors.New("order not found")
		usecase.Log.Warn().Msg(errors.Error())
		return []review.ListOfNonReviewedOrder{}, errors
	}

	for _, fixedOrder := range fixedOrderResult {
		costumeIdResult, err := usecase.OrderRepository.FindCostumeIdByOrderId(ctx, tx, fixedOrder.Order_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return reviewResponse, err
		}

		costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, costumeIdResult)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return reviewResponse, err
		}

		review := review.ListOfNonReviewedOrder{
			Order_id:      fixedOrder.Order_id,
			Costume_name:  costumeResult.Name,
			Costume_id:    costumeResult.Id,
			Costume_price: costumeResult.Price,
			Costume_size:  *costumeResult.Ukuran,
			Order_date:    fixedOrder.Created_at.Format("2006-01-02 15:04:05"),
		}

		if costumeResult.Picture != nil {
			value := imageEnv + *costumeResult.Picture
			review.Costume_picture = &value
		}

		reviewResponse = append(reviewResponse, review)
	}

	return reviewResponse, nil
}

func (usecase *ReviewUsecase) FindAllReviewedOrder(ctx context.Context, userid string) ([]review.ListOfReviewedOrder, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	userReviewResult, err := usecase.ReviewRepository.FindAllReviewedOrder(ctx, tx, userid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return []review.ListOfReviewedOrder{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	var reviewResponse []review.ListOfReviewedOrder

	for _, reviewResult := range userReviewResult {
		costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, reviewResult.Costume_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return []review.ListOfReviewedOrder{}, err
		}

		review := review.ListOfReviewedOrder{
			Order_id:           reviewResult.Order_id,
			Costume_name:       costumeResult.Name,
			Costume_id:         costumeResult.Id,
			Review_date:        reviewResult.Created_at.Format("2006-01-02 15:04:05"),
			Review_description: reviewResult.Description,
			Review_rating:      reviewResult.Rating,
		}

		if costumeResult.Picture != nil {
			value := imageEnv + *costumeResult.Picture
			review.Costume_picture = &value
		}

		reviewResponse = append(reviewResponse, review)
	}

	return reviewResponse, nil
}

func (usecase *ReviewUsecase) FindReviewInfoByOrderId(ctx context.Context, userid string, orderid string) (review.ReviewInfo, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	reviewInfoResponse, err := usecase.ReviewRepository.FindReviewInfoByOrderId(ctx, tx, userid, orderid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return review.ReviewInfo{}, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	reviewInfoResponse.Review_picture = imageEnv + reviewInfoResponse.Review_picture

	return reviewInfoResponse, nil
}

func (usecase *ReviewUsecase) DeleteUserReviewByReviewID(ctx context.Context, userid string, reviewid int) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	//err = usecase.ReviewRepository.DeleteUserReviewByReviewID(ctx, tx, userid, orderid)
	//if err != nil {
	//	usecase.Log.Warn().Msg(err.Error())
	//	return err
	//}

	usecase.ReviewRepository.DeleteUserReviewByReviewID(ctx, tx, userid, reviewid)

	return nil
}

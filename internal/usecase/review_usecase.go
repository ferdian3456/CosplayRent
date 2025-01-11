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

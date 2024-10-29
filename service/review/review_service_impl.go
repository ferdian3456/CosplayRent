package review

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/review"
	reviews "cosplayrent/repository/review"
	"database/sql"
	"github.com/go-playground/validator"
	"time"
)

type ReviewServiceImpl struct {
	ReviewRepository reviews.ReviewRepository
	DB               *sql.DB
	Validate         *validator.Validate
}

func NewReviewService(reviewRepository reviews.ReviewRepository, DB *sql.DB, validate *validator.Validate) ReviewService {
	return &ReviewServiceImpl{
		ReviewRepository: reviewRepository,
		DB:               DB,
		Validate:         validate,
	}
}

func (service *ReviewServiceImpl) Create(ctx context.Context, request review.ReviewCreateRequest) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	reviewDomain := domain.Review{
		User_id:     request.User_id,
		Costume_id:  request.Costume_id,
		Description: request.Description,
		Rating:      request.Rating,
		Created_at:  &now,
	}

	service.ReviewRepository.Create(ctx, tx, reviewDomain)
}

func (service *ReviewServiceImpl) FindByCostumeId(ctx context.Context, id int) []review.ReviewResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	review := []review.ReviewResponse{}
	review, err = service.ReviewRepository.FindByCostumeId(ctx, tx, id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return review
}

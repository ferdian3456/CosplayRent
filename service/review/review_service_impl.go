package review

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/review"
	"cosplayrent/repository/costume"
	reviews "cosplayrent/repository/review"
	"cosplayrent/repository/user"
	"database/sql"
	"github.com/go-playground/validator"
	"log"
	"time"
)

type ReviewServiceImpl struct {
	ReviewRepository  reviews.ReviewRepository
	UserRepository    user.UserRepository
	CostumeRepository costume.CostumeRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewReviewService(reviewRepository reviews.ReviewRepository, costumeRepository costume.CostumeRepository, userRepository user.UserRepository, DB *sql.DB, validate *validator.Validate) ReviewService {
	return &ReviewServiceImpl{
		ReviewRepository:  reviewRepository,
		CostumeRepository: costumeRepository,
		UserRepository:    userRepository,
		DB:                DB,
		Validate:          validate,
	}
}

func (service *ReviewServiceImpl) Create(ctx context.Context, request review.ReviewCreateRequest) {
	log.Printf("User with uuid: %s enter Review Service: Create", request.User_id)

	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()
	reviewDomain := domain.Review{
		User_id:        request.User_id,
		Costume_id:     request.Costume_id,
		Description:    request.Description,
		Review_picture: request.Review_picture,
		Rating:         request.Rating,
		Created_at:     &now,
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

	for i := range review {
		userResult, err := service.UserRepository.FindByUUID(ctx, tx, review[i].User_id)
		helper.PanicIfError(err)
		review[i].Name = userResult.Name
		review[i].Profile_picture = userResult.Profile_picture
	}

	return review
}

func (service *ReviewServiceImpl) FindUserReview(ctx context.Context, uuid string) []review.OwnReviewResponse {
	log.Printf("User with uuid: %s enter Review Service: FindUserReview", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	review := []review.OwnReviewResponse{}
	review, err = service.ReviewRepository.FindUserReview(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	for i := range review {
		costumeResult, err := service.CostumeRepository.FindById(ctx, tx, review[i].Id)
		helper.PanicIfError(err)
		review[i].Costume_name = costumeResult.Name
		review[i].Costume_picture = costumeResult.Picture
		review[i].Ukuran = costumeResult.Ukuran
	}

	return review
}

func (service *ReviewServiceImpl) Update(ctx context.Context, request review.ReviewUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter Review Service: Update", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	now := time.Now()

	updateRequest := review.ReviewUpdateRequest{
		ReviewId:       request.ReviewId,
		Review_picture: request.Review_picture,
		Description:    request.Description,
		Rating:         request.Rating,
		Updated_at:     &now,
	}

	service.ReviewRepository.Update(ctx, tx, updateRequest, uuid)
}

func (service *ReviewServiceImpl) FindUserReviewByReviewID(ctx context.Context, uuid string, reviewid int) review.OwnReviewByReviewID {
	log.Printf("User with uuid: %s enter Review Service: FindUserReviewByReviewID", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	reviewResult := review.OwnReviewByReviewID{}

	reviewResult, err = service.ReviewRepository.FindUserReviewByReviewID(ctx, tx, uuid, reviewid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userResult, err := service.UserRepository.FindByUUID(ctx, tx, reviewResult.User_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeResult, err := service.CostumeRepository.FindById(ctx, tx, reviewResult.Costume_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	reviewResult.Username = userResult.Name
	reviewResult.Costume_name = costumeResult.Name

	return reviewResult
}

func (service *ReviewServiceImpl) DeleteUserReviewByReviewID(ctx context.Context, uuid string, reviewid int) {
	log.Printf("User with uuid: %s enter Review Service: FindUserReviewByReviewID", uuid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	service.ReviewRepository.DeleteUserReviewByReviewID(ctx, tx, uuid, reviewid)
}

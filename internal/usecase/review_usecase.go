package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/review"
	"cosplayrent/internal/repository"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type ReviewUsecase struct {
	ReviewRepository *repository.ReviewRepository
	DB               *sql.DB
	Validate         *validator.Validate
	Log              *zerolog.Logger
}

func NewReviewUsecase(reviewRepository *repository.ReviewRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *ReviewUsecase {
	return &ReviewUsecase{
		ReviewRepository: reviewRepository,
		DB:               DB,
		Validate:         validate,
		Log:              zerolog,
	}
}

func (Usecase *ReviewUsecase) Create(ctx context.Context, request review.ReviewCreateRequest) {
	log.Printf("User with uuid: %s enter Review Usecase: Create", request.User_id)

	err := Usecase.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, request.User_id)
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

	Usecase.ReviewRepository.Create(ctx, tx, reviewDomain)
}

func (Usecase *ReviewUsecase) FindByCostumeId(ctx context.Context, id int) []review.ReviewResponse {
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	review := []review.ReviewResponse{}
	review, err = Usecase.ReviewRepository.FindByCostumeId(ctx, tx, id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range review {
		userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, review[i].User_id)
		helper.PanicIfError(err)
		review[i].Name = userResult.Name
		if userResult.Profile_picture != nil {
			value := imageEnv + *userResult.Profile_picture
			review[i].Profile_picture = &value
		}
	}

	return review
}

func (Usecase *ReviewUsecase) FindUserReview(ctx context.Context, uuid string) []review.OwnReviewResponse {
	log.Printf("User with uuid: %s enter Review Usecase: FindUserReview", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	review := []review.OwnReviewResponse{}
	review, err = Usecase.ReviewRepository.FindUserReview(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range review {
		costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, review[i].Costume_Id)
		helper.PanicIfError(err)
		review[i].Costume_name = costumeResult.Name
		if review[i].Costume_picture != nil {
			value := imageEnv + *review[i].Costume_picture
			review[i].Costume_picture = &value
		}
		//log.Println(review[i].Rating)
		review[i].Ukuran = costumeResult.Ukuran
	}

	return review
}

func (Usecase *ReviewUsecase) Update(ctx context.Context, request review.ReviewUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter Review Usecase: Update", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = Usecase.ReviewRepository.FindUserReviewByReviewID(ctx, tx, uuid, request.ReviewId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()

	updateRequest := review.ReviewUpdateRequest{
		ReviewId:       request.ReviewId,
		Review_picture: request.Review_picture,
		Description:    request.Description,
		Rating:         request.Rating,
		Updated_at:     &now,
	}

	Usecase.ReviewRepository.Update(ctx, tx, updateRequest, uuid)
}

func (Usecase *ReviewUsecase) FindUserReviewByReviewID(ctx context.Context, uuid string, reviewid int) review.OwnReviewByReviewID {
	log.Printf("User with uuid: %s enter Review Usecase: FindUserReviewByReviewID", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	reviewResult := review.OwnReviewByReviewID{}

	reviewResult, err = Usecase.ReviewRepository.FindUserReviewByReviewID(ctx, tx, uuid, reviewid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, reviewResult.User_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, reviewResult.Costume_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	reviewResult.Username = userResult.Name
	reviewResult.Costume_name = costumeResult.Name

	return reviewResult
}

func (Usecase *ReviewUsecase) DeleteUserReviewByReviewID(ctx context.Context, uuid string, reviewid int) {
	log.Printf("User with uuid: %s enter Review Usecase: FindUserReviewByReviewID", uuid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	reviewResult, err := Usecase.ReviewRepository.FindUserReviewByReviewID(ctx, tx, uuid, reviewid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Usecase.ReviewRepository.DeleteUserReviewByReviewID(ctx, tx, uuid, reviewid)

	finalReviewPicturePath := ".." + *reviewResult.Review_picture

	err = os.Remove(finalReviewPicturePath)
	helper.PanicIfError(err)
}

package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/costume"
	"cosplayrent/internal/repository"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/go-playground/validator"
)

type CostumeUsecase struct {
	CostumeRepository *repository.CostumeRepository
	DB                *sql.DB
	Validate          *validator.Validate
	Log               *zerolog.Logger
}

func NewCostumeUsecase(costumeRepository *repository.CostumeRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *CostumeUsecase {
	return &CostumeUsecase{
		CostumeRepository: costumeRepository,
		DB:                DB,
		Validate:          validate,
		Log:               zerolog,
	}
}

func (Usecase *CostumeUsecase) Create(ctx context.Context, request costume.CostumeCreateRequest, userUUID string) {
	log.Printf("User with uuid: %s enter Costume Usecase: Create", request.User_id)

	err := Usecase.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userUUID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	now := time.Now()
	costumeDomain := domain.Costume{
		User_id:     request.User_id,
		Name:        request.Name,
		Description: request.Description,
		Bahan:       request.Bahan,
		Ukuran:      request.Ukuran,
		Berat:       request.Berat,
		Kategori:    request.Kategori,
		Price:       request.Price,
		Available:   request.Available,
		Picture:     request.Picture,
		Created_at:  &now,
	}

	Usecase.CostumeRepository.Create(ctx, tx, costumeDomain)
}

func (Usecase *CostumeUsecase) FindById(ctx context.Context, id int) costume.CostumeResponse {
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume := costume.CostumeResponse{}
	costume, err = Usecase.CostumeRepository.FindById(ctx, tx, id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	user, err := Usecase.UserRepository.FindByUUID(ctx, tx, costume.User_id)
	helper.PanicIfError(err)

	costume.Username = user.Name

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	if user.Profile_picture != nil {
		value := imageEnv + *user.Profile_picture
		costume.Profile_picture = &value
	}

	if costume.Picture != nil {
		value := imageEnv + *costume.Picture
		costume.Picture = &value
	}

	return costume
}

func (Usecase *CostumeUsecase) FindAll(ctx context.Context) []costume.CostumeResponse {
	var err error
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := Usecase.CostumeRepository.FindAll(ctx, tx)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range costume {
		userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, costume[i].User_id)
		helper.PanicIfError(err)
		costume[i].Username = userResult.Name
		if costume[i].Picture != nil {
			value := imageEnv + *costume[i].Picture
			costume[i].Picture = &value
		}
	}
	return costume
}

func (Usecase *CostumeUsecase) FindByName(ctx context.Context, name string) []costume.CostumeResponse {
	var err error
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := Usecase.CostumeRepository.FindByName(ctx, tx, name)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return costume
}

func (Usecase *CostumeUsecase) Update(ctx context.Context, costumeRequest costume.CostumeUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter Costume Usecase: Update", uuid)
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	result, err1 := Usecase.CostumeRepository.FindById(ctx, tx, costumeRequest.Id)
	if err1 != nil {
		panic(exception.NewNotFoundError(err1.Error()))
	}

	now := time.Now()

	updateRequest := costume.CostumeUpdateRequest{
		Id:          result.Id,
		Name:        costumeRequest.Name,
		Description: costumeRequest.Description,
		Bahan:       costumeRequest.Bahan,
		Ukuran:      costumeRequest.Ukuran,
		Berat:       costumeRequest.Berat,
		Kategori:    costumeRequest.Kategori,
		Price:       costumeRequest.Price,
		Picture:     costumeRequest.Picture,
		Available:   costumeRequest.Available,
		Update_at:   &now,
	}

	Usecase.CostumeRepository.Update(ctx, tx, updateRequest, uuid)
}

func (Usecase *CostumeUsecase) Delete(ctx context.Context, id int, uuid string) {
	log.Printf("User with uuid: %s enter Costume Usecase: Delete", uuid)
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeResult, err := Usecase.CostumeRepository.FindById(ctx, tx, id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	finalCostumePicturePath := ".." + *costumeResult.Picture

	err = os.Remove(finalCostumePicturePath)
	helper.PanicIfError(err)

	Usecase.CostumeRepository.Delete(ctx, tx, id, uuid)
}

func (Usecase *CostumeUsecase) FindByUserUUID(ctx context.Context, userUUID string) []costume.CostumeResponse {
	var err error
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := Usecase.CostumeRepository.FindByUserUUID(ctx, tx, userUUID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return costume
}

func (Usecase *CostumeUsecase) FindSellerCostumeByCostumeID(ctx context.Context, userUUID string, costumeID int) costume.CostumeResponse {
	log.Printf("User with uuid: %s enter Costume Usecase: FindSellerCostumeByCostumeID", userUUID)
	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume := costume.CostumeResponse{}
	costume, err = Usecase.CostumeRepository.FindSellerCostumeByCostumeID(ctx, tx, userUUID, costumeID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userResult, err := Usecase.UserRepository.FindByUUID(ctx, tx, costume.User_id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costume.Username = userResult.Name

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	if userResult.Profile_picture != nil {
		value := imageEnv + *userResult.Profile_picture
		costume.Profile_picture = &value
	}

	if costume.Picture != nil {
		value := imageEnv + *costume.Picture
		costume.Picture = &value
	}

	return costume
}

func (Usecase *CostumeUsecase) FindSellerCostume(ctx context.Context, userUUID string) []costume.CostumeResponse {
	log.Printf("User with uuid: %s enter Costume Usecase: FindSellerCostume", userUUID)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := Usecase.CostumeRepository.FindSellerCostume(ctx, tx, userUUID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range costume {
		if costume[i].Picture != nil {
			value := imageEnv + *costume[i].Picture
			costume[i].Picture = &value
		}
	}

	return costume
}

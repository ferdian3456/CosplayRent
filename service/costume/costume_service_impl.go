package costume

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/costume"
	costumes "cosplayrent/repository/costume"
	"cosplayrent/repository/user"
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"
)

type CostumeServiceImpl struct {
	CostumeRepository costumes.CostumeRepository
	UserRepository    user.UserRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewCostumeService(costumeRepository costumes.CostumeRepository, userRepository user.UserRepository, DB *sql.DB, validate *validator.Validate) CostumeService {
	return &CostumeServiceImpl{
		CostumeRepository: costumeRepository,
		UserRepository:    userRepository,
		DB:                DB,
		Validate:          validate,
	}
}

func (service *CostumeServiceImpl) Create(ctx context.Context, request costume.CostumeCreateRequest, userUUID string) {
	log.Printf("User with uuid: %s enter Costume Service: Create", request.User_id)

	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByUUID(ctx, tx, userUUID)
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

	service.CostumeRepository.Create(ctx, tx, costumeDomain)
}

func (service *CostumeServiceImpl) FindById(ctx context.Context, id int) costume.CostumeResponse {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume := costume.CostumeResponse{}
	costume, err = service.CostumeRepository.FindById(ctx, tx, id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	user, err := service.UserRepository.FindByUUID(ctx, tx, costume.User_id)
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

func (service *CostumeServiceImpl) FindAll(ctx context.Context) []costume.CostumeResponse {
	var err error
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := service.CostumeRepository.FindAll(ctx, tx)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i := range costume {
		userResult, err := service.UserRepository.FindByUUID(ctx, tx, costume[i].User_id)
		helper.PanicIfError(err)
		costume[i].Username = userResult.Name
		if costume[i].Picture != nil {
			value := imageEnv + *costume[i].Picture
			costume[i].Picture = &value
		}
	}
	return costume
}

func (service *CostumeServiceImpl) FindByName(ctx context.Context, name string) []costume.CostumeResponse {
	var err error
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := service.CostumeRepository.FindByName(ctx, tx, name)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return costume
}

func (service *CostumeServiceImpl) Update(ctx context.Context, costumeRequest costume.CostumeUpdateRequest, uuid string) {
	log.Printf("User with uuid: %s enter Costume Service: Update", uuid)
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	result, err1 := service.CostumeRepository.FindById(ctx, tx, costumeRequest.Id)
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

	service.CostumeRepository.Update(ctx, tx, updateRequest, uuid)
}

func (service *CostumeServiceImpl) Delete(ctx context.Context, id int, uuid string) {
	log.Printf("User with uuid: %s enter Costume Service: Delete", uuid)
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByUUID(ctx, tx, uuid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeResult, err := service.CostumeRepository.FindById(ctx, tx, id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	finalCostumePicturePath := ".." + *costumeResult.Picture

	err = os.Remove(finalCostumePicturePath)
	helper.PanicIfError(err)

	service.CostumeRepository.Delete(ctx, tx, id, uuid)
}

func (service *CostumeServiceImpl) FindByUserUUID(ctx context.Context, userUUID string) []costume.CostumeResponse {
	var err error
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := service.CostumeRepository.FindByUserUUID(ctx, tx, userUUID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return costume
}

func (service *CostumeServiceImpl) FindSellerCostumeByCostumeID(ctx context.Context, userUUID string, costumeID int) costume.CostumeResponse {
	log.Printf("User with uuid: %s enter Costume Service: FindSellerCostumeByCostumeID", userUUID)
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume := costume.CostumeResponse{}
	costume, err = service.CostumeRepository.FindSellerCostumeByCostumeID(ctx, tx, userUUID, costumeID)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	userResult, err := service.UserRepository.FindByUUID(ctx, tx, costume.User_id)
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

func (service *CostumeServiceImpl) FindSellerCostume(ctx context.Context, userUUID string) []costume.CostumeResponse {
	log.Printf("User with uuid: %s enter Costume Service: FindSellerCostume", userUUID)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	costume, err := service.CostumeRepository.FindSellerCostume(ctx, tx, userUUID)
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

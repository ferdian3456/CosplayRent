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

func (service *CostumeServiceImpl) Create(ctx context.Context, request costume.CostumeCreateRequest) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

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
	costume.Profile_picture = user.Profile_picture

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

	for i := range costume {
		userResult, err := service.UserRepository.FindByUUID(ctx, tx, costume[i].User_id)
		helper.PanicIfError(err)
		costume[i].Username = userResult.Name
		costume[i].Profile_picture = userResult.Profile_picture
	}

	return helper.ToCostumeResponse(costume)
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

func (service *CostumeServiceImpl) Update(ctx context.Context, costumeRequest costume.CostumeUpdateRequest) {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

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

	service.CostumeRepository.Update(ctx, tx, updateRequest)
}

func (service *CostumeServiceImpl) Delete(ctx context.Context, id int) {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	service.CostumeRepository.Delete(ctx, tx, id)
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

	return costume
}

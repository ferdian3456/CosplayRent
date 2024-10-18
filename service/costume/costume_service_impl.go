package costume

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/model/domain"
	"cosplayrent/model/web/costume"
	costumes "cosplayrent/repository/costume"
	"database/sql"
	"github.com/go-playground/validator"
	"time"
)

type CostumeServiceImpl struct {
	CostumeRepository costumes.CostumeRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewCostumeService(costumeRepository costumes.CostumeRepository, DB *sql.DB, validate *validator.Validate) CostumeService {
	return &CostumeServiceImpl{
		CostumeRepository: costumeRepository,
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
		Price:       request.Price,
		Available:   request.Available,
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

	return helper.ToCostumeResponse(costume)
}

func (service *CostumeServiceImpl) Update(ctx context.Context, costumeRequest costume.CostumeUpdateRequest) {
	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	result, err1 := service.CostumeRepository.FindById(ctx, tx, costumeRequest.Id)
	if err1 != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	updateRequest := costume.CostumeUpdateRequest{
		Id:          result.Id,
		Name:        costumeRequest.Name,
		Description: costumeRequest.Description,
		Price:       costumeRequest.Price,
		Picture:     costumeRequest.Picture,
		Available:   costumeRequest.Available,
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

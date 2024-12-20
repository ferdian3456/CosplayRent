package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/costume"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"time"
)

type CostumeUsecase struct {
	CostumeRepository *repository.CostumeRepository
	DB                *sql.DB
	Validate          *validator.Validate
	Log               *zerolog.Logger
	Config            *koanf.Koanf
}

func NewCostumeUsecase(costumeRepository *repository.CostumeRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *CostumeUsecase {
	return &CostumeUsecase{
		CostumeRepository: costumeRepository,
		DB:                DB,
		Validate:          validate,
		Log:               zerolog,
		Config:            koanf,
	}
}

func (usecase *CostumeUsecase) Create(ctx context.Context, userRequest costume.CostumeCreateRequest, uuid string) error {
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

	now := time.Now()

	costumeDomain := domain.Costume{
		User_id:     uuid,
		Name:        userRequest.Name,
		Description: userRequest.Description,
		Bahan:       userRequest.Bahan,
		Ukuran:      userRequest.Ukuran,
		Berat:       userRequest.Berat,
		Kategori:    userRequest.Kategori,
		Price:       userRequest.Price,
		Picture:     *userRequest.Picture,
		Created_at:  &now,
		Updated_at:  &now,
	}

	usecase.CostumeRepository.Create(ctx, tx, costumeDomain)

	return nil
}

func (usecase *CostumeUsecase) Update(ctx context.Context, userRequest costume.CostumeUpdateRequest, uuid string) error {
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

	now := time.Now()

	costumeDomain := domain.Costume{
		Id:          userRequest.Id,
		User_id:     uuid,
		Name:        userRequest.Name,
		Description: userRequest.Description,
		Bahan:       userRequest.Bahan,
		Ukuran:      userRequest.Ukuran,
		Berat:       userRequest.Berat,
		Kategori:    userRequest.Kategori,
		Price:       userRequest.Price,
		Picture:     *userRequest.Picture,
		Updated_at:  &now,
	}

	usecase.CostumeRepository.Update(ctx, tx, costumeDomain)

	return nil
}

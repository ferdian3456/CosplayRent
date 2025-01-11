package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web/category"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type CategoryUsecase struct {
	CategoryRepository *repository.CategoryRepository
	DB                 *sql.DB
	Validate           *validator.Validate
	Log                *zerolog.Logger
	Config             *koanf.Koanf
}

func NewCategoryUsecase(categoryRepository *repository.CategoryRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *CategoryUsecase {
	return &CategoryUsecase{
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
		Log:                zerolog,
		Config:             koanf,
	}
}

func (usecase *CategoryUsecase) FindAllCategory(ctx context.Context) ([]category.CategoryResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	categories, err := usecase.CategoryRepository.FindAllCategory(ctx, tx)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return categories, err
	}

	return categories, nil
}

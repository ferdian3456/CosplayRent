package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/domain"
	"cosplayrent/internal/model/web/costume"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type CostumeUsecase struct {
	UserRepository     *repository.UserRepository
	CostumeRepository  *repository.CostumeRepository
	CategoryRepository *repository.CategoryRepository
	DB                 *sql.DB
	Validate           *validator.Validate
	Log                *zerolog.Logger
	Config             *koanf.Koanf
}

func NewCostumeUsecase(userRepository *repository.UserRepository, costumeRepository *repository.CostumeRepository, categoryRepository *repository.CategoryRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *CostumeUsecase {
	return &CostumeUsecase{
		UserRepository:     userRepository,
		CostumeRepository:  costumeRepository,
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
		Log:                zerolog,
		Config:             koanf,
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
		Available:   userRequest.Available,
		Price:       userRequest.Price,
		Picture:     *userRequest.Picture,
		Updated_at:  &now,
	}

	err = usecase.CostumeRepository.CheckCostume(ctx, tx, uuid, costumeDomain.Id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	usecase.CostumeRepository.Update(ctx, tx, costumeDomain)

	return nil
}

func (usecase *CostumeUsecase) FindAll(ctx context.Context) ([]costume.CostumeResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	costume, err := usecase.CostumeRepository.FindAll(ctx, tx)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range costume {
		userResult, err := usecase.UserRepository.FindByUUID(ctx, tx, costume[i].User_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return costume, err
		}
		costume[i].Username = userResult.Name
		if costume[i].Picture != nil {
			value := imageEnv + *costume[i].Picture
			costume[i].Picture = &value
		}
	}

	return costume, nil
}

func (usecase *CostumeUsecase) FindSellerCostume(ctx context.Context, userUUID string) ([]costume.SellerCostumeResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	costume, err := usecase.CostumeRepository.FindSellerCostume(ctx, tx, userUUID)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range costume {
		if costume[i].Picture != nil {
			value := imageEnv + *costume[i].Picture
			costume[i].Picture = &value
		}
	}

	return costume, nil
}

func (usecase *CostumeUsecase) FindById(ctx context.Context, id int) (costume.CostumeResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	costume := costume.CostumeResponse{}
	costume, err = usecase.CostumeRepository.FindById(ctx, tx, id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	categoryName, err := usecase.CategoryRepository.FindCategoryNameById(ctx, tx, costume.Kategori_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	user, err := usecase.UserRepository.FindNameAndProfile(ctx, tx, costume.User_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	costume.Username = user.Name

	imageEnv := usecase.Config.String("IMAGE_ENV")

	if user.Profile_picture != nil {
		value := imageEnv + *user.Profile_picture
		costume.Profile_picture = &value
	}

	if costume.Picture != nil {
		value := imageEnv + *costume.Picture
		costume.Picture = &value
	}

	costume.Kategori = categoryName

	return costume, nil
}

func (usecase *CostumeUsecase) FindSellerCostumeByCostumeID(ctx context.Context, userUUID string, costumeID int) (costume.CostumeResponse, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	costume := costume.CostumeResponse{}
	costume, err = usecase.CostumeRepository.FindSellerCostumeByCostumeID(ctx, tx, userUUID, costumeID)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	userResult, err := usecase.UserRepository.FindNameAndProfile(ctx, tx, costume.User_id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return costume, err
	}

	costume.Username = userResult.Name

	imageEnv := usecase.Config.String("IMAGE_ENV")

	if userResult.Profile_picture != nil {
		value := imageEnv + *userResult.Profile_picture
		costume.Profile_picture = &value
	}

	if costume.Picture != nil {
		value := imageEnv + *costume.Picture
		costume.Picture = &value
	}

	return costume, nil
}

func (usecase *CostumeUsecase) Delete(ctx context.Context, id int, uuid string) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.CostumeRepository.CheckDeleteCostume(ctx, tx, uuid, id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	costumepicture, err := usecase.CostumeRepository.FindPictureById(ctx, tx, id)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	usecase.CostumeRepository.Delete(ctx, tx, id, uuid)

	finalCostumePicturePath := ".." + *costumepicture

	err = os.Remove(finalCostumePicturePath)
	if err != nil {
		respErr := errors.New("failed to remove costume")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	return nil
}

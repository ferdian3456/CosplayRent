package usecase

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web/wishlist"
	"cosplayrent/internal/repository"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type WishlistUsecase struct {
	WishlistRepository *repository.WishlistRepository
	CostumeRepository  *repository.CostumeRepository
	DB                 *sql.DB
	Validate           *validator.Validate
	Log                *zerolog.Logger
	Config             *koanf.Koanf
}

func NewWishlistUsecase(wishlistRepository *repository.WishlistRepository, costumeRepository *repository.CostumeRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger, koanf *koanf.Koanf) *WishlistUsecase {
	return &WishlistUsecase{
		WishlistRepository: wishlistRepository,
		CostumeRepository:  costumeRepository,
		DB:                 DB,
		Validate:           validate,
		Log:                zerolog,
		Config:             koanf,
	}
}

func (usecase *WishlistUsecase) AddWishlist(ctx context.Context, uuid string, costumeid int) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.CostumeRepository.CheckOwnership(ctx, tx, uuid, costumeid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	usecase.WishlistRepository.AddWishlist(ctx, tx, uuid, costumeid)
	return nil
}

func (usecase *WishlistUsecase) DeleteWishlist(ctx context.Context, uuid string, costumeid int) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.CostumeRepository.CheckOwnership(ctx, tx, uuid, costumeid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}

	usecase.WishlistRepository.DeleteWishlist(ctx, tx, uuid, costumeid)
	return nil
}

func (usecase *WishlistUsecase) CheckWishlistStatus(ctx context.Context, uuid string, costumeid int) error {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	err = usecase.WishlistRepository.FindWishlistById(ctx, tx, uuid, costumeid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return err
	}
	return nil
}

func (usecase *WishlistUsecase) FindAllWishListByUserId(ctx context.Context, uuid string) ([]wishlist.WishListResponses, error) {
	tx, err := usecase.DB.Begin()
	if err != nil {
		respErr := errors.New("failed to start transaction")
		usecase.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer helper.CommitOrRollback(tx)

	wishlistResponse, err := usecase.WishlistRepository.FindCostumeIdByOrderId(ctx, tx, uuid)
	if err != nil {
		usecase.Log.Warn().Msg(err.Error())
		return wishlistResponse, err
	}

	imageEnv := usecase.Config.String("IMAGE_ENV")

	for i := range wishlistResponse {
		costumeResult, err := usecase.CostumeRepository.FindById(ctx, tx, wishlistResponse[i].Costume_id)
		if err != nil {
			usecase.Log.Warn().Msg(err.Error())
			return wishlistResponse, err
		}

		wishlistResponse[i].Costume_name = costumeResult.Name
		wishlistResponse[i].Costume_price = costumeResult.Price
		wishlistResponse[i].Costume_size = *costumeResult.Ukuran
		if costumeResult.Picture != nil {
			value := imageEnv + *costumeResult.Picture
			wishlistResponse[i].Costume_picture = &value
		}
	}

	return wishlistResponse, nil
}

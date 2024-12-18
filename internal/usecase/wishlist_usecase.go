package usecase

import (
	"context"
	"cosplayrent/internal/exception"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web/wishlist"
	"cosplayrent/internal/repository"
	"database/sql"
	"log"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type WishlistUsecase struct {
	WishlistRepository *repository.WishlistRepository
	DB                 *sql.DB
	Validate           *validator.Validate
	Log                *zerolog.Logger
}

func NewWishlistUsecase(wishlistRepository *repository.WishlistRepository, DB *sql.DB, validate *validator.Validate, zerolog *zerolog.Logger) *WishlistUsecase {
	return &WishlistUsecase{
		WishlistRepository: wishlistRepository,
		DB:                 DB,
		Validate:           validate,
		Log:                zerolog,
	}
}

func (Usecase *WishlistUsecase) AddWishList(ctx context.Context, costumeid int, userid string) {
	log.Printf("User with uuid: %s enter Wishlist Usecase: AddWishlist", userid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = Usecase.CostumeRepository.FindById(ctx, tx, costumeid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Usecase.WishlistRepository.AddWishList(ctx, tx, costumeid, userid)
}

func (Usecase *WishlistUsecase) DeleteWishList(ctx context.Context, costumeid int, userid string) {
	log.Printf("User with uuid: %s enter Wishlist Usecase: DeleteWishList", userid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = Usecase.CostumeRepository.FindById(ctx, tx, costumeid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	Usecase.WishlistRepository.DeleteWishList(ctx, tx, costumeid, userid)
}

func (Usecase *WishlistUsecase) FindAllWishlistByUserId(ctx context.Context, userid string) []wishlist.WishListResponses {
	log.Printf("User with uuid: %s enter Wishlist Usecase: FindAllWishlistByUserId", userid)

	tx, err := Usecase.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	defer helper.CommitOrRollback(tx)

	_, err = Usecase.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = Usecase.WishlistRepository.CheckWishlistUserByUserId(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeIds, err := Usecase.WishlistRepository.FindUserWishListByUserId(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	WishListResponses := make([]wishlist.WishListResponses, len(costumeIds))

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i, costumeId := range costumeIds {
		costumeDetail, err := Usecase.CostumeRepository.FindById(ctx, tx, costumeId)
		helper.PanicIfError(err)

		response := wishlist.WishListResponses{
			CostumeId:     costumeDetail.Id,
			Costume_name:  costumeDetail.Name,
			Costume_price: costumeDetail.Price,
			Costume_size:  *costumeDetail.Ukuran,
		}

		if costumeDetail.Picture != nil {
			value := imageEnv + *costumeDetail.Picture
			response.Costume_picture = &value
		}

		WishListResponses[i] = response
	}

	return WishListResponses
}

package wishlist

import (
	"context"
	"cosplayrent/exception"
	"cosplayrent/helper"
	wishlists "cosplayrent/model/web/wishlist"
	"cosplayrent/repository/costume"
	"cosplayrent/repository/user"
	"cosplayrent/repository/wishlist"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type WishlistServiceImpl struct {
	WishlistRepository wishlist.WishListRepository
	UserRepository     user.UserRepository
	CostumeRepository  costume.CostumeRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewWishlistService(wishlistRepository wishlist.WishListRepository, userRepository user.UserRepository, costumeRepository costume.CostumeRepository, DB *sql.DB, validate *validator.Validate) WishlistService {
	return &WishlistServiceImpl{
		WishlistRepository: wishlistRepository,
		UserRepository:     userRepository,
		CostumeRepository:  costumeRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *WishlistServiceImpl) AddWishList(ctx context.Context, costumeid int, userid string) {
	log.Printf("User with uuid: %s enter Wishlist Service: AddWishlist", userid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = service.CostumeRepository.FindById(ctx, tx, costumeid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.WishlistRepository.AddWishList(ctx, tx, costumeid, userid)
}

func (service *WishlistServiceImpl) DeleteWishList(ctx context.Context, costumeid int, userid string) {
	log.Printf("User with uuid: %s enter Wishlist Service: DeleteWishList", userid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = service.CostumeRepository.FindById(ctx, tx, costumeid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.WishlistRepository.DeleteWishList(ctx, tx, costumeid, userid)
}

func (service *WishlistServiceImpl) FindAllWishlistByUserId(ctx context.Context, userid string) []wishlists.WishListResponses {
	log.Printf("User with uuid: %s enter Wishlist Service: FindAllWishlistByUserId", userid)

	tx, err := service.DB.Begin()
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	defer helper.CommitOrRollback(tx)

	_, err = service.UserRepository.FindByUUID(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	_, err = service.WishlistRepository.CheckWishlistUserByUserId(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	costumeIds, err := service.WishlistRepository.FindUserWishListByUserId(ctx, tx, userid)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	
	WishListResponses := make([]wishlists.WishListResponses, len(costumeIds))

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	for i, costumeId := range costumeIds {
		costumeDetail, err := service.CostumeRepository.FindById(ctx, tx, costumeId)
		helper.PanicIfError(err)

		response := wishlists.WishListResponses{
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

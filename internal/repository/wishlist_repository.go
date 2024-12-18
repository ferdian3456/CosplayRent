package repository

import (
	"context"
	"cosplayrent/internal/helper"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
	"log"
)

type WishlistRepository struct {
	Log *zerolog.Logger
}

func NewWishlistRepository(zerolog *zerolog.Logger) *WishlistRepository {
	return &WishlistRepository{
		Log: zerolog,
	}
}

func (repository *WishlistRepository) AddWishList(ctx context.Context, tx *sql.Tx, costumeid int, userid string) {
	log.Printf("User with uuid: %s enter Wishlist Repository: AddWishList", userid)

	query := "INSERT INTO wishlists (user_id, costume_id) VALUES($1,$2)"
	_, err := tx.ExecContext(ctx, query, userid, costumeid)
	helper.PanicIfError(err)
}

func (repository *WishlistRepository) DeleteWishList(ctx context.Context, tx *sql.Tx, costumeid int, userid string) {
	log.Printf("User with uuid: %s enter Wishlist Repository: DeleteWishList", userid)

	query := "DELETE FROM wishlists wishlists WHERE user_id=$1 AND costume_id=$1 "
	_, err := tx.ExecContext(ctx, query, userid, costumeid)
	helper.PanicIfError(err)
}

func (repository *WishlistRepository) FindUserWishListByUserId(ctx context.Context, tx *sql.Tx, userid string) ([]int, error) {
	log.Printf("User with uuid: %s enter Wishlist Repository: FindUserWishListByUserId", userid)

	query := "SELECT costume_id FROM wishlists WHERE user_id=$1"
	rows, err := tx.QueryContext(ctx, query, userid)
	helper.PanicIfError(err)

	defer rows.Close()

	var costumesId []int
	for rows.Next() {
		var costumeId int
		err = rows.Scan(&costumeId)
		helper.PanicIfError(err)
		costumesId = append(costumesId, costumeId)
	}

	if len(costumesId) == 0 {
		return nil, errors.New("costume is not found")
	}

	return costumesId, nil
}

func (repository *WishlistRepository) CheckWishlistUserByUserId(ctx context.Context, tx *sql.Tx, userid string) ([]int, error) {
	log.Printf("User with uuid: %s enter Wishlist Repository: CheckWishlistUserByUserId", userid)

	query := "SELECT id FROM wishlists WHERE user_id=$1"
	rows, err := tx.QueryContext(ctx, query, userid)
	helper.PanicIfError(err)
	hasData := false

	defer rows.Close()

	var wishlistsId []int

	for rows.Next() {
		var wishlistid int
		err = rows.Scan(&wishlistid)
		helper.PanicIfError(err)
		wishlistsId = append(wishlistsId, wishlistid)
		hasData = true
	}

	if hasData == false {
		return wishlistsId, errors.New("wishlist is not found")
	}

	return wishlistsId, nil
}

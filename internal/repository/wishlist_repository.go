package repository

import (
	"context"
	"cosplayrent/internal/model/web/wishlist"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
)

type WishlistRepository struct {
	Log *zerolog.Logger
}

func NewWishlistRepository(zerolog *zerolog.Logger) *WishlistRepository {
	return &WishlistRepository{
		Log: zerolog,
	}
}

func (repository *WishlistRepository) AddWishlist(ctx context.Context, tx *sql.Tx, uuid string, costumeid int) {
	query := "INSERT INTO wishlists (user_id,costume_id) VALUES ($1,$2)"
	_, err := tx.ExecContext(ctx, query, uuid, costumeid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *WishlistRepository) FindCostumeIdByOrderId(ctx context.Context, tx *sql.Tx, uuid string) ([]wishlist.WishListResponses, error) {
	query := "SELECT costume_id FROM wishlists WHERE user_id=$1"
	rows, err := tx.QueryContext(ctx, query, uuid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	defer rows.Close()

	hasData := false

	wishlists := []wishlist.WishListResponses{}
	for rows.Next() {
		wishlist := wishlist.WishListResponses{}
		err = rows.Scan(&wishlist.Costume_id)
		if err != nil {
			respErr := errors.New("failed to scan query result")
			repository.Log.Panic().Err(err).Msg(respErr.Error())
		}

		wishlists = append(wishlists, wishlist)
		hasData = true
	}

	if hasData == false {
		return wishlists, errors.New("user not found")
	}

	return wishlists, nil
}

func (repository *WishlistRepository) DeleteWishlist(ctx context.Context, tx *sql.Tx, uuid string, costumeid int) {
	query := "DELETE FROM wishlists WHERE user_id=$1 AND costume_id=$2"
	_, err := tx.ExecContext(ctx, query, uuid, costumeid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}
}

func (repository *WishlistRepository) FindWishlistById(ctx context.Context, tx *sql.Tx, uuid string, costumeid int) error {
	query := "SELECT id FROM wishlists WHERE user_id=$1 AND costume_id=$2"
	row, err := tx.QueryContext(ctx, query, uuid, costumeid)
	if err != nil {
		respErr := errors.New("failed to query into database")
		repository.Log.Panic().Err(err).Msg(respErr.Error())
	}

	if row.Next() {
		return nil
	} else {
		return errors.New("wishlist not found")
	}
}

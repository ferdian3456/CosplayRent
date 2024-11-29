package wishlist

import (
	"context"
	"database/sql"
)

type WishListRepository interface {
	AddWishList(ctx context.Context, tx *sql.Tx, costumeid int, userid string)
	DeleteWishList(ctx context.Context, tx *sql.Tx, costumeid int, userid string)
	FindUserWishListByUserId(ctx context.Context, tx *sql.Tx, userid string) ([]int, error)
	CheckWishlistUserByUserId(ctx context.Context, tx *sql.Tx, userid string) ([]int, error)
}

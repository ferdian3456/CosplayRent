package wishlist

import (
	"context"
	"cosplayrent/model/web/wishlist"
)

type WishlistService interface {
	AddWishList(ctx context.Context, costumeid int, userid string)
	DeleteWishList(ctx context.Context, costumeid int, userid string)
	FindAllWishlistByUserId(ctx context.Context, userid string) []wishlist.WishListResponses
}

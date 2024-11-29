package wishlist

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type WishlistController interface {
	AddWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAllWishListByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

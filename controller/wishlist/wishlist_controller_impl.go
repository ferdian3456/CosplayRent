package wishlist

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/service/wishlist"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type WishlistControllerImpl struct {
	WishlistService wishlist.WishlistService
}

func NewWishlistController(wishlistService wishlist.WishlistService) WishlistController {
	return &WishlistControllerImpl{
		WishlistService: wishlistService,
	}
}

func (controller WishlistControllerImpl) AddWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Wishlist Controller: AddWishlist", userUUID)

	costumeID := params.ByName(`costumeID`)
	fixCostumeID, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	controller.WishlistService.AddWishList(request.Context(), fixCostumeID, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller WishlistControllerImpl) DeleteWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Wishlist Controller: DeleteWishList", userUUID)

	costumeID := params.ByName("costumeID")
	fixCostumeID, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	controller.WishlistService.DeleteWishList(request.Context(), fixCostumeID, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller WishlistControllerImpl) FindAllWishListByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Wishlist Controller: FindAllWishListByUserId", userUUID)

	wishlistResult := controller.WishlistService.FindAllWishlistByUserId(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   wishlistResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

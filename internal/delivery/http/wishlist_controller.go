package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type WishlistController struct {
	WishlistUsecase *usecase.WishlistUsecase
	Log             *zerolog.Logger
}

func NewWishlistController(WishlistUsecase *usecase.WishlistUsecase, zerolog *zerolog.Logger) *WishlistController {
	return &WishlistController{
		WishlistUsecase: WishlistUsecase,
		Log:             zerolog,
	}
}

func (controller WishlistController) AddWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	controller.WishlistUsecase.AddWishList(request.Context(), fixCostumeID, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller WishlistController) DeleteWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	controller.WishlistUsecase.DeleteWishList(request.Context(), fixCostumeID, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller WishlistController) FindAllWishListByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	wishlistResult := controller.WishlistUsecase.FindAllWishlistByUserId(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   wishlistResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

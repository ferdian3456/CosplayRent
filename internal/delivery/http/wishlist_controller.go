package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/usecase"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

type WishlistController struct {
	WishlistUsecase *usecase.WishlistUsecase
	Log             *zerolog.Logger
}

func NewWishlistController(wishlistUsecase *usecase.WishlistUsecase, zerolog *zerolog.Logger) *WishlistController {
	return &WishlistController{
		WishlistUsecase: wishlistUsecase,
		Log:             zerolog,
	}
}

func (controller WishlistController) AddWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeId := params.ByName("costumeID")
	fixCostumeId, err := strconv.Atoi(costumeId)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	err = controller.WishlistUsecase.AddWishlist(request.Context(), userUUID, fixCostumeId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller WishlistController) DeleteWishlist(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeId := params.ByName("costumeID")
	fixCostumeId, err := strconv.Atoi(costumeId)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	err = controller.WishlistUsecase.DeleteWishlist(request.Context(), userUUID, fixCostumeId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller WishlistController) CheckWishlistStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeId := params.ByName("costumeID")
	fixCostumeId, err := strconv.Atoi(costumeId)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	err = controller.WishlistUsecase.CheckWishlistStatus(request.Context(), userUUID, fixCostumeId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	wishlistResponse := web.WishlistStatusResponse{
		Status_wishlist: "True",
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   wishlistResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *WishlistController) FindAllWishListByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	wishlistResponse, err := controller.WishlistUsecase.FindAllWishListByUserId(request.Context(), userUUID)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   wishlistResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

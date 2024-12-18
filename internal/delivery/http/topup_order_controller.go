package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/usecase"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type TopUpOrderController struct {
	TopUpOrderUsecase *usecase.TopUpOrderUsecase
	Log               *zerolog.Logger
}

func NewTopUpOrderController(TopUpOrderUsecase *usecase.TopUpOrderUsecase, zerolog *zerolog.Logger) *TopUpOrderController {
	return &TopUpOrderController{
		TopUpOrderUsecase: TopUpOrderUsecase,
		Log:               zerolog,
	}
}

func (controller TopUpOrderController) CreateTopUpOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: CreateTopUpOrder", userUUID)

	topUpEMoneyRequest := user.TopUpEmoney{}
	helper.ReadFromRequestBody(request, &topUpEMoneyRequest)

	midtransResult := controller.TopUpOrderUsecase.CreateTopUpOrder(request.Context(), topUpEMoneyRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   midtransResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller TopUpOrderController) CheckTopUpOrderByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderID := params.ByName("orderID")

	topuporderResult := controller.TopUpOrderUsecase.CheckTopUpOrderByOrderId(request.Context(), orderID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   topuporderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

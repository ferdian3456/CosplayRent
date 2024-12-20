package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type TopUpOrderController struct {
	TopUpOrderUsecase *usecase.TopUpOrderUsecase
	Log               *zerolog.Logger
}

func NewTopUpOrderController(topUpOrderUsecase *usecase.TopUpOrderUsecase, zerolog *zerolog.Logger) *TopUpOrderController {
	return &TopUpOrderController{
		TopUpOrderUsecase: topUpOrderUsecase,
		Log:               zerolog,
	}
}

func (controller TopUpOrderController) CreateTopUpOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	userRequest := user.TopUpEmoney{}
	helper.ReadFromRequestBody(request, &userRequest)

	midtransResponse, err := controller.TopUpOrderUsecase.CreateTopUpOrder(request.Context(), userRequest, userUUID)
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
		Data:   midtransResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

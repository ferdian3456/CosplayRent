package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type OrderController struct {
	OrderUsecase *usecase.OrderUsecase
	Log          *zerolog.Logger
}

func NewOrderController(orderUsecase *usecase.OrderUsecase, zerolog *zerolog.Logger) *OrderController {
	return &OrderController{
		OrderUsecase: orderUsecase,
		Log:          zerolog,
	}
}

func (controller OrderController) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	orderCreateRequest := order.OrderCreateRequest{}
	helper.ReadFromRequestBody(request, &orderCreateRequest)

	err := controller.OrderUsecase.Create(request.Context(), orderCreateRequest, userUUID)
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

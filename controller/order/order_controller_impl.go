package order

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/order"
	orders "cosplayrent/service/order"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OrderControllerImpl struct {
	OrderService orders.OrderService
}

func NewOrderController(orderService orders.OrderService) OrderController {
	return &OrderControllerImpl{
		OrderService: orderService,
	}
}

func (controller OrderControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderCreateRequest := order.OrderCreateRequest{}
	helper.ReadFromRequestBody(request, &orderCreateRequest)
	controller.OrderService.Create(request.Context(), orderCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")

	orderDomain := controller.OrderService.FindByUserId(request.Context(), userID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   orderDomain,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

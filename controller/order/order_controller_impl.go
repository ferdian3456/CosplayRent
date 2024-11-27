package order

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/order"
	orders "cosplayrent/service/order"
	"github.com/julienschmidt/httprouter"
	"log"
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

	log.Printf("User with uuid: %s enter User Controller: TopUp", userUUID)

	orderCreateRequest := order.OrderCreateRequest{}
	helper.ReadFromRequestBody(request, &orderCreateRequest)

	orderCreateRequest.Id = userUUID
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

func (controller OrderControllerImpl) DirectlyOrderToMidtrans(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Order Controller: DirectlyOrderToMidtrans", userUUID)

	directOrderToMidtransRequest := order.DirectlyOrderToMidtrans{}
	helper.ReadFromRequestBody(request, &directOrderToMidtransRequest)

	directOrderToMidtransRequest.Costumer_id = userUUID
	midtransResult := controller.OrderService.DirectlyOrderToMidtrans(request.Context(), userUUID, directOrderToMidtransRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   midtransResult,
	}

	helper.WriteToResponseBody(writer, webResponse)

}

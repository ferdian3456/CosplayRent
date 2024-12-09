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

func (controller OrderControllerImpl) CheckStatusPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderId := params.ByName("orderID")

	orderResult := controller.OrderService.FindOrderDetailByOrderId(request.Context(), orderId)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   orderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) GetAllSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Order Controller: GetAllSellerOrder", userUUID)

	sellerOrderResult := controller.OrderService.GetAllSellerOrder(request.Context(), userUUID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   sellerOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) GetAllUserOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Order Controller: GetAllUserOrder", userUUID)

	userOrderResult := controller.OrderService.GetAllUserOrder(request.Context(), userUUID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) UpdateSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	orderId := params.ByName("orderID")

	updateRequest := order.OrderUpdateRequest{}
	helper.ReadFromRequestBody(request, &updateRequest)

	log.Printf("User with uuid: %s enter Order Controller: UpdateSellerOrder", userUUID)

	controller.OrderService.UpdateSellerOrder(request.Context(), updateRequest, userUUID, orderId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) GetDetailOrderByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	orderId := params.ByName("orderID")

	log.Printf("User with uuid: %s enter Order Controller: GetDetailOrderByOrderId", userUUID)

	detailOrderResult := controller.OrderService.GetDetailOrderByOrderId(request.Context(), userUUID, orderId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   detailOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) GetUserDetailOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	orderId := params.ByName("orderID")

	log.Printf("User with uuid: %s enter Order Controller: GetUserOrder", userUUID)

	userOrderResult := controller.OrderService.GetUserDetailOrder(request.Context(), userUUID, orderId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderControllerImpl) CheckBalanceWithOrderAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//bodyBytes, _ := ioutil.ReadAll(request.Body)
	//log.Println("Raw request body:", string(bodyBytes))

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

	log.Printf("User with uuid: %s enter Order Controller: CheckBalanceWithOrderAmount", userUUID)

	checkBalanceWithOrderAmountRequest := order.CheckBalanceWithOrderAmount{}
	helper.ReadFromRequestBody(request, &checkBalanceWithOrderAmountRequest)

	log.Println(checkBalanceWithOrderAmountRequest.Order_amount)

	userOrderResult := controller.OrderService.CheckBalanceWithOrderAmount(request.Context(), checkBalanceWithOrderAmountRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

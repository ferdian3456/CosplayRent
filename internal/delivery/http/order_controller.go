package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/order"
	"cosplayrent/internal/usecase"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type OrderController struct {
	OrderUsecase *usecase.OrderUsecase
	Log          *zerolog.Logger
}

func NewOrderController(OrderUsecase *usecase.OrderUsecase, zerolog *zerolog.Logger) *OrderController {
	return &OrderController{
		OrderUsecase: OrderUsecase,
		Log:          zerolog,
	}
}

func (controller OrderController) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	controller.OrderUsecase.Create(request.Context(), orderCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID := params.ByName("userID")

	orderDomain := controller.OrderUsecase.FindByUserId(request.Context(), userID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   orderDomain,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) DirectlyOrderToMidtrans(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	midtransResult := controller.OrderUsecase.DirectlyOrderToMidtrans(request.Context(), userUUID, directOrderToMidtransRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   midtransResult,
	}

	helper.WriteToResponseBody(writer, webResponse)

}

func (controller OrderController) CheckStatusPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderId := params.ByName("orderID")

	orderResult := controller.OrderUsecase.FindOrderDetailByOrderId(request.Context(), orderId)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   orderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetAllSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	sellerOrderResult := controller.OrderUsecase.GetAllSellerOrder(request.Context(), userUUID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   sellerOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetAllUserOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	userOrderResult := controller.OrderUsecase.GetAllUserOrder(request.Context(), userUUID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) UpdateSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	controller.OrderUsecase.UpdateSellerOrder(request.Context(), updateRequest, userUUID, orderId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetDetailOrderByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	detailOrderResult := controller.OrderUsecase.GetDetailOrderByOrderId(request.Context(), userUUID, orderId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   detailOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetUserDetailOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	userOrderResult := controller.OrderUsecase.GetUserDetailOrder(request.Context(), userUUID, orderId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) CheckBalanceWithOrderAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	userOrderResult := controller.OrderUsecase.CheckBalanceWithOrderAmount(request.Context(), checkBalanceWithOrderAmountRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

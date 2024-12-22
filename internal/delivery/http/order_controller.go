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

	directOrderToMidtransRequest := order.DirectlyOrderToMidtrans{}
	helper.ReadFromRequestBody(request, &directOrderToMidtransRequest)

	directOrderToMidtransRequest.Costumer_id = userUUID
	midtransResult, err := controller.OrderUsecase.Create(request.Context(), userUUID, directOrderToMidtransRequest)
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
		Data:   midtransResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) CheckStatusPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderId := params.ByName("orderID")

	orderStatus, err := controller.OrderUsecase.CheckStatusPayment(request.Context(), orderId)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	orderResponse := web.OrderStatusResponse{
		Status_payment: orderStatus,
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   orderResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetAllSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	sellerOrderResult, err := controller.OrderUsecase.GetAllSellerOrder(request.Context(), userUUID)
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
		Data:   sellerOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetDetailOrderByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)
	orderId := params.ByName("orderID")

	detailOrderResult, err := controller.OrderUsecase.GetDetailOrderByOrderId(request.Context(), userUUID, orderId)
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
		Data:   detailOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetUserDetailOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	orderId := params.ByName("orderID")

	userOrderResult, err := controller.OrderUsecase.GetUserDetailOrder(request.Context(), userUUID, orderId)
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
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) GetAllUserOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	userOrderResult, err := controller.OrderUsecase.GetAllUserOrder(request.Context(), userUUID)
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
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller OrderController) CheckBalanceWithOrderAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//bodyBytes, _ := ioutil.ReadAll(request.Body)
	//log.Println("Raw request body:", string(bodyBytes))

	userUUID, _ := request.Context().Value("user_uuid").(string)

	checkBalanceWithOrderAmountRequest := order.CheckBalanceWithOrderAmount{}
	helper.ReadFromRequestBody(request, &checkBalanceWithOrderAmountRequest)

	userOrderResult, err := controller.OrderUsecase.CheckBalanceWithOrderAmount(request.Context(), checkBalanceWithOrderAmountRequest, userUUID)
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
		Data:   userOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

package topup_order

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/user"
	"cosplayrent/service/topup_order"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type TopUpOrderControllerImpl struct {
	TopUpOrderService topup_order.TopupOrderService
}

func NewTopUpOrderController(topuporderService topup_order.TopupOrderService) *TopUpOrderControllerImpl {
	return &TopUpOrderControllerImpl{
		TopUpOrderService: topuporderService,
	}
}

func (controller TopUpOrderControllerImpl) CreateTopUpOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	midtransResult := controller.TopUpOrderService.CreateTopUpOrder(request.Context(), topUpEMoneyRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   midtransResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller TopUpOrderControllerImpl) CheckTopUpOrderByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderID := params.ByName("orderID")

	topuporderResult := controller.TopUpOrderService.CheckTopUpOrderByOrderId(request.Context(), orderID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   topuporderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

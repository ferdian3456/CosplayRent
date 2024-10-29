package midtrans

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/midtrans"
	midtranss "cosplayrent/service/midtrans"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type MidtransControllerImpl struct {
	MidtransService midtranss.MidtransService
}

func NewMidtransController(midtransService midtranss.MidtransService) MidtransController {
	return &MidtransControllerImpl{
		MidtransService: midtransService,
	}
}

func (controller MidtransControllerImpl) CreateTransaction(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderID := params.ByName("orderID")

	midtransCreateRequest := midtrans.MidtransRequest{}
	helper.ReadFromRequestBody(request, &midtransCreateRequest)
	midtransResponse := controller.MidtransService.CreateTransaction(request.Context(), midtransCreateRequest, orderID)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   midtransResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller MidtransControllerImpl) MidtransCallBack(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	midtransCallBack := midtrans.MidtransCallback{}
	helper.ReadFromRequestBody(request, &midtransCallBack)
	//log.Println(midtransCallBack.OrderID)

	fmt.Println("Transaction Success")
	fmt.Println(midtransCallBack.Status_Code)
	fmt.Println(midtransCallBack.TransactionStatus)
	fmt.Println(midtransCallBack.OrderID)
	fmt.Println(midtransCallBack.GrossAmount)
	fmt.Println(midtransCallBack.PaymentType)
	fmt.Println(midtransCallBack.TransactionTime)
	fmt.Println(midtransCallBack.TransactionID)
	fmt.Println(midtransCallBack.SignatureKey)

	controller.MidtransService.MidtransCallBack(request.Context(), midtransCallBack.OrderID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

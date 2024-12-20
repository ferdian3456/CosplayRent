package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	midtransWeb "cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type MidtransController struct {
	MidtransUsecase *usecase.MidtransUsecase
	Log             *zerolog.Logger
}

func NewMidtransController(midtransUsecase *usecase.MidtransUsecase, zerolog *zerolog.Logger) *MidtransController {
	return &MidtransController{
		MidtransUsecase: midtransUsecase,
		Log:             zerolog,
	}
}

func (controller MidtransController) MidtransCallBack(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	midtransCallBack := midtransWeb.MidtransCallback{}
	helper.ReadFromRequestBody(request, &midtransCallBack)
	//log.Println(midtransCallBack.OrderID)
	
	if midtransCallBack.Status_Code == "200" {
		//fmt.Println("Transaction Success")
		//fmt.Println(midtransCallBack.Status_Code)
		//fmt.Println(midtransCallBack.TransactionStatus)
		//fmt.Println(midtransCallBack.OrderID)
		//fmt.Println(midtransCallBack.GrossAmount)
		//fmt.Println(midtransCallBack.PaymentType)
		//fmt.Println(midtransCallBack.TransactionTime)
		//fmt.Println(midtransCallBack.TransactionID)
		//fmt.Println(midtransCallBack.SignatureKey)

		controller.Log.Debug().Msg("OrderId:" + midtransCallBack.OrderID)

		controller.MidtransUsecase.MidtransCallBack(request.Context(), midtransCallBack)

		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

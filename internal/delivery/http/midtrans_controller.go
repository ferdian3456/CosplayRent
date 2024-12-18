package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/midtrans"
	"cosplayrent/internal/usecase"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type MidtransController struct {
	MidtransUsecase *usecase.MidtransUsecase
	Log             *zerolog.Logger
}

func NewMidtransController(MidtransUsecase *usecase.MidtransUsecase, zerolog *zerolog.Logger) *MidtransController {
	return &MidtransController{
		MidtransUsecase: MidtransUsecase,
		Log:             zerolog,
	}
}

func (controller MidtransController) CreateTransaction(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Review Controller: FindUserReview", userUUID)

	midtransCreateRequest := midtrans.MidtransRequest{}
	helper.ReadFromRequestBody(request, &midtransCreateRequest)
	//midtransResponse := controller.MidtransUsecase.CreateTransaction(request.Context(), midtransCreateRequest, userUUID)
	//webResponse := web.WebResponse{
	//	Code:   200,
	//	Status: "OK",
	//	Data:   midtransResponse,
	//}
	//
	//helper.WriteToResponseBody(writer, webResponse)
}

func (controller MidtransController) MidtransCallBack(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	midtransCallBack := midtrans.MidtransCallback{}
	helper.ReadFromRequestBody(request, &midtransCallBack)
	//log.Println(midtransCallBack.OrderID)

	if midtransCallBack.TransactionStatus == "settlement" {
		//fmt.Println("Transaction Success")
		//fmt.Println(midtransCallBack.Status_Code)
		//fmt.Println(midtransCallBack.TransactionStatus)
		//fmt.Println(midtransCallBack.OrderID)
		//fmt.Println(midtransCallBack.GrossAmount)
		//fmt.Println(midtransCallBack.PaymentType)
		//fmt.Println(midtransCallBack.TransactionTime)
		//fmt.Println(midtransCallBack.TransactionID)
		//fmt.Println(midtransCallBack.SignatureKey)

		controller.MidtransUsecase.MidtransCallBack(request.Context(), midtransCallBack.OrderID, midtransCallBack.GrossAmount)

		webResponse := web.WebResponse{
			Code:   200,
			Status: "OK",
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}

package order

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type OrderController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DirectlyOrderToMidtrans(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckStatusPayment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateSellerOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetDetailOrderByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetUserDetailOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllUserOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckBalanceWithOrderAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

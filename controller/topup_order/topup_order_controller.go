package topup_order

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type TopupOrderController interface {
	CreateTopUpOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindUserIdByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

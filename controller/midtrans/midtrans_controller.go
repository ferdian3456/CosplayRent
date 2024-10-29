package midtrans

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MidtransController interface {
	CreateTransaction(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	MidtransCallBack(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

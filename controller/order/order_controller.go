package order

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type OrderController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
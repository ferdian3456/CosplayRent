package rajaongkir

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type RajaOngkirController interface {
	FindProvince(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindCity(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckShippment(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

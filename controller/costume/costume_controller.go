package costume

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CostumeController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByName(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByUserUUID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindSellerCostumeByCostumeID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

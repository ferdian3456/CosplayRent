package rajaongkir

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/rajaongkir"
	rajaongkirs "cosplayrent/service/rajaongkir"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type RajaOngkirControllerImpl struct {
	RajaOngkirService rajaongkirs.RajaOngkirService
}

func NewRajaOngkirController(rajaongkirService rajaongkirs.RajaOngkirService) RajaOngkirController {
	return &RajaOngkirControllerImpl{
		RajaOngkirService: rajaongkirService,
	}
}

func (controller RajaOngkirControllerImpl) FindProvince(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rajaOngkirProvinceResponse, err := controller.RajaOngkirService.FindProvince(request.Context())
	var data interface{}
	if err != nil {
		data = err
	} else {

		data = rajaOngkirProvinceResponse.Rajaongkir.Results
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller RajaOngkirControllerImpl) FindCity(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	provinceID := params.ByName("provinceID")

	rajaOngkirCityResponse, err := controller.RajaOngkirService.FindCity(request.Context(), provinceID)
	var data interface{}
	if err != nil {
		data = err
	} else {
		data = rajaOngkirCityResponse.Rajaongkir.Results
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller RajaOngkirControllerImpl) CheckShippment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	shipmentRequest := rajaongkir.RajaOngkirSendShipmentRequest{}
	helper.ReadFromRequestBody(request, &shipmentRequest)
	rajaOngkirCostResponse, err := controller.RajaOngkirService.CheckShippment(request.Context(), shipmentRequest)

	var data interface{}
	if err != nil {
		data = err
	} else {
		data = rajaOngkirCostResponse.Rajaongkir.Results
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

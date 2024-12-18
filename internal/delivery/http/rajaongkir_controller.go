package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/rajaongkir"
	"cosplayrent/internal/usecase"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type RajaOngkirController struct {
	RajaOngkirUsecase *usecase.RajaOngkirUsecase
	Log               *zerolog.Logger
}

func NewRajaOngkirController(RajaOngkirUsecase *usecase.RajaOngkirUsecase, zerolog *zerolog.Logger) *RajaOngkirController {
	return &RajaOngkirController{
		RajaOngkirUsecase: RajaOngkirUsecase,
		Log:               zerolog,
	}
}

func (controller RajaOngkirController) FindProvince(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rajaOngkirProvinceResponse, err := controller.RajaOngkirUsecase.FindProvince(request.Context())
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

func (controller RajaOngkirController) FindCity(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	provinceID := params.ByName("provinceID")

	rajaOngkirCityResponse, err := controller.RajaOngkirUsecase.FindCity(request.Context(), provinceID)
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

func (controller RajaOngkirController) CheckShippment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	shipmentRequest := rajaongkir.RajaOngkirSendShipmentRequest{}
	helper.ReadFromRequestBody(request, &shipmentRequest)

	rajaOngkirCostResponse, err := controller.RajaOngkirUsecase.CheckShippment(request.Context(), shipmentRequest)
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

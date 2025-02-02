package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/rajaongkir"
	"cosplayrent/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type RajaOngkirController struct {
	RajaOngkirUsecase *usecase.RajaOngkirUsecase
	Log               *zerolog.Logger
}

func NewRajaOngkirController(rajaongkirUsecase *usecase.RajaOngkirUsecase, zerolog *zerolog.Logger) *RajaOngkirController {
	return &RajaOngkirController{
		RajaOngkirUsecase: rajaongkirUsecase,
		Log:               zerolog,
	}
}

func (controller RajaOngkirController) FindProvince(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	rajaOngkirProvinceResponse := controller.RajaOngkirUsecase.FindProvince(request.Context())

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   rajaOngkirProvinceResponse.Rajaongkir.Results,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller RajaOngkirController) FindCity(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	provinceID := params.ByName("provinceID")

	rajaOngkirCityResponse := controller.RajaOngkirUsecase.FindCity(request.Context(), provinceID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   rajaOngkirCityResponse.Rajaongkir.Results,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller RajaOngkirController) CheckShippment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	shipmentRequest := rajaongkir.RajaOngkirSendShipmentRequest{}
	helper.ReadFromRequestBody(request, &shipmentRequest)

	rajaOngkirCostResponse, err := controller.RajaOngkirUsecase.CheckShippment(request.Context(), shipmentRequest)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   rajaOngkirCostResponse.Rajaongkir.Results,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

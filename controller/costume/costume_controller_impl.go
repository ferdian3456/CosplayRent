package costume

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/costume"
	costumes "cosplayrent/service/costume"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type CostumeControllerImpl struct {
	CostumeService costumes.CostumeService
}

func NewCostumeController(costumeService costumes.CostumeService) CostumeController {
	return &CostumeControllerImpl{
		CostumeService: costumeService,
	}
}

func (controller CostumeControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeCreateRequest := costume.CostumeCreateRequest{}
	helper.ReadFromRequestBody(request, &costumeCreateRequest)

	controller.CostumeService.Create(request.Context(), costumeCreateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	id, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	costumeResponse := controller.CostumeService.FindById(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeResponse := controller.CostumeService.FindAll(request.Context())

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	id, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	costumeUpdateRequest := costume.CostumeUpdateRequest{}
	helper.ReadFromRequestBody(request, &costumeUpdateRequest)

	costumeUpdateRequest.Id = id
	controller.CostumeService.Update(request.Context(), costumeUpdateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	id, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	controller.CostumeService.Delete(request.Context(), id)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

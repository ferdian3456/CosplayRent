package costume

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/costume"
	costumes "cosplayrent/service/costume"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
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
	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	costume_userId := request.FormValue("user_id")
	costumeName := request.FormValue("name")
	costumeDescription := request.FormValue("description")
	costumeBahan := request.FormValue("bahan")
	costumeUkuran := request.FormValue("ukuran")
	costumeBerat := request.FormValue("berat")
	costumeKategori := request.FormValue("kategori")
	costumePrice := request.FormValue("price")

	file, handler, err := request.FormFile("costume_picture")
	helper.PanicIfError(err)
	defer file.Close()

	if _, err := os.Stat("../static/costume/"); os.IsNotExist(err) {
		err = os.MkdirAll("../static/costume/", os.ModePerm)
		helper.PanicIfError(err)
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
	costumeImagePath := filepath.Join("../static/costume/", fileName)

	destFile, err := os.Create(costumeImagePath)

	helper.PanicIfError(err)

	_, err = io.Copy(destFile, file)
	helper.PanicIfError(err)

	defer destFile.Close()

	var costumeFixPrice float64
	costumeFixPrice, err = strconv.ParseFloat(costumePrice, 64)
	costumeImageTrimPath := strings.TrimPrefix(costumeImagePath, "..")

	err = godotenv.Load("../.env")
	helper.PanicIfError(err)

	imageEnv := os.Getenv("IMAGE_ENV")

	costumeFinalPath := fmt.Sprintf(imageEnv + costumeImageTrimPath)
	costumeRequest := costume.CostumeCreateRequest{
		User_id:     costume_userId,
		Name:        costumeName,
		Description: costumeDescription,
		Bahan:       costumeBahan,
		Ukuran:      costumeUkuran,
		Berat:       costumeBerat,
		Kategori:    costumeKategori,
		Price:       costumeFixPrice,
		Picture:     costumeFinalPath,
	}

	controller.CostumeService.Create(request.Context(), costumeRequest)

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

func (controller CostumeControllerImpl) FindByName(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeName := params.ByName("costumeName")

	costumeResponse := controller.CostumeService.FindByName(request.Context(), costumeName)

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

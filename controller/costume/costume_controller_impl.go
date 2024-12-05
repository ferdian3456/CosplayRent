package costume

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/costume"
	costumes "cosplayrent/service/costume"
	"fmt"
	"io"
	"log"
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

	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Costume Controller: Create", userUUID)

	costumeName := request.FormValue("name")
	costumeDescription := request.FormValue("description")
	costumeBahan := request.FormValue("bahan")
	costumeUkuran := request.FormValue("ukuran")
	costumeBerat := request.FormValue("berat")
	costumeKategori := request.FormValue("kategori")
	costumePrice := request.FormValue("price")

	var costumePicturePath *string

	if file, handler, err := request.FormFile("costume_picture"); err == nil {
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

		costumeImageTrimPath := strings.TrimPrefix(costumeImagePath, "..")

		costumePicturePath = &costumeImageTrimPath
	}
	fixPrice, err := strconv.ParseFloat(costumePrice, 64)
	helper.PanicIfError(err)

	fixBerat, err := strconv.Atoi(costumeBerat)
	helper.PanicIfError(err)

	costumeRequest := costume.CostumeCreateRequest{
		User_id:     userUUID,
		Name:        costumeName,
		Description: costumeDescription,
		Bahan:       costumeBahan,
		Ukuran:      costumeUkuran,
		Berat:       fixBerat,
		Kategori:    costumeKategori,
		Price:       fixPrice,
		Picture:     costumePicturePath,
	}

	controller.CostumeService.Create(request.Context(), costumeRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Costume Controller: Update", userUUID)

	costumeID := params.ByName("costumeID")
	costumeId, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)
	costumeName := request.FormValue("name")
	costumeDescription := request.FormValue("description")
	costumeBahan := request.FormValue("bahan")
	costumeUkuran := request.FormValue("ukuran")
	costumeBerat := request.FormValue("berat")
	costumeKategori := request.FormValue("kategori")
	costumePrice := request.FormValue("price")

	var costumePicturePath *string

	if file, handler, err := request.FormFile("costume_picture"); err == nil {
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

		costumeImageTrimPath := strings.TrimPrefix(costumeImagePath, "..")

		costumePicturePath = &costumeImageTrimPath
	}

	fixPrice, err := strconv.ParseFloat(costumePrice, 64)
	helper.PanicIfError(err)

	fixBerat, err := strconv.Atoi(costumeBerat)
	helper.PanicIfError(err)

	costumeRequest := costume.CostumeUpdateRequest{
		Id:          costumeId,
		Name:        &costumeName,
		Description: &costumeDescription,
		Bahan:       &costumeBahan,
		Ukuran:      &costumeUkuran,
		Berat:       &fixBerat,
		Kategori:    &costumeKategori,
		Price:       &fixPrice,
		Picture:     costumePicturePath,
	}

	controller.CostumeService.Update(request.Context(), costumeRequest, userUUID)

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

func (controller CostumeControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Costume Controller: Delete", userUUID)

	costumeID := params.ByName("costumeID")
	id, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	controller.CostumeService.Delete(request.Context(), id, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) FindByUserUUID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID := params.ByName("userUUID")

	costumeReturn := controller.CostumeService.FindByUserUUID(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeReturn,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) FindSellerCostumeByCostumeID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Costume Controller: FindSellerCostumeByCostumeID", userUUID)

	costumeID := params.ByName("costumeID")
	finalCostumeID, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	costumeReturn := controller.CostumeService.FindSellerCostumeByCostumeID(request.Context(), userUUID, finalCostumeID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeReturn,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeControllerImpl) FindSellerCostume(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	log.Printf("User with uuid: %s enter Costume Controller: FindSellerCostume", userUUID)

	costumeReturn := controller.CostumeService.FindSellerCostume(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeReturn,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

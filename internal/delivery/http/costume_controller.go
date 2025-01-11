package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/costume"
	"cosplayrent/internal/usecase"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CostumeController struct {
	CostumeUsecase *usecase.CostumeUsecase
	Log            *zerolog.Logger
}

func NewCostumeController(costumeUsecase *usecase.CostumeUsecase, zerolog *zerolog.Logger) *CostumeController {
	return &CostumeController{
		CostumeUsecase: costumeUsecase,
		Log:            zerolog,
	}
}

func (controller CostumeController) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	request.Body = http.MaxBytesReader(writer, request.Body, 5*1024*1024) // 5 MB

	file, fileHeader, err := request.FormFile("costume_picture")

	var costumePicturePath *string

	if err != nil {
		if err.Error() == "http: no such file" {

		} else if err.Error() == "http: request body too large" {
			respErr := errors.New("request exceeded 5 mb")
			webResponse := web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   respErr.Error(),
			}

			controller.Log.Warn().Err(err).Msg(respErr.Error())
			helper.WriteToResponseBody(writer, webResponse)
			return
		} else {
			respErr := errors.New("unexpected error handling file upload")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	} else if file != nil {
		defer file.Close()

		fileType := fileHeader.Header.Get("Content-Type")
		if fileType != "image/jpeg" && fileType != "image/png" {
			respErr := errors.New("file is not image")
			webResponse := web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   respErr.Error(),
			}

			controller.Log.Warn().Err(err).Msg(respErr.Error())
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		_, err = os.Stat("../static/costume/")
		if os.IsNotExist(err) {
			err = os.MkdirAll("../static/costume/", os.ModePerm)
			if err != nil {
				respErr := errors.New("failed to create directory")
				controller.Log.Panic().Err(err).Msg(respErr.Error())
			}
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		costumeImagePath := filepath.Join("../static/costume/", fileName)

		destFile, err := os.Create(costumeImagePath)
		if err != nil {
			respErr := errors.New("failed to create file in the directory path")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			respErr := errors.New("failed to copy a created file from request's file")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}

		costumeImageTrimPath := strings.TrimPrefix(costumeImagePath, "..")
		costumePicturePath = &costumeImageTrimPath
	}

	costumeName := request.FormValue("name")
	costumeDescription := request.FormValue("description")
	costumeBahan := request.FormValue("bahan")
	costumeUkuran := request.FormValue("ukuran")
	costumeBerat := request.FormValue("berat")
	costumeKategori := request.FormValue("kategori")
	costumePrice := request.FormValue("price")

	var fixPrice float64
	if costumePrice != "" {
		fixPrice, err = strconv.ParseFloat(costumePrice, 64)
		if err != nil {
			respErr := errors.New("error converting string to float64")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	var fixBerat int
	if costumeBerat != "" {
		fixBerat, err = strconv.Atoi(costumeBerat)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	var fixKategoriId int
	if costumeKategori != "" {
		fixKategoriId, err = strconv.Atoi(costumeKategori)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	userRequest := costume.CostumeCreateRequest{
		Name:        costumeName,
		Description: costumeDescription,
		Bahan:       costumeBahan,
		Ukuran:      costumeUkuran,
		Berat:       fixBerat,
		Kategori:    fixKategoriId,
		Price:       fixPrice,
		Picture:     costumePicturePath,
	}

	err = controller.CostumeUsecase.Create(request.Context(), userRequest, userUUID)
	if err != nil {
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
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeController) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		if request.MultipartForm == nil || len(request.MultipartForm.Value) == 0 {
			respErr := errors.New("request contains no data or only empty fields")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		} else {
			respErr := errors.New("request exceeded 10 MB")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	costumeID := params.ByName("costumeID")
	costumeName := request.FormValue("name")
	costumeDescription := request.FormValue("description")
	costumeBahan := request.FormValue("bahan")
	costumeUkuran := request.FormValue("ukuran")
	costumeBerat := request.FormValue("berat")
	costumeKategoriId := request.FormValue("kategori")
	costumeAvailable := request.FormValue("available")
	costumePrice := request.FormValue("price")

	var costumePicturePath *string

	file, handler, err := request.FormFile("costume_picture")
	if err == nil {
		defer file.Close()

		_, err := os.Stat("../static/costume/")
		if os.IsNotExist(err) {
			err = os.MkdirAll("../static/costume/", os.ModePerm)
			if err != nil {
				respErr := errors.New("failed to create directory")
				controller.Log.Panic().Err(err).Msg(respErr.Error())
			}
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		costumeImagePath := filepath.Join("../static/costume/", fileName)

		destFile, err := os.Create(costumeImagePath)
		if err != nil {
			respErr := errors.New("failed to create file in the directory path")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}

		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			respErr := errors.New("failed to copy a created file from request's file")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}

		costumeImageTrimPath := strings.TrimPrefix(costumeImagePath, "..")

		costumePicturePath = &costumeImageTrimPath
	} else {
		var emptyPicture string = ""
		costumePicturePath = &emptyPicture
	}

	fixId, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	var fixPrice float64
	if costumePrice != "" {
		fixPrice, err = strconv.ParseFloat(costumePrice, 64)
		if err != nil {
			respErr := errors.New("error converting string to float64")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	var fixBerat int
	if costumeBerat != "" {
		fixBerat, err = strconv.Atoi(costumeBerat)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	var fixKategoriId int
	if costumeKategoriId != "" {
		fixKategoriId, err = strconv.Atoi(costumeKategoriId)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	costumeRequest := costume.CostumeUpdateRequest{
		Id:          fixId,
		Name:        costumeName,
		Description: costumeDescription,
		Bahan:       costumeBahan,
		Ukuran:      costumeUkuran,
		Berat:       fixBerat,
		Kategori:    fixKategoriId,
		Available:   costumeAvailable,
		Price:       fixPrice,
		Picture:     costumePicturePath,
	}

	err = controller.CostumeUsecase.Update(request.Context(), costumeRequest, userUUID)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeController) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeResponse, err := controller.CostumeUsecase.FindAll(request.Context())
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeController) FindSellerCostume(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeReturn, err := controller.CostumeUsecase.FindSellerCostume(request.Context(), userUUID)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeReturn,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeController) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	id, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	costumeResponse, err := controller.CostumeUsecase.FindById(request.Context(), id)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeController) FindSellerCostumeByCostumeID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeID := params.ByName("costumeID")
	finalCostumeID, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	costumeReturn, err := controller.CostumeUsecase.FindSellerCostumeByCostumeID(request.Context(), userUUID, finalCostumeID)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   costumeReturn,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller CostumeController) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeID := params.ByName("costumeID")
	id, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	err = controller.CostumeUsecase.Delete(request.Context(), id, userUUID)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Data:   err.Error(),
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

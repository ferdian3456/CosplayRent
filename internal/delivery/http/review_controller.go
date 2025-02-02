package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/review"
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

type ReviewController struct {
	ReviewUsecase *usecase.ReviewUsecase
	OrderUsecase  *usecase.OrderUsecase
	Log           *zerolog.Logger
}

func NewReviewController(reviewUsecase *usecase.ReviewUsecase, orderUsecase *usecase.OrderUsecase, zerolog *zerolog.Logger) *ReviewController {
	return &ReviewController{
		ReviewUsecase: reviewUsecase,
		OrderUsecase:  orderUsecase,
		Log:           zerolog,
	}
}

func (controller ReviewController) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	request.Body = http.MaxBytesReader(writer, request.Body, 5*1024*1024) // 5 MB

	file, fileHeader, err := request.FormFile("review_picture")

	var reviewPicturePath *string

	if err != nil {
		if err.Error() == "http: no such file" {

		} else if err.Error() == "http: request body too large" {
			respErr := errors.New("request exceeded 5 mb")

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)

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

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)

			webResponse := web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   respErr.Error(),
			}

			controller.Log.Warn().Err(err).Msg(respErr.Error())
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		_, err = os.Stat("../static/review/")
		if os.IsNotExist(err) {
			err = os.MkdirAll("../static/review/", os.ModePerm)
			if err != nil {
				respErr := errors.New("failed to create directory")
				controller.Log.Panic().Err(err).Msg(respErr.Error())
			}
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		reviewImagePath := filepath.Join("../static/review/", fileName)

		destFile, err := os.Create(reviewImagePath)
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

		reviewImageTrimPath := strings.TrimPrefix(reviewImagePath, "..")
		reviewPicturePath = &reviewImageTrimPath
	}

	costumeID := request.FormValue("costume_id")
	costumeDescription := request.FormValue("description")
	costumeRating := request.FormValue("rating")
	orderID := request.FormValue("order_id")

	finalCostumeID, err := strconv.Atoi(costumeID)
	fmt.Println(costumeID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	finalRating, err := strconv.Atoi(costumeRating)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	reviewRequest := review.ReviewCreateRequest{
		Order_id:       orderID,
		Customer_id:    userUUID,
		Costume_id:     finalCostumeID,
		Description:    costumeDescription,
		Review_picture: reviewPicturePath,
		Rating:         finalRating,
	}

	err = controller.ReviewUsecase.Create(request.Context(), reviewRequest, userUUID)
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
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	request.Body = http.MaxBytesReader(writer, request.Body, 5*1024*1024) // 5 MB

	file, fileHeader, err := request.FormFile("review_picture")

	var reviewPicturePath *string

	if err != nil {
		if err.Error() == "http: no such file" {

		} else if err.Error() == "http: request body too large" {
			respErr := errors.New("request exceeded 5 mb")

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)

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

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)

			webResponse := web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   respErr.Error(),
			}

			controller.Log.Warn().Err(err).Msg(respErr.Error())
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		_, err = os.Stat("../static/review/")
		if os.IsNotExist(err) {
			err = os.MkdirAll("../static/review/", os.ModePerm)
			if err != nil {
				respErr := errors.New("failed to create directory")
				controller.Log.Panic().Err(err).Msg(respErr.Error())
			}
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		reviewImagePath := filepath.Join("../static/review/", fileName)

		destFile, err := os.Create(reviewImagePath)
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

		reviewImageTrimPath := strings.TrimPrefix(reviewImagePath, "..")
		reviewPicturePath = &reviewImageTrimPath
	}

	reviewID := params.ByName("reviewID")
	reviewDescription := request.FormValue("description")
	reviewRating := request.FormValue("rating")

	finalReviewId, err := strconv.Atoi(reviewID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	finalRating, err := strconv.Atoi(reviewRating)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	reviewRequest := review.ReviewUpdateRequest{
		Review_picture: reviewPicturePath,
		Description:    reviewDescription,
		Rating:         finalRating,
	}

	err = controller.ReviewUsecase.Update(request.Context(), reviewRequest, userUUID, finalReviewId)
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
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindUserReview(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	reviewResponse, err := controller.ReviewUsecase.FindUserReview(request.Context(), userUUID)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

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
		Data:   reviewResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindByCostumeId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	costumeId, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	reviewDomain, err := controller.ReviewUsecase.FindByCostumeId(request.Context(), costumeId)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

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
		Data:   reviewDomain,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindAllUserReview(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	detailOrderResult, err := controller.ReviewUsecase.FindAllUserReview(request.Context(), userUUID)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

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
		Data:   detailOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindReviewInfoByOrderId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	orderID := params.ByName("orderID")

	detailOrderResult, err := controller.ReviewUsecase.FindReviewInfoByOrderId(request.Context(), userUUID, orderID)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

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
		Data:   detailOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindAllReviewedOrder(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	detailOrderResult, err := controller.ReviewUsecase.FindAllReviewedOrder(request.Context(), userUUID)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

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
		Data:   detailOrderResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) DeleteUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	reviewID := params.ByName("reviewID")
	finalReviewID, err := strconv.Atoi(reviewID)
	if err != nil {
		respErr := errors.New("error converting string to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	err = controller.ReviewUsecase.DeleteUserReviewByReviewID(request.Context(), userUUID, finalReviewID)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

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

package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/review"
	"cosplayrent/internal/usecase"
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
	"github.com/rs/zerolog"
)

type ReviewController struct {
	ReviewUsecase *usecase.ReviewUsecase
	Log           *zerolog.Logger
}

func NewReviewController(ReviewUsecase *usecase.ReviewUsecase, zerolog *zerolog.Logger) *ReviewController {
	return &ReviewController{
		ReviewUsecase: ReviewUsecase,
		Log:           zerolog,
	}
}

func (controller ReviewController) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Review Controller: Create", userUUID)

	costumeID := request.FormValue("costume_id")
	costumeDescription := request.FormValue("description")
	costumeRating := request.FormValue("rating")

	var reviewPicturePath *string

	if file, handler, err := request.FormFile("review_picture"); err == nil {
		defer file.Close()

		if _, err := os.Stat("../static/review/"); os.IsNotExist(err) {
			err = os.MkdirAll("../static/review/", os.ModePerm)
			helper.PanicIfError(err)
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		reviewImagePath := filepath.Join("../static/review/", fileName)

		destFile, err := os.Create(reviewImagePath)
		helper.PanicIfError(err)

		_, err = io.Copy(destFile, file)
		helper.PanicIfError(err)

		defer destFile.Close()

		reviewImageTrimPath := strings.TrimPrefix(reviewImagePath, "..")

		reviewPicturePath = &reviewImageTrimPath
	}

	finalCostumeID, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)
	finalRating, err := strconv.Atoi(costumeRating)
	helper.PanicIfError(err)

	reviewRequest := review.ReviewCreateRequest{
		User_id:        userUUID,
		Costume_id:     finalCostumeID,
		Description:    costumeDescription,
		Review_picture: reviewPicturePath,
		Rating:         finalRating,
	}

	controller.ReviewUsecase.Create(request.Context(), reviewRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Review Controller: Update", userUUID)

	reviewID := params.ByName("reviewID")
	reviewDescription := request.FormValue("description")
	reviewRating := request.FormValue("rating")

	var reviewPicturePath *string

	if file, handler, err := request.FormFile("review_picture"); err == nil {
		defer file.Close()

		if _, err := os.Stat("../static/review/"); os.IsNotExist(err) {
			err = os.MkdirAll("../static/review/", os.ModePerm)
			helper.PanicIfError(err)
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		reviewImagePath := filepath.Join("../static/review/", fileName)

		destFile, err := os.Create(reviewImagePath)
		helper.PanicIfError(err)

		_, err = io.Copy(destFile, file)
		helper.PanicIfError(err)

		defer destFile.Close()

		reviewImageTrimPath := strings.TrimPrefix(reviewImagePath, "..")

		reviewPicturePath = &reviewImageTrimPath
	}

	finalReviewID, err := strconv.Atoi(reviewID)
	helper.PanicIfError(err)

	reviewRequest := review.ReviewUpdateRequest{
		ReviewId:       finalReviewID,
		Review_picture: reviewPicturePath,
		Description:    &reviewDescription,
		Rating:         &reviewRating,
	}

	controller.ReviewUsecase.Update(request.Context(), reviewRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindByCostumeId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	costumeId, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	reviewDomain := controller.ReviewUsecase.FindByCostumeId(request.Context(), costumeId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   reviewDomain,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindUserReview(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Review Controller: FindUserReview", userUUID)

	reviewResponse := controller.ReviewUsecase.FindUserReview(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   reviewResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) FindUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Review Controller: FindUserReview", userUUID)

	reviewID := params.ByName("reviewID")
	finalReviewID, err := strconv.Atoi(reviewID)
	helper.PanicIfError(err)

	reviewResults := controller.ReviewUsecase.FindUserReviewByReviewID(request.Context(), userUUID, finalReviewID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   reviewResults,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewController) DeleteUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter Review Controller: FindUserReview", userUUID)

	reviewID := params.ByName("reviewID")
	finalReviewID, err := strconv.Atoi(reviewID)
	helper.PanicIfError(err)

	controller.ReviewUsecase.DeleteUserReviewByReviewID(request.Context(), userUUID, finalReviewID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

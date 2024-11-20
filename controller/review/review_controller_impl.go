package review

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/review"
	reviews "cosplayrent/service/review"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ReviewControllerImpl struct {
	ReviewService reviews.ReviewService
}

func NewReviewController(reviewService reviews.ReviewService) ReviewController {
	return &ReviewControllerImpl{
		ReviewService: reviewService,
	}
}

func (controller ReviewControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

		err = godotenv.Load("../.env")
		helper.PanicIfError(err)

		imageEnv := os.Getenv("IMAGE_ENV")

		reviewFinalPath := fmt.Sprintf(imageEnv + reviewImageTrimPath)
		reviewPicturePath = &reviewFinalPath
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

	controller.ReviewService.Create(request.Context(), reviewRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

		err = godotenv.Load("../.env")
		helper.PanicIfError(err)

		imageEnv := os.Getenv("IMAGE_ENV")

		reviewFinalPath := fmt.Sprintf(imageEnv + reviewImageTrimPath)
		reviewPicturePath = &reviewFinalPath
	}

	finalReviewID, err := strconv.Atoi(reviewID)
	helper.PanicIfError(err)

	reviewRequest := review.ReviewUpdateRequest{
		ReviewId:       finalReviewID,
		Review_picture: reviewPicturePath,
		Description:    &reviewDescription,
		Rating:         &reviewRating,
	}

	controller.ReviewService.Update(request.Context(), reviewRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewControllerImpl) FindByCostumeId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	costumeID := params.ByName("costumeID")
	costumeId, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	reviewDomain := controller.ReviewService.FindByCostumeId(request.Context(), costumeId)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   reviewDomain,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewControllerImpl) FindUserReview(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	reviewResponse := controller.ReviewService.FindUserReview(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   reviewResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewControllerImpl) FindUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	reviewResults := controller.ReviewService.FindUserReviewByReviewID(request.Context(), userUUID, finalReviewID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   reviewResults,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller ReviewControllerImpl) DeleteUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	controller.ReviewService.DeleteUserReviewByReviewID(request.Context(), userUUID, finalReviewID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

package review

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/review"
	reviews "cosplayrent/service/review"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
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
	costumeID := params.ByName("costumeID")
	costumeId, err := strconv.Atoi(costumeID)
	helper.PanicIfError(err)

	reviewCreateRequest := review.ReviewCreateRequest{}
	helper.ReadFromRequestBody(request, &reviewCreateRequest)
	reviewCreateRequest.Costume_id = costumeId
	controller.ReviewService.Create(request.Context(), reviewCreateRequest)

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

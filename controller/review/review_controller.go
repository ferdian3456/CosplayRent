package review

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ReviewController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByCostumeId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindUserReview(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteUserReviewByReviewID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
)

type CategoryController struct {
	CategoryUsecase *usecase.CategoryUsecase
	Log             *zerolog.Logger
}

func NewCategoryController(categoryUsecase *usecase.CategoryUsecase, zerolog *zerolog.Logger) *CategoryController {
	return &CategoryController{
		CategoryUsecase: categoryUsecase,
		Log:             zerolog,
	}
}

func (controller CategoryController) FindAllCategory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryResponse, err := controller.CategoryUsecase.FindAllCategory(request.Context())
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
		Data:   categoryResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

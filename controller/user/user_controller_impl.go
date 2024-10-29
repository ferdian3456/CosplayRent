package user

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/user"
	users "cosplayrent/service/user"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

type UserControllerImpl struct {
	UserService users.UserService
}

func NewUserController(userService users.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller UserControllerImpl) Register(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userCreateRequest := user.UserCreateRequest{}
	helper.ReadFromRequestBody(request, &userCreateRequest)

	token := controller.UserService.Create(request.Context(), userCreateRequest)
	tokenResponse := web.TokenResponse{
		Token: token,
	}
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   tokenResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userLoginRequest := user.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &userLoginRequest)

	token := controller.UserService.Login(request.Context(), userLoginRequest)
	log.Println(token)
	tokenResponse := web.TokenResponse{
		Token: token,
	}
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   tokenResponse,
	}
	log.Println(token)
	log.Println(tokenResponse.Token)

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) FindByUUID(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID := params.ByName("userUUID")

	userResponse := controller.UserService.FindByUUID(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userResponse := controller.UserService.FindAll(request.Context())

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID := params.ByName("userUUID")

	userUpdateRequest := user.UserUpdateRequest{}
	helper.ReadFromRequestBody(request, &userUpdateRequest)

	userUpdateRequest.Id = userUUID
	controller.UserService.Update(request.Context(), userUpdateRequest)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID := params.ByName("userUUID")

	controller.UserService.Delete(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) VerifyAndRetrieve(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	tokenHeader := request.Header.Get("Authorization")
	if tokenHeader == "" {
		webResponsel := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "No authorization in header ",
		}
		helper.WriteToResponseBody(writer, webResponsel)
		return
	}

	tokenAfter := strings.TrimPrefix(tokenHeader, "Bearer ")
	if tokenAfter == "" {
		webResponsel := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "No token after beerer ",
		}
		helper.WriteToResponseBody(writer, webResponsel)
	}
	userDomain, _ := controller.UserService.VerifyAndRetrieve(request.Context(), tokenAfter)

	webResponsel := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userDomain,
	}

	helper.WriteToResponseBody(writer, webResponsel)
}

package user

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/user"
	users "cosplayrent/service/user"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	userName := request.FormValue("name")
	userEmail := request.FormValue("email")
	userAddress := request.FormValue("address")

	var profilePicturePath *string

	if file, handler, err := request.FormFile("profile_picture"); err == nil {
		defer file.Close()

		if _, err := os.Stat("../static/profile/"); os.IsNotExist(err) {
			err = os.MkdirAll("../static/profile/", os.ModePerm)
			helper.PanicIfError(err)
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		profileImagePath := filepath.Join("../static/profile/", fileName)

		destFile, err := os.Create(profileImagePath)
		helper.PanicIfError(err)

		_, err = io.Copy(destFile, file)
		helper.PanicIfError(err)

		defer destFile.Close()

		userImageTrimPath := strings.TrimPrefix(profileImagePath, "..")

		err = godotenv.Load("../.env")
		helper.PanicIfError(err)

		imageEnv := os.Getenv("IMAGE_ENV")

		userFinalPath := fmt.Sprintf(imageEnv + userImageTrimPath)
		profilePicturePath = &userFinalPath
	}

	// Create the user update request
	userRequest := user.UserUpdateRequest{
		Id:              userUUID,
		Name:            &userName,
		Email:           &userEmail,
		Address:         &userAddress,
		Profile_picture: profilePicturePath,
	}

	controller.UserService.Update(request.Context(), userRequest)

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

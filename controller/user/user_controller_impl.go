package user

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"cosplayrent/model/web/user"
	users "cosplayrent/service/user"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
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

	token, err := controller.UserService.Login(request.Context(), userLoginRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Unauthorized",
			Data:   "Wrong email or password",
		}

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

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

	log.Printf("User with uuid: %s enter User Controller: FindByUUID", userUUID)

	userResponse := controller.UserService.FindByUUID(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: FindAll", userUUID)

	userResponse := controller.UserService.FindAll(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, ok := request.Context().Value("user_uuid").(string)
	if !ok {
		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "Unauthorized",
			Data:   "Invalid Token",
		}
		helper.WriteToResponseBody(writer, webResponse)
	}

	log.Printf("User with uuid: %s enter User Controller: Update", userUUID)

	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	userName := request.FormValue("name")
	userEmail := request.FormValue("email")
	userAddress := request.FormValue("address")
	userOriginCityName := request.FormValue("origin_city_name")
	userOriginProvinceName := request.FormValue("origin_province_name")

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

		profilePicturePath = &userImageTrimPath
	}

	userRequest := user.UserUpdateRequest{
		Id:                   userUUID,
		Name:                 &userName,
		Email:                &userEmail,
		Address:              &userAddress,
		Profile_picture:      profilePicturePath,
		Origin_province_name: &userOriginProvinceName,
		Origin_city_name:     &userOriginCityName,
	}

	controller.UserService.Update(request.Context(), userRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: Delete", userUUID)

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
			Status: "No Authorization in header ",
		}
		helper.WriteToResponseBody(writer, webResponsel)
		return
	}

	userDomain, _ := controller.UserService.VerifyAndRetrieve(request.Context(), tokenHeader)

	webResponsel := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   userDomain,
	}

	helper.WriteToResponseBody(writer, webResponsel)
}

func (controller *UserControllerImpl) AddIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: AddIdentityCard", userUUID)

	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	log.Printf("User with uuid: %s enter User Controller: Update", userUUID)

	err = request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	var IdentityCardPicturePath *string

	if file, handler, err := request.FormFile("identity_card"); err == nil {
		defer file.Close()

		if _, err := os.Stat("../static/identity_card/"); os.IsNotExist(err) {
			err = os.MkdirAll("../static/identity_card/", os.ModePerm)
			helper.PanicIfError(err)
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		IdentityCardImagePath := filepath.Join("../static/identity_card/", fileName)

		destFile, err := os.Create(IdentityCardImagePath)
		helper.PanicIfError(err)

		_, err = io.Copy(destFile, file)
		helper.PanicIfError(err)

		defer destFile.Close()

		IdentityCardImageTrimPath := strings.TrimPrefix(IdentityCardImagePath, "..")

		IdentityCardPicturePath = &IdentityCardImageTrimPath
	}

	controller.UserService.AddIdentityCard(request.Context(), userUUID, *IdentityCardPicturePath)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) GetIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: GetIdentityCard", userUUID)

	identityCardResult := controller.UserService.GetIdentityCard(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   identityCardResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) UpdateIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: UpdateIdentityCard", userUUID)

	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	var IdentityCardPicturePath *string

	if file, handler, err := request.FormFile("identity_card"); err == nil {
		defer file.Close()

		if _, err := os.Stat("../static/identity_card/"); os.IsNotExist(err) {
			err = os.MkdirAll("../static/identity_card/", os.ModePerm)
			helper.PanicIfError(err)
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		IdentityCardImagePath := filepath.Join("../static/identity_card/", fileName)

		destFile, err := os.Create(IdentityCardImagePath)
		helper.PanicIfError(err)

		_, err = io.Copy(destFile, file)
		helper.PanicIfError(err)

		defer destFile.Close()

		IdentityCardImageTrimPath := strings.TrimPrefix(IdentityCardImagePath, "..")

		IdentityCardPicturePath = &IdentityCardImageTrimPath
	}

	controller.UserService.UpdateIdentityCard(request.Context(), userUUID, *IdentityCardPicturePath)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) GetEMoneyAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: GetEMoneyAmount", userUUID)

	eMoneyResult := controller.UserService.GetEMoneyAmount(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   eMoneyResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserControllerImpl) TopUp(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

	log.Printf("User with uuid: %s enter User Controller: TopUp", userUUID)

	topUpEMoneyRequest := user.TopUpEmoney{}
	helper.ReadFromRequestBody(request, &topUpEMoneyRequest)

	controller.UserService.TopUp(request.Context(), topUpEMoneyRequest, userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

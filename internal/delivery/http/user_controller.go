package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/user"
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

type UserController struct {
	UserUsecase *usecase.UserUsecase
	Log         *zerolog.Logger
}

func NewUserController(UserUsecase *usecase.UserUsecase, zerolog *zerolog.Logger) *UserController {
	return &UserController{
		UserUsecase: UserUsecase,
		Log:         zerolog,
	}
}

func (controller UserController) Register(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userCreateRequest := user.UserCreateRequest{}
	helper.ReadFromRequestBody(request, &userCreateRequest)

	token, err := controller.UserUsecase.Create(request.Context(), userCreateRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
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

func (controller UserController) Login(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userLoginRequest := user.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &userLoginRequest)

	token, err := controller.UserUsecase.Login(request.Context(), userLoginRequest)
	if err != nil {
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   err.Error(),
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

func (controller UserController) FindByUUID(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	userResponse, err := controller.UserUsecase.FindByUUID(request.Context(), userUUID)

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
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) FindAll(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	userResponse, err := controller.UserUsecase.FindAll(request.Context(), userUUID)
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
		Data:   userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		respErr := errors.New("request exceed 10 mb")
		controller.Log.Panic().Err(respErr).Msg(err.Error())
	}

	userName := request.FormValue("name")
	userEmail := request.FormValue("email")
	userAddress := request.FormValue("address")
	userOriginProvinceName := request.FormValue("origin_province_name")
	userOriginProvinceId := request.FormValue("origin_province_id")
	userOriginCityName := request.FormValue("origin_city_name")
	userOriginCityId := request.FormValue("origin_city_id")

	var profilePicturePath *string

	file, handler, err := request.FormFile("profile_picture")
	if err == nil {
		defer file.Close()

		_, err := os.Stat("../static/profile/")
		if os.IsNotExist(err) {
			err = os.MkdirAll("../static/profile", os.ModePerm)
			if err != nil {
				respErr := errors.New("failed to create directory")
				controller.Log.Panic().Err(respErr).Msg(err.Error())
			}
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		profileImagePath := filepath.Join("../static/profile/", fileName)

		destFile, err := os.Create(profileImagePath)
		if err != nil {
			respErr := errors.New("failed to create file in the directory path")
			controller.Log.Panic().Err(respErr).Msg(err.Error())
		}

		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			respErr := errors.New("failed to copy a created file from request's file")
			controller.Log.Panic().Err(respErr).Msg(err.Error())
		}

		userImageTrimPath := strings.TrimPrefix(profileImagePath, "..")

		profilePicturePath = &userImageTrimPath
	} else {
		var emptyPicture string = ""
		profilePicturePath = &emptyPicture
	}

	var originCityIdFinal int
	if userOriginCityId != "" {
		originCityIdFinal, err = strconv.Atoi(userOriginCityId)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(respErr).Msg(err.Error())
		}
	}

	var originProvinceIdFinal int
	if userOriginProvinceId != "" {
		originProvinceIdFinal, err = strconv.Atoi(userOriginProvinceId)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(respErr).Msg(err.Error())
		}
	}

	userRequest := user.UserPatchRequest{
		Name:                 &userName,
		Email:                &userEmail,
		Address:              &userAddress,
		Profile_picture:      profilePicturePath,
		Origin_province_name: &userOriginProvinceName,
		Origin_province_id:   &originProvinceIdFinal,
		Origin_city_name:     &userOriginCityName,
		Origin_city_id:       &originCityIdFinal,
	}

	err = controller.UserUsecase.Update(request.Context(), userRequest, userUUID)

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

//
//func (controller UserController) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: Delete", userUUID)
//
//	controller.UserUsecase.Delete(request.Context(), userUUID)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) VerifyAndRetrieve(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	tokenHeader := request.Header.Get("Authorization")
//	if tokenHeader == "" {
//		webResponsel := web.WebResponse{
//			Code:   http.StatusBadRequest,
//			Status: "No Authorization in header ",
//		}
//		helper.WriteToResponseBody(writer, webResponsel)
//		return
//	}
//
//	userDomain, _ := controller.UserUsecase.VerifyAndRetrieve(request.Context(), tokenHeader)
//
//	webResponsel := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   userDomain,
//	}
//
//	helper.WriteToResponseBody(writer, webResponsel)
//}
//
//func (controller UserController) AddIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: AddIdentityCard", userUUID)
//
//	err := request.ParseMultipartForm(10 << 20)
//	helper.PanicIfError(err)
//
//	var IdentityCardPicturePath *string
//
//	if file, handler, err := request.FormFile("identity_card"); err == nil {
//		defer file.Close()
//
//		if _, err := os.Stat("../static/identity_card/"); os.IsNotExist(err) {
//			err = os.MkdirAll("../static/identity_card/", os.ModePerm)
//			helper.PanicIfError(err)
//		}
//
//		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
//		IdentityCardImagePath := filepath.Join("../static/identity_card/", fileName)
//
//		destFile, err := os.Create(IdentityCardImagePath)
//		helper.PanicIfError(err)
//
//		_, err = io.Copy(destFile, file)
//		helper.PanicIfError(err)
//
//		defer destFile.Close()
//
//		IdentityCardImageTrimPath := strings.TrimPrefix(IdentityCardImagePath, "..")
//
//		IdentityCardPicturePath = &IdentityCardImageTrimPath
//	}
//
//	controller.UserUsecase.AddIdentityCard(request.Context(), userUUID, *IdentityCardPicturePath)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) GetIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: GetIdentityCard", userUUID)
//
//	identityCardResult := controller.UserUsecase.GetIdentityCard(request.Context(), userUUID)
//
//	identityCardResponse := user.IdentityCardRequest{
//		IdentityCard_picture: identityCardResult,
//	}
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   identityCardResponse,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) UpdateIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: UpdateIdentityCard", userUUID)
//
//	err := request.ParseMultipartForm(10 << 20)
//	helper.PanicIfError(err)
//
//	var IdentityCardPicturePath *string
//
//	if file, handler, err := request.FormFile("identity_card"); err == nil {
//		defer file.Close()
//
//		if _, err := os.Stat("../static/identity_card/"); os.IsNotExist(err) {
//			err = os.MkdirAll("../static/identity_card/", os.ModePerm)
//			helper.PanicIfError(err)
//		}
//
//		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
//		IdentityCardImagePath := filepath.Join("../static/identity_card/", fileName)
//
//		destFile, err := os.Create(IdentityCardImagePath)
//		helper.PanicIfError(err)
//
//		_, err = io.Copy(destFile, file)
//		helper.PanicIfError(err)
//
//		defer destFile.Close()
//
//		IdentityCardImageTrimPath := strings.TrimPrefix(IdentityCardImagePath, "..")
//
//		IdentityCardPicturePath = &IdentityCardImageTrimPath
//	}
//
//	controller.UserUsecase.UpdateIdentityCard(request.Context(), userUUID, *IdentityCardPicturePath)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) GetEMoneyAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: GetEMoneyAmount", userUUID)
//
//	eMoneyResult := controller.UserUsecase.GetEMoneyAmount(request.Context(), userUUID)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   eMoneyResult,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) GetEMoneyTransactionHistory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: GetEMoneyChangeHistory", userUUID)
//
//	eMoneyChangeResult := controller.UserUsecase.GetEMoneyTransactionHistory(request.Context(), userUUID)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   eMoneyChangeResult,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) CheckUserStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: CheckUserStatus", userUUID)
//
//	costumeID := params.ByName("costumeID")
//	finalCostumeID, err := strconv.Atoi(costumeID)
//	helper.PanicIfError(err)
//
//	statusResult := controller.UserUsecase.CheckUserStatus(request.Context(), userUUID, finalCostumeID)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   statusResult,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) GetSellerAddressDetailByCostumeId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: GetSellerAddressDetailByCostumeId", userUUID)
//
//	costumeID := params.ByName("costumeID")
//	finalCostumeId, err := strconv.Atoi(costumeID)
//	helper.PanicIfError(err)
//
//	sellerAddressResult := controller.UserUsecase.GetSellerAddressDetailByCostumeId(request.Context(), userUUID, finalCostumeId)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   sellerAddressResult,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//
//}
//
//func (controller UserController) CheckSellerStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	userUUID, ok := request.Context().Value("user_uuid").(string)
//	if !ok {
//		webResponse := web.WebResponse{
//			Code:   http.StatusInternalServerError,
//			Status: "Unauthorized",
//			Data:   "Invalid Token",
//		}
//		helper.WriteToResponseBody(writer, webResponse)
//		return
//	}
//
//	log.Printf("User with uuid: %s enter User Controller: CheckSellerStatus", userUUID)
//
//	sellerStatusResult := controller.UserUsecase.CheckSellerStatus(request.Context(), userUUID)
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   sellerStatusResult,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}
//
//func (controller UserController) CheckAppVersion(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
//	var err error = godotenv.Load("../.env")
//	helper.PanicIfError(err)
//
//	APP_VERSION := os.Getenv("APP_VERSION")
//
//	AppVersion := web.AppResponse{
//		AppVersion: APP_VERSION,
//	}
//
//	webResponse := web.WebResponse{
//		Code:   200,
//		Status: "OK",
//		Data:   AppVersion,
//	}
//
//	helper.WriteToResponseBody(writer, webResponse)
//}

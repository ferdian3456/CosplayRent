package controller

import (
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/model/web/user"
	"cosplayrent/internal/usecase"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type UserController struct {
	UserUsecase *usecase.UserUsecase
	Log         *zerolog.Logger
}

func NewUserController(userUsecase *usecase.UserUsecase, zerolog *zerolog.Logger) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
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

func (controller UserController) VerifyCode(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userVerificationCodeRequest := user.UserVerificationCode{}
	helper.ReadFromRequestBody(request, &userVerificationCodeRequest)

	userUUID, _ := request.Context().Value("user_uuid").(string)

	err := controller.UserUsecase.VerifyCode(request.Context(), userVerificationCodeRequest, userUUID)
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

	request.Body = http.MaxBytesReader(writer, request.Body, 5*1024*1024) // 5 MB

	file, fileHeader, err := request.FormFile("profile_picture")

	var profilePicturePath *string

	if err != nil {
		if err.Error() == "http: no such file" {

		} else if err.Error() == "http: request body too large" {
			respErr := errors.New("request exceeded 5 mb")
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
			webResponse := web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   respErr.Error(),
			}

			controller.Log.Warn().Err(err).Msg(respErr.Error())
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		_, err = os.Stat("../static/profile/")
		if os.IsNotExist(err) {
			err = os.MkdirAll("../static/profile/", os.ModePerm)
			if err != nil {
				respErr := errors.New("failed to create directory")
				controller.Log.Panic().Err(err).Msg(respErr.Error())
			}
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		profileImagePath := filepath.Join("../static/profile/", fileName)

		destFile, err := os.Create(profileImagePath)
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

		userImageTrimPath := strings.TrimPrefix(profileImagePath, "..")
		profilePicturePath = &userImageTrimPath
	}

	userName := request.FormValue("name")
	userEmail := request.FormValue("email")
	userAddress := request.FormValue("address")
	userOriginProvinceName := request.FormValue("origin_province_name")
	userOriginProvinceId := request.FormValue("origin_province_id")
	userOriginCityName := request.FormValue("origin_city_name")
	userOriginCityId := request.FormValue("origin_city_id")

	var originCityIdFinal int
	if userOriginCityId != "" {
		fmt.Println("id :", userOriginCityId)
		originCityIdFinal, err = strconv.Atoi(userOriginCityId)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	var originProvinceIdFinal int
	if userOriginProvinceId != "" {
		originProvinceIdFinal, err = strconv.Atoi(userOriginProvinceId)
		if err != nil {
			respErr := errors.New("error converting string to int")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
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

func (controller UserController) AddOrUpdateIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	request.Body = http.MaxBytesReader(writer, request.Body, 5*1024*1024) // 5 Mb

	file, fileHeader, err := request.FormFile("identity_card")

	if err != nil {
		if err.Error() == "http: no such file" {
			respErr := errors.New("identity_card is empty")
			webResponse := web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   respErr.Error(),
			}

			controller.Log.Warn().Err(err).Msg(respErr.Error())

			helper.WriteToResponseBody(writer, webResponse)
			return
		}
		respErr := errors.New("request exceeded 5 mb")
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   respErr.Error(),
		}

		controller.Log.Warn().Err(err).Msg(respErr.Error())

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")
	if fileType != "image/jpeg" && fileType != "image/png" {
		respErr := errors.New("file is not image")
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   respErr.Error(),
		}

		controller.Log.Warn().Err(err).Msg(respErr.Error())

		helper.WriteToResponseBody(writer, webResponse)
		return
	}

	defer file.Close()

	var IdentityCardPicturePath *string

	_, err = os.Stat("../static/identity_card/")
	if os.IsNotExist(err) {
		err = os.MkdirAll("../static/identity_card/", os.ModePerm)
		if err != nil {
			respErr := errors.New("failed to create directory")
			controller.Log.Panic().Err(err).Msg(respErr.Error())
		}
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	IdentityCardImagePath := filepath.Join("../static/identity_card/", fileName)

	destFile, err := os.Create(IdentityCardImagePath)
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

	IdentityCardImageTrimPath := strings.TrimPrefix(IdentityCardImagePath, "..")

	IdentityCardPicturePath = &IdentityCardImageTrimPath

	userRequest := user.IdentityCardRequest{
		IdentityCard_picture: IdentityCardPicturePath,
	}

	err = controller.UserUsecase.AddOrUpdateIdentityCard(request.Context(), userUUID, userRequest)
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

func (controller UserController) GetIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	userResponse, err := controller.UserUsecase.GetIdentityCard(request.Context(), userUUID)
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

func (controller UserController) GetEMoneyAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	eMoneyRespone := controller.UserUsecase.GetEMoneyAmount(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   eMoneyRespone,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) GetEMoneyTransactionHistory(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	eMoneyChangeResponse, err := controller.UserUsecase.GetEMoneyTransactionHistory(request.Context(), userUUID)
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
		Data:   eMoneyChangeResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) CheckUserStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeID := params.ByName("costumeID")
	finalCostumeID, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("failed to convert costume id to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	statusResult, err := controller.UserUsecase.CheckUserStatus(request.Context(), userUUID, finalCostumeID)
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
		Data:   statusResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) CheckSellerStatus(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	statusResult, err := controller.UserUsecase.CheckSellerStatus(request.Context(), userUUID)
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
		Data:   statusResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) FindSellerAddressDetailByCostumeId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	costumeID := params.ByName("costumeID")
	finalCostumeId, err := strconv.Atoi(costumeID)
	if err != nil {
		respErr := errors.New("failed to convert costume id to int")
		controller.Log.Panic().Err(err).Msg(respErr.Error())
	}

	sellerAddressResult, err := controller.UserUsecase.FindSellerAddressDetailByCostumeId(request.Context(), userUUID, finalCostumeId)
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
		Data:   sellerAddressResult,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller UserController) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUUID, _ := request.Context().Value("user_uuid").(string)

	controller.UserUsecase.Delete(request.Context(), userUUID)

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

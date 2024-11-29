package user

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UserController interface {
	Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByUUID(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	VerifyAndRetrieve(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	AddIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateIdentityCard(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetEMoneyAmount(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	//TopUp(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

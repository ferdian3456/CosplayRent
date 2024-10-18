package middleware

import (
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strings"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

//func (middleware *AuthMiddleware) ServeHTTP(next httprouter.Handle) httprouter.Handle {
//	return func(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
//		if request.Header.Get("X-API-Key") == "RAHASIA" {
//			log.Print("ini middleware")
//			next(writer, request, p)
//		} else {
//			writer.Header().Set("Content-Type", "application/json")
//			writer.WriteHeader(http.StatusUnauthorized)
//
//			webResponse := web.WebResponse{
//				Code:   http.StatusBadRequest,
//				Status: "Unauthorized",
//			}
//
//			helper.WriteToResponseBody(writer, webResponse)
//		}
//	}
//}

func (middleware *AuthMiddleware) ServeHTTP(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
		headerToken := request.Header.Get("Authorization")
		if headerToken == "" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)

			webResponse := web.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Data:   "No token provided",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		tokenString := strings.TrimPrefix(headerToken, "Bearer ")
		if tokenString == headerToken {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("Unauthorized: Invalid Authorization header format"))
			return
		}

		secretKey := os.Getenv("SECRET_KEY")
		secretKeyByte := []byte(secretKey)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return secretKeyByte, nil
		})

		if err != nil || !token.Valid {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)

			webResponse := web.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Data:   "Invalid token",
			}
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		next(writer, request, p)
	}
}

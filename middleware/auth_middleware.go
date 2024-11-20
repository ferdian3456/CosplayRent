package middleware

import (
	"context"
	"cosplayrent/helper"
	"cosplayrent/model/web"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

const (
	userUUIDkey = "user_uuid"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

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

		var err error = godotenv.Load("../.env")
		helper.PanicIfError(err)

		secretKey := os.Getenv("SECRET_KEY")
		secretKeyByte := []byte(secretKey)

		token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
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

		var id string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if val, exists := claims["id"]; exists {
				if strVal, ok := val.(string); ok {
					id = strVal
				}
			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusUnauthorized)

				webResponse := web.WebResponse{
					Code:   http.StatusUnauthorized,
					Status: "Unauthorized",
					Data:   "Invalid Token",
				}
				helper.WriteToResponseBody(writer, webResponse)
				return
			}
		}

		log.Printf("User with uuid: %s enter Middleware", id)
		ctx := context.WithValue(request.Context(), userUUIDkey, id)
		next(writer, request.WithContext(ctx), p)
	}
}

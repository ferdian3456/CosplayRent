package middleware

import (
	"context"
	"cosplayrent/internal/helper"
	"cosplayrent/internal/model/web"
	"cosplayrent/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

const (
	userUUIDkey = "user_uuid"
)

type AuthMiddleware struct {
	Handler     http.Handler
	Log         *zerolog.Logger
	Config      *koanf.Koanf
	UserUsecase *usecase.UserUsecase
}

func NewAuthMiddleware(handler http.Handler, zerolog *zerolog.Logger, koanf *koanf.Koanf, userUsecase *usecase.UserUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		Handler:     handler,
		Log:         zerolog,
		Config:      koanf,
		UserUsecase: userUsecase,
	}
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

			middleware.Log.Warn().Msg("Unauthorized, no token provided")
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		splitToken := strings.Split(headerToken, "Bearer ")
		if len(splitToken) != 2 {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)

			webResponse := web.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Data:   "Token format is not match",
			}

			middleware.Log.Warn().Msg("Unauthorized, token format is not match")
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		secretKey := middleware.Config.String("application.secret_key")
		secretKeyByte := []byte(secretKey)

		token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
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

			middleware.Log.Warn().Msg("Unauthorized, invalid token")
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

				middleware.Log.Warn().Msg("Unauthorized, invalid token")
				helper.WriteToResponseBody(writer, webResponse)
				return
			}
		}

		err = middleware.UserUsecase.CheckUserExistance(request.Context(), id)
		if err != nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)

			webResponse := web.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "Unauthorized",
				Data:   "User not found, please register",
			}

			middleware.Log.Warn().Msg("Unauthorized, user not found")
			helper.WriteToResponseBody(writer, webResponse)
			return
		}

		middleware.Log.Debug().Msg("User:" + id)
		ctx := context.WithValue(request.Context(), userUUIDkey, id)
		next(writer, request.WithContext(ctx), p)
	}
}

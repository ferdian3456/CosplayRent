package config

import (
	"cosplayrent/internal/delivery/http"
	"cosplayrent/internal/delivery/http/middleware"
	"cosplayrent/internal/delivery/http/route"
	"cosplayrent/internal/repository"
	"cosplayrent/internal/usecase"
	"database/sql"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type ServerConfig struct {
	Router   *httprouter.Router
	DB       *sql.DB
	Memcache *memcache.Client
	Log      *zerolog.Logger
	Validate *validator.Validate
	Config   *koanf.Koanf
}

func Server(config *ServerConfig) {
	userRepository := repository.NewUserRepository(config.Log)
	userUsecase := usecase.NewUserUsecase(userRepository, config.DB, config.Validate, config.Log, config.Config)
	userController := controller.NewUserController(userUsecase, config.Log)

	authMiddleware := middleware.NewAuthMiddleware(config.Router, config.Log, config.Config, userUsecase)

	routeConfig := route.RouteConfig{
		Router:         config.Router,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.SetupRoute()
}

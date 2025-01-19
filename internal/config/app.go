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
	notificationRepository := repository.NewNotificationRepository(config.Log)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepository, config.DB, config.Validate, config.Log, config.Config)

	userRepository := repository.NewUserRepository(config.Log)
	userUsecase := usecase.NewUserUsecase(userRepository, repository.NewCostumeRepository(config.Log), notificationUsecase, config.DB, config.Validate, config.Log, config.Config)
	userController := controller.NewUserController(userUsecase, config.Log)

	costumeRepository := repository.NewCostumeRepository(config.Log)
	costumeUsecase := usecase.NewCostumeUsecase(userRepository, costumeRepository, repository.NewCategoryRepository(config.Log), config.DB, config.Validate, config.Log, config.Config)
	costumeController := controller.NewCostumeController(costumeUsecase, config.Log)

	categoryRepository := repository.NewCategoryRepository(config.Log)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository, config.DB, config.Validate, config.Log, config.Config)
	categoryController := controller.NewCategoryController(categoryUsecase, config.Log)

	wishlistRepository := repository.NewWishlistRepository(config.Log)
	wishlistUsecase := usecase.NewWishlistUsecase(wishlistRepository, costumeRepository, config.DB, config.Validate, config.Log, config.Config)
	wishlistController := controller.NewWishlistController(wishlistUsecase, config.Log)

	midtransUsecase := usecase.NewMidtransUsecase(userRepository, repository.NewOrderRepository(config.Log), repository.NewTopUpOrderRepository(config.Log), config.DB, config.Validate, config.Log, config.Config)
	midtransController := controller.NewMidtransController(midtransUsecase, config.Log)

	topUpOrderRepository := repository.NewTopUpOrderRepository(config.Log)
	topUpOrderUsecase := usecase.NewTopUpOrderUsecase(userRepository, topUpOrderRepository, midtransUsecase, config.DB, config.Validate, config.Log, config.Config)
	topUpOrderController := controller.NewTopUpOrderController(topUpOrderUsecase, config.Log)

	orderRepository := repository.NewOrderRepository(config.Log)
	orderUsecase := usecase.NewOrderUsecase(userRepository, costumeRepository, categoryRepository, orderRepository, midtransUsecase, config.DB, config.Validate, config.Log, config.Config)
	orderController := controller.NewOrderController(orderUsecase, config.Log)

	reviewRepository := repository.NewReviewRepository(config.Log)
	reviewUsecase := usecase.NewReviewUsecase(userRepository, costumeRepository, reviewRepository, config.DB, config.Validate, config.Log, config.Config)
	reviewController := controller.NewReviewController(reviewUsecase, orderUsecase, config.Log)

	rajaongkirUsecase := usecase.NewRajaOngkirUsecase(config.Memcache, config.Validate, config.Log, config.Config)
	rajaongkirController := controller.NewRajaOngkirController(rajaongkirUsecase, config.Log)

	authMiddleware := middleware.NewAuthMiddleware(config.Router, config.Log, config.Config, userUsecase)

	routeConfig := route.RouteConfig{
		Router:               config.Router,
		UserController:       userController,
		CostumeController:    costumeController,
		CategoryController:   categoryController,
		WishlistController:   wishlistController,
		OrderController:      orderController,
		ReviewController:     reviewController,
		TopUpOrderController: topUpOrderController,
		MidtransController:   midtransController,
		RajaOngkirController: rajaongkirController,
		AuthMiddleware:       authMiddleware,
	}

	routeConfig.SetupRoute()
}

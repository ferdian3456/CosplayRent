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
	DB       *sql.DB
	Memcache *memcache.Client
	Log      *zerolog.Logger
	Validate *validator.Validate
	Config   *koanf.Koanf
}

func Server(config *ServerConfig) {
	userRepository := repository.NewUserRepository(config.Log)
	UserUsecase := usecase.NewUserUsecase(userRepository, config.DB, config.Validate, config.Log)
	userController := controller.NewUserController(UserUsecase, config.Log)

	costumerRepository := repository.NewCostumeRepository(config.Log)
	CostumeUsecase := usecase.NewCostumeUsecase(costumerRepository, config.DB, config.Validate, config.Log)
	costumeController := controller.NewCostumeController(CostumeUsecase, config.Log)

	reviewRepository := repository.NewReviewRepository(config.Log)
	ReviewUsecase := usecase.NewReviewUsecase(reviewRepository, config.DB, config.Validate, config.Log)
	reviewController := controller.NewReviewController(ReviewUsecase, config.Log)

	orderRepository := repository.NewOrderRepository(config.Log)
	OrderUsecase := usecase.NewOrderUsecase(orderRepository, config.DB, config.Validate, config.Log)
	orderController := controller.NewOrderController(OrderUsecase, config.Log)

	topuporderRepository := repository.NewTopUpOrderRepository(config.Log)
	TopUpOrderUsecase := usecase.NewTopUpOrderUsecase(topuporderRepository, config.DB, config.Validate, config.Log)
	topuporderController := controller.NewTopUpOrderController(TopUpOrderUsecase, config.Log)

	wishlistRepository := repository.NewWishlistRepository(config.Log)
	WishlistUsecase := usecase.NewWishlistUsecase(wishlistRepository, config.DB, config.Validate, config.Log)
	wishlistController := controller.NewWishlistController(WishlistUsecase, config.Log)

	midtransRepository := repository.NewMidtransRepository(config.Log)
	MidtransUsecase := usecase.NewMidtransUsecase(midtransRepository, config.DB, config.Validate, config.Log)
	midtransController := controller.NewMidtransController(MidtransUsecase, config.Log)

	RajaOngkirUsecase := usecase.NewRajaOngkirUsecase(config.Validate, config.Memcache, config.Log)
	rajaongkirController := controller.NewRajaOngkirController(RajaOngkirUsecase, config.Log)

	router := httprouter.New()
	authMiddleware := middleware.NewAuthMiddleware(router)

	routeConfig := route.RouteConfig{
		Router:               router,
		UserController:       userController,
		CostumeController:    costumeController,
		ReviewController:     reviewController,
		OrderController:      orderController,
		TopUpOrderController: topuporderController,
		WishlistController:   wishlistController,
		MidtransController:   midtransController,
		RajaOngkirController: rajaongkirController,
		AuthMiddleware:       authMiddleware,
	}

	routeConfig.SetupRoute()
}

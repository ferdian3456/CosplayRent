package main

import (
	"cosplayrent/app"
	costume_controller "cosplayrent/controller/costume"
	midtrans_controller "cosplayrent/controller/midtrans"
	order_controller "cosplayrent/controller/order"
	rajaongkir_controller "cosplayrent/controller/rajaongkir"
	review_controller "cosplayrent/controller/review"
	topup_order_controller "cosplayrent/controller/topup_order"
	user_controller "cosplayrent/controller/user"
	wishlist_controller "cosplayrent/controller/wishlist"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/middleware"
	costume_repository "cosplayrent/repository/costume"
	midtrans_repository "cosplayrent/repository/midtrans"
	order_repository "cosplayrent/repository/order"
	review_repository "cosplayrent/repository/review"
	"cosplayrent/repository/topup_order"
	topup_order_repository "cosplayrent/repository/topup_order"
	user_repository "cosplayrent/repository/user"
	wishlist_repository "cosplayrent/repository/wishlist"
	costume_service "cosplayrent/service/costume"
	midtrans_service "cosplayrent/service/midtrans"
	order_service "cosplayrent/service/order"
	rajaongkir_service "cosplayrent/service/rajaongkir"
	review_service "cosplayrent/service/review"
	topup_order_service "cosplayrent/service/topup_order"
	user_service "cosplayrent/service/user"
	wishlist_service "cosplayrent/service/wishlist"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "True")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	DB := app.NewDB()
	memcacheClient := app.NewClient()
	validate := validator.New()
	userRepository := user_repository.NewUserRepository()
	userService := user_service.NewUserService(userRepository, DB, validate)
	userController := user_controller.NewUserController(userService)

	costumerRepository := costume_repository.NewCostumeRepository()
	costumeService := costume_service.NewCostumeService(costumerRepository, user_repository.NewUserRepository(), DB, validate)
	costumeController := costume_controller.NewCostumeController(costumeService)

	reviewRepository := review_repository.NewReviewRepository()
	reviewService := review_service.NewReviewService(reviewRepository, costume_repository.NewCostumeRepository(), user_repository.NewUserRepository(), DB, validate)
	reviewController := review_controller.NewReviewController(reviewService)

	orderRepository := order_repository.NewOrderRepository()
	orderService := order_service.NewOrderService(orderRepository, user_repository.NewUserRepository(), midtrans_service.NewMidtransService(midtrans_repository.NewMidtransRepository(), topup_order.NewTopUpOrderRepository(), user_repository.NewUserRepository(), order_repository.NewOrderRepository(), DB, validate), DB, validate)
	orderController := order_controller.NewOrderController(orderService)
	//log.Println(orderController)

	topuporderRepository := topup_order_repository.NewTopUpOrderRepository()
	topuporderService := topup_order_service.NewTopUpOrderService(topuporderRepository, user_repository.NewUserRepository(), midtrans_service.NewMidtransService(midtrans_repository.NewMidtransRepository(), topup_order.NewTopUpOrderRepository(), user_repository.NewUserRepository(), order_repository.NewOrderRepository(), DB, validate), DB, validate)
	topuporderController := topup_order_controller.NewTopUpOrderController(topuporderService)

	wishlistRepository := wishlist_repository.NewWishListRepository()
	wishlistService := wishlist_service.NewWishlistService(wishlistRepository, user_repository.NewUserRepository(), costume_repository.NewCostumeRepository(), DB, validate)
	wishlistController := wishlist_controller.NewWishlistController(wishlistService)

	midtransRepository := midtrans_repository.NewMidtransRepository()
	midtransService := midtrans_service.NewMidtransService(midtransRepository, topup_order.NewTopUpOrderRepository(), user_repository.NewUserRepository(), order_repository.NewOrderRepository(), DB, validate)
	midtransController := midtrans_controller.NewMidtransController(midtransService)

	rajaongkirService := rajaongkir_service.NewRajaOngkirService(validate, memcacheClient)
	rajaongkirController := rajaongkir_controller.NewRajaOngkirController(rajaongkirService)

	router := httprouter.New()
	authMiddleware := middleware.NewAuthMiddleware(router)

	router.POST("/api/register", userController.Register)
	router.GET("/api/identitycard", authMiddleware.ServeHTTP(userController.GetIdentityCard))
	router.POST("/api/identitycard", authMiddleware.ServeHTTP(userController.AddIdentityCard))
	router.PUT("/api/identitycard", authMiddleware.ServeHTTP(userController.UpdateIdentityCard))
	router.GET("/api/emoney", authMiddleware.ServeHTTP(userController.GetEMoneyAmount))
	//router.PUT("/api/emoney", authMiddleware.ServeHTTP(userController.TopUp))
	router.POST("/api/login", userController.Login)
	router.GET("/api/userdetail", authMiddleware.ServeHTTP(userController.FindByUUID))
	router.GET("/api/user", authMiddleware.ServeHTTP(userController.FindAll))
	router.PUT("/api/userdetail", authMiddleware.ServeHTTP(userController.Update))
	router.DELETE("/api/useraccount", authMiddleware.ServeHTTP(userController.Delete))

	//router.GET("/api/search/:costumeName", authMiddleware.ServeHTTP(costumeController.FindByName))
	router.POST("/api/costume", authMiddleware.ServeHTTP(costumeController.Create))
	router.GET("/api/costume", costumeController.FindAll)
	router.GET("/api/seller", authMiddleware.ServeHTTP(costumeController.FindSellerCostume))
	router.GET("/api/costume/:costumeID", costumeController.FindById)
	router.GET("/api/seller/:costumeID", authMiddleware.ServeHTTP(costumeController.FindSellerCostumeByCostumeID)) // find by costume id
	router.PUT("/api/seller/:costumeID", authMiddleware.ServeHTTP(costumeController.Update))
	router.DELETE("/api/seller/:costumeID", authMiddleware.ServeHTTP(costumeController.Delete))

	router.GET("/api/review", authMiddleware.ServeHTTP(reviewController.FindUserReview))
	router.POST("/api/review", authMiddleware.ServeHTTP(reviewController.Create))
	router.GET("/api/review/:reviewID", authMiddleware.ServeHTTP(reviewController.FindUserReviewByReviewID))
	router.PUT("/api/review/:reviewID", authMiddleware.ServeHTTP(reviewController.Update))
	router.DELETE("/api/review/:reviewID", authMiddleware.ServeHTTP(reviewController.DeleteUserReviewByReviewID))
	router.GET("/api/costume/:costumeID/review", reviewController.FindByCostumeId)

	router.POST("/api/order", authMiddleware.ServeHTTP(orderController.Create))
	router.POST("/api/order/midtrans", authMiddleware.ServeHTTP(orderController.DirectlyOrderToMidtrans))
	router.GET("/api/checkorder/:orderID", orderController.CheckStatusPayment)

	router.PUT("/api/topup", authMiddleware.ServeHTTP(topuporderController.CreateTopUpOrder))
	//router.GET("/api/check/topuporder", topup_order.n)

	router.GET("/api/wishlist", authMiddleware.ServeHTTP(wishlistController.FindAllWishListByUserId))
	router.POST("/api/wishlist/:costumeID", authMiddleware.ServeHTTP(wishlistController.AddWishlist))
	router.DELETE("/api/wishlist/:costumeID", authMiddleware.ServeHTTP(wishlistController.DeleteWishlist))

	router.GET("/api/provinces", rajaongkirController.FindProvince)
	router.GET("/api/city/:provinceID", rajaongkirController.FindCity)
	router.POST("/api/checkshippment", rajaongkirController.CheckShippment)

	router.POST("/api/midtrans/transaction/:orderID", authMiddleware.ServeHTTP(midtransController.CreateTransaction))
	router.POST("/api/midtrans/callback", midtransController.MidtransCallBack)

	router.ServeFiles("/static/*filepath", http.Dir("../static"))
	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:8081",
		Handler: CORS(router),
	}

	err := server.ListenAndServe()
	log.Println(err)
	helper.PanicIfError(err)
}

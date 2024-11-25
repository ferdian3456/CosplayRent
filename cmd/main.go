package main

import (
	"cosplayrent/app"
	costume_controller "cosplayrent/controller/costume"
	midtrans_controller "cosplayrent/controller/midtrans"
	order_controller "cosplayrent/controller/order"
	rajaongkir_controller "cosplayrent/controller/rajaongkir"
	review_controller "cosplayrent/controller/review"
	user_controller "cosplayrent/controller/user"
	"cosplayrent/exception"
	"cosplayrent/helper"
	"cosplayrent/middleware"
	costume_repository "cosplayrent/repository/costume"
	midtrans_repository "cosplayrent/repository/midtrans"
	order_repository "cosplayrent/repository/order"
	review_repository "cosplayrent/repository/review"
	user_repository "cosplayrent/repository/user"
	costume_service "cosplayrent/service/costume"
	midtrans_service "cosplayrent/service/midtrans"
	order_service "cosplayrent/service/order"
	rajaongkir_service "cosplayrent/service/rajaongkir"
	review_service "cosplayrent/service/review"
	user_service "cosplayrent/service/user"
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
	orderService := order_service.NewOrderService(orderRepository, DB, validate)
	orderController := order_controller.NewOrderController(orderService)
	//log.Println(orderController)

	midtransRepository := midtrans_repository.NewMidtransRepository()
	midtransService := midtrans_service.NewMidtransService(midtransRepository, user_repository.NewUserRepository(), order_repository.NewOrderRepository(), DB, validate)
	midtransController := midtrans_controller.NewMidtransController(midtransService)

	rajaongkirService := rajaongkir_service.NewRajaOngkirService(validate, memcacheClient)
	rajaongkirController := rajaongkir_controller.NewRajaOngkirController(rajaongkirService)

	router := httprouter.New()
	authMiddleware := middleware.NewAuthMiddleware(router)

	router.POST("/api/register", userController.Register)
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

	router.POST("/api/order", orderController.Create)

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

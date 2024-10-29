package main

import (
	"cosplayrent/app"
	costume_controller "cosplayrent/controller/costume"
	midtrans_controller "cosplayrent/controller/midtrans"
	order_controller "cosplayrent/controller/order"
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
	review_service "cosplayrent/service/review"
	user_service "cosplayrent/service/user"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Content-Type, Authorization")
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
	validate := validator.New()
	userRepository := user_repository.NewUserRepository()
	userService := user_service.NewUserService(userRepository, DB, validate)
	userController := user_controller.NewUserController(userService)

	costumerRepository := costume_repository.NewCostumeRepository()
	costumeService := costume_service.NewCostumeService(costumerRepository, DB, validate)
	costumeController := costume_controller.NewCostumeController(costumeService)

	reviewRepository := review_repository.NewReviewRepository()
	reviewService := review_service.NewReviewService(reviewRepository, DB, validate)
	reviewController := review_controller.NewReviewController(reviewService)

	orderRepository := order_repository.NewOrderRepository()
	orderService := order_service.NewOrderService(orderRepository, DB, validate)
	orderController := order_controller.NewOrderController(orderService)
	//log.Println(orderController)

	midtransRepository := midtrans_repository.NewMidtransRepository()
	midtransService := midtrans_service.NewMidtransService(midtransRepository, DB, validate)
	midtransController := midtrans_controller.NewMidtransController(midtransService)

	router := httprouter.New()
	authMiddleware := middleware.NewAuthMiddleware(router)

	router.POST("/api/register", userController.Register)
	router.POST("/api/login", userController.Login)
	router.GET("/api/user/:userUUID", userController.FindByUUID)
	router.GET("/api/user", userController.FindAll)
	router.PUT("/api/user/:userUUID", userController.Update)
	router.DELETE("/api/user/:userUUID", userController.Delete)
	router.GET("/api/verifytoken", userController.VerifyAndRetrieve)

	router.GET("/api/search/:costumeName", costumeController.FindByName)
	router.POST("/api/costume", authMiddleware.ServeHTTP(costumeController.Create))
	router.GET("/api/costume", costumeController.FindAll)
	router.GET("/api/costume/:costumeID", costumeController.FindById)
	router.PUT("/api/costume/:costumeID", authMiddleware.ServeHTTP(costumeController.Update))
	router.DELETE("/api/costume/:costumeID", authMiddleware.ServeHTTP(costumeController.Delete))

	router.GET("/api/review/:costumeID", reviewController.FindByCostumeId)
	router.POST("/api/review/:costumeID", reviewController.Create)

	router.POST("/api/order", orderController.Create)

	router.POST("/api/midtrans/transaction/:orderID", midtransController.CreateTransaction)
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

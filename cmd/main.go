package main

import (
	"cosplayrent/app"
	user_controller "cosplayrent/controller/user"
	"cosplayrent/exception"
	"cosplayrent/helper"
	user_repository "cosplayrent/repository/user"
	user_service "cosplayrent/service/user"
	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	DB := app.NewDB()
	validate := validator.New()
	userRepository := user_repository.NewUserRepository()
	userService := user_service.NewUserService(userRepository, DB, validate)
	userController := user_controller.NewUserController(userService)

	router := httprouter.New()

	router.POST("/api/register", userController.Register)
	router.POST("/api/login", userController.Login)
	router.GET("/api/user/:userUUID", userController.FindByUUID)
	router.GET("/api/user", userController.FindAll)
	router.PUT("/api/user/:userUUID", userController.Update)
	router.DELETE("/api/user/:userUUID", userController.Delete)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}

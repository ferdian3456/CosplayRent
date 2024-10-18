package main

import (
	"cosplayrent/app"
	costume_controller "cosplayrent/controller/costume"
	user_controller "cosplayrent/controller/user"
	"cosplayrent/exception"
	"cosplayrent/helper"
	costume_repository "cosplayrent/repository/costume"
	user_repository "cosplayrent/repository/user"
	costume_service "cosplayrent/service/costume"
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

	costumerRepository := costume_repository.NewCostumeRepository()
	costumeService := costume_service.NewCostumeService(costumerRepository, DB, validate)
	costumeController := costume_controller.NewCostumeController(costumeService)

	router := httprouter.New()

	router.POST("/api/register", userController.Register)
	router.POST("/api/login", userController.Login)
	router.GET("/api/user/:userUUID", userController.FindByUUID)
	router.GET("/api/user", userController.FindAll)
	router.PUT("/api/user/:userUUID", userController.Update)
	router.DELETE("/api/user/:userUUID", userController.Delete)

	router.POST("/api/costume", costumeController.Create)
	router.GET("/api/costume", costumeController.FindAll)
	router.GET("/api/costume/:costumeID", costumeController.FindById)
	router.PUT("/api/costume/:costumeID", costumeController.Update)
	router.DELETE("/api/costume/:costumeID", costumeController.Delete)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:8081",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}

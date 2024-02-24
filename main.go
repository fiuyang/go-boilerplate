package main

import (
	"gin-boilerplate/config"
	"gin-boilerplate/controller"
	_ "gin-boilerplate/docs"
	"gin-boilerplate/helper"
	"gin-boilerplate/migrations"
	"gin-boilerplate/repository"
	"gin-boilerplate/router"
	"gin-boilerplate/service"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

// @title 	Boilerplate API
// @version	1.0
// @description A Boilerplate API in Go using Gin framework

// @host 	localhost:8000
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	//Database
	db := config.ConnectionDB(&loadConfig)

	//Validation
	validate := validator.New()
	_ = validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		return helper.ValidateUnique(db, fl)
	})

	//Migrations
	migrations.AutoMigrate(db)

	//Init Repository
	userRepository := repository.NewUserRepositoryImpl(db)
	tagRepository := repository.NewTagRepositoryImpl(db)

	//Init Service
	authService := service.NewAuthServiceImpl(userRepository, validate)
	tagService := service.NewTagServiceImpl(tagRepository, validate)
	userService := service.NewUserServiceImpl(userRepository, validate)

	//Init controller
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	tagController := controller.NewTagController(tagService)

	//Router
	routes := router.NewRouter(
		userRepository,
		authController,
		userController,
		tagController,
	)

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	error := server.ListenAndServe()
	helper.ErrorPanic(error)
}

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

// @host 	localhost:8888
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
	migrations.PerformMigrations(db)

	//Init Repository
	userRepository := repository.NewUsersRepositoryImpl(db)
	tagsRepository := repository.NewTagsRepositoryImpl(db)

	//Init Service
	authenticationService := service.NewAuthenticationServiceImpl(userRepository)
	tagsService := service.NewTagsServiceImpl(tagsRepository)
	usersService := service.NewUsersServiceImpl(userRepository)

	//Init controller
	authenticationController := controller.NewAuthenticationController(authenticationService, validate)
	usersController := controller.NewUsersController(usersService, validate)
	tagsController := controller.NewTagsController(tagsService, validate)

	//Router
	routes := router.NewRouter(
		userRepository,
		authenticationController,
		usersController,
		tagsController,
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

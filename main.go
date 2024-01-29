package main

import (
	_ "gin-boilerplate/docs"
	"gin-boilerplate/config"
	"gin-boilerplate/controller"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"
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
	validate := validator.New()

	if !db.Migrator().HasTable(&model.Users{}) {
		db.Table("users").AutoMigrate(&model.Users{})
	}

	if !db.Migrator().HasTable(&model.Tags{}) {
		db.Table("tags").AutoMigrate(&model.Tags{})
	}

	if !db.Migrator().HasTable(&model.PasswordResets{}) {
		db.Table("password_resets").AutoMigrate(&model.PasswordResets{})
	}
	
	// db.Table("users").AutoMigrate(&model.Users{})
	// db.Table("tags").AutoMigrate(&model.Tags{})

	//Init Repository
	userRepository := repository.NewUsersRepositoryImpl(db)
	tagsRepository := repository.NewTagsREpositoryImpl(db)

	//Init Service
	authenticationService := service.NewAuthenticationServiceImpl(userRepository)
	tagsService := service.NewTagsServiceImpl(tagsRepository)
	usersService := service.NewUsersServiceImpl(userRepository)

	//Init controller
	authenticationController := controller.NewAuthenticationController(authenticationService, validate)
	usersController := controller.NewUsersController(usersService)
	tagsController := controller.NewTagsController(tagsService, validate)

	routes := router.NewRouter(userRepository, authenticationController, usersController, tagsController)

	

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server_err := server.ListenAndServe()
	helper.ErrorPanic(server_err)
}

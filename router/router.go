package router

import (
	"gin-boilerplate/controller"
	"gin-boilerplate/middleware"
	"gin-boilerplate/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	userRepository repository.UsersRepository, 
	authenticationController *controller.AuthenticationController, 
	usersController *controller.UserController, 
	tagsController *controller.TagsController,
	 ) *gin.Engine {
	service := gin.Default()

	//add swagger docs
	service.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})

	router := service.Group("/api")
	authenticationRouter := router.Group("/auth")
	authenticationRouter.POST("/register", authenticationController.Register)
	authenticationRouter.POST("/login", authenticationController.Login)
	authenticationRouter.POST("/forgot-password", authenticationController.ForgotPassword)
	authenticationRouter.POST("/check-otp", authenticationController.CheckOtp)
	authenticationRouter.PATCH("/reset-password", authenticationController.ResetPassword)
	authenticationRouter.POST("/logout", middleware.JwtMiddleware(userRepository), authenticationController.Logout)

	usersRouter := router.Group("/users")
	usersRouter.GET("", middleware.JwtMiddleware(userRepository), usersController.GetUsers)
	usersRouter.POST("", middleware.JwtMiddleware(userRepository), usersController.Create)
	usersRouter.PATCH("/:userId", middleware.JwtMiddleware(userRepository), usersController.Update)
	usersRouter.GET("/:userId", middleware.JwtMiddleware(userRepository), usersController.FindById)
	usersRouter.POST("/bulk", middleware.JwtMiddleware(userRepository), usersController.BulkDelete)
	usersRouter.GET("/export", middleware.JwtMiddleware(userRepository), usersController.Export)
	usersRouter.POST("/import", middleware.JwtMiddleware(userRepository), usersController.Import)

	tagsRouter := router.Group("/tags")
	tagsRouter.GET("", middleware.JwtMiddleware(userRepository), tagsController.FindAll)
	tagsRouter.GET("/:tagId", middleware.JwtMiddleware(userRepository), tagsController.FindById)
	tagsRouter.POST("", middleware.JwtMiddleware(userRepository), tagsController.Create)
	tagsRouter.PATCH("/:tagId", middleware.JwtMiddleware(userRepository), tagsController.Update)
	tagsRouter.DELETE("/:tagId", middleware.JwtMiddleware(userRepository), tagsController.Delete)

	return service
}

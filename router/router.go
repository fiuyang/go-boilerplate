package router

import (
	"gin-boilerplate/controller"
	"gin-boilerplate/exception"
	"gin-boilerplate/middleware"
	"gin-boilerplate/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	userRepository repository.UserRepository,
	authController *controller.AuthController,
	userController *controller.UserController,
	tagController *controller.TagController,
) *gin.Engine {
	service := gin.Default()

	service.Use(gin.CustomRecovery(exception.ErrorHandlers))
	//add swagger docs
	service.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})

	router := service.Group("/api")
	authRouter := router.Group("/auth")
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Login)
	authRouter.POST("/forgot-password", authController.ForgotPassword)
	authRouter.POST("/check-otp", authController.CheckOtp)
	authRouter.PATCH("/reset-password", authController.ResetPassword)
	authRouter.POST("/logout", middleware.JwtMiddleware(userRepository), authController.Logout)
	authRouter.GET("/refresh-token", middleware.JwtMiddleware(userRepository), authController.RefreshToken)

	userRouter := router.Group("/users")
	userRouter.GET("", middleware.JwtMiddleware(userRepository), userController.GetUsers)
	userRouter.POST("", middleware.JwtMiddleware(userRepository), userController.Create)
	userRouter.PATCH("/:userId", middleware.JwtMiddleware(userRepository), userController.Update)
	userRouter.GET("/:userId", middleware.JwtMiddleware(userRepository), userController.FindById)
	userRouter.POST("/bulk", middleware.JwtMiddleware(userRepository), userController.BulkDelete)
	userRouter.GET("/export", middleware.JwtMiddleware(userRepository), userController.Export)
	userRouter.POST("/import", middleware.JwtMiddleware(userRepository), userController.Import)

	tagRouter := router.Group("/tags")
	tagRouter.GET("", middleware.JwtMiddleware(userRepository), tagController.FindAll)
	tagRouter.GET("/:tagId", middleware.JwtMiddleware(userRepository), tagController.FindById)
	tagRouter.POST("", middleware.JwtMiddleware(userRepository), tagController.Create)
	tagRouter.PATCH("/:tagId", middleware.JwtMiddleware(userRepository), tagController.Update)
	tagRouter.DELETE("/:tagId", middleware.JwtMiddleware(userRepository), tagController.Delete)

	return service
}

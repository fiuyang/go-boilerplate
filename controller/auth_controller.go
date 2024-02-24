package controller

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{
		authService: service,
	}

}

// CreateTags		godoc
// @Summary			Login
// @Description		Login.
// @Param			login body request.LoginRequest true "login"
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/login [post]
func (controller *AuthController) Login(ctx *gin.Context) {
	loginRequest := request.LoginRequest{}
	err := ctx.ShouldBindJSON(&loginRequest)
	helper.ErrorPanic(err)

	accessToken, err := controller.authService.Login(loginRequest)

	helper.ErrorPanic(err)
	
	resp := response.LoginResponse{
		TokenType:   "Bearer",
		AccessToken: accessToken,
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully log in!",
		Data:    resp,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			Register
// @Description		Register.
// @Param			register body request.CreateUsersRequest true "Register"
// @Produce			application/json
// @Tags			auth
// @Success			201 {object} response.Response{}
// @Router			/auth/register [post]
func (controller *AuthController) Register(ctx *gin.Context) {
	createUsersRequest := request.CreateUsersRequest{}
	err := ctx.ShouldBindJSON(&createUsersRequest)
	helper.ErrorPanic(err)

	controller.authService.Register(createUsersRequest)

	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Ok",
		Message: "Successfully created user!",
		Data:    nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)
}

// CreateTags		godoc
// @Summary			ForgotPassword
// @Description		ForgotPassword.
// @Param			forgot-password body request.ForgotPasswordRequest true "ForgotPassword"
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/forgot-password [post]
func (controller *AuthController) ForgotPassword(ctx *gin.Context) {
	forgotPasswordRequest := request.ForgotPasswordRequest{}
	err := ctx.ShouldBindJSON(&forgotPasswordRequest)
	helper.ErrorPanic(err)

	otp, err := controller.authService.ForgotPassword(forgotPasswordRequest)
	helper.ErrorPanic(err)

	otpStr, err := strconv.Atoi(otp)
	helper.ErrorPanic(err)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Forgot Password successfully",
		Data:    otpStr,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			CheckOtp
// @Description		CheckOtp.
// @Param			check-otp body request.CheckOtpRequest true "CheckOtp"
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/check-otp [post]
func (controller *AuthController) CheckOtp(ctx *gin.Context) {

	checkOtpRequest := request.CheckOtpRequest{}
	err := ctx.ShouldBindJSON(&checkOtpRequest)
	helper.ErrorPanic(err)

	controller.authService.CheckOtp(checkOtpRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Otp Valid",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			ResetPassword
// @Description		ResetPassword.
// @Param			reset-password body request.ResetPasswordRequest true "ResetPassword"
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/reset-password [patch]
func (controller *AuthController) ResetPassword(ctx *gin.Context) {

	otpStr := ctx.Query("otp")

	otp, err := strconv.Atoi(otpStr)
	helper.ErrorPanic(err)

	resetPasswordRequest := request.ResetPasswordRequest{}
	err = ctx.ShouldBindJSON(&resetPasswordRequest)
	helper.ErrorPanic(err)

	controller.authService.ResetPassword(otp, resetPasswordRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Reset Password successfully",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			Logout
// @Description		Logout.
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/logout [post]
// @Security        Bearer
func (controller *AuthController) Logout(ctx *gin.Context) {
	token := extractTokenFromRequest(ctx)

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty token"})
		return
	}

	err := controller.authService.Logout(token)
	helper.ErrorPanic(err)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Logout Successfully",
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			RefreshToken
// @Description		RefreshToken.
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/refresh-token [get]
// @Security        Bearer
func (controller *AuthController) RefreshToken(ctx *gin.Context) {
	token := extractTokenFromRequest(ctx)

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty token"})
		return
	}

	refreshToken, err := controller.authService.RefreshToken(token)
	
	helper.ErrorPanic(err)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "",
		Data:    refreshToken,
	}
	ctx.JSON(http.StatusOK, webResponse)
}

func extractTokenFromRequest(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if token != "" {
		return strings.TrimPrefix(token, "Bearer ")
	}
	return ""
}

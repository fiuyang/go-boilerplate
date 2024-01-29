package controller

import (
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"gin-boilerplate/helper"
	"gin-boilerplate/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthenticationController struct {
	authenticationService service.AuthenticationService
	Validate       *validator.Validate
}

func NewAuthenticationController(service service.AuthenticationService,  validate *validator.Validate) *AuthenticationController {
	return &AuthenticationController{
		authenticationService: service,
		Validate:       validate,
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
func (controller *AuthenticationController) Login(ctx *gin.Context) {
	loginRequest := request.LoginRequest{}
	err := ctx.ShouldBindJSON(&loginRequest)
	helper.ErrorPanic(err)
	
	
	if err := controller.Validate.Struct(loginRequest); err != nil {
		if helper.ValidationError(err, ctx, loginRequest) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	
	token, err_token := controller.authenticationService.Login(loginRequest)
	fmt.Println(err_token)
	if err_token != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "BadRequest",
			Message: "Invalid email or password",
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}
	
	resp := response.LoginResponse{
		TokenType: "Bearer",
		Token:     token,
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
// @Success			200 {object} response.Response{}
// @Router			/auth/register [post]
func (controller *AuthenticationController) Register(ctx *gin.Context) {
	createUsersRequest := request.CreateUsersRequest{}
	err := ctx.ShouldBindJSON(&createUsersRequest)
	helper.ErrorPanic(err)
	
	if err := controller.Validate.Struct(createUsersRequest); err != nil {
		if helper.ValidationError(err, ctx, createUsersRequest) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	controller.authenticationService.Register(createUsersRequest)
	
	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully created user!",
		Data:    nil,
	}
	
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			ForgotPassword
// @Description		ForgotPassword.
// @Param			forgot-password body request.ForgotPasswordRequest true "ForgotPassword"
// @Produce			application/json
// @Tags			auth
// @Success			200 {object} response.Response{}
// @Router			/auth/forgot-password [post]
func (controller *AuthenticationController) ForgotPassword(ctx *gin.Context) {
	
	forgotPasswordRequest := request.ForgotPasswordRequest{}
	err := ctx.ShouldBindJSON(&forgotPasswordRequest)
	helper.ErrorPanic(err)
	
	if err := controller.Validate.Struct(forgotPasswordRequest); err != nil {
		if helper.ValidationError(err, ctx, forgotPasswordRequest) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	otp, err := controller.authenticationService.ForgotPassword(forgotPasswordRequest)
	
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to process forgot password request",
			Data:    nil,
		}
		ctx.JSON(http.StatusInternalServerError, webResponse)
		return
	}
	
	otpStr, err := strconv.Atoi(otp)
	helper.ErrorPanic(err)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Forgot Password successfully",
		Data:   otpStr,
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
func (controller *AuthenticationController) CheckOtp(ctx *gin.Context) {
	
	checkOtpRequest := request.CheckOtpRequest{}
	err := ctx.ShouldBindJSON(&checkOtpRequest)
	helper.ErrorPanic(err)
	
	if err := controller.Validate.Struct(checkOtpRequest); err != nil {
		if helper.ValidationError(err, ctx, checkOtpRequest) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	controller.authenticationService.CheckOtp(checkOtpRequest)
		
	helper.ErrorPanic(err)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Otp Valid",
		Data:   nil,
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
func (controller *AuthenticationController) ResetPassword(ctx *gin.Context) {
	
	otpStr := ctx.Query("otp")
	
	otp, err := strconv.Atoi(otpStr)
	helper.ErrorPanic(err)
	
	resetPasswordRequest := request.ResetPasswordRequest{}
	err = ctx.ShouldBindJSON(&resetPasswordRequest)
	helper.ErrorPanic(err)
	
	if err := controller.Validate.Struct(resetPasswordRequest); err != nil {
		if helper.ValidationError(err, ctx, resetPasswordRequest) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	controller.authenticationService.ResetPassword(otp, resetPasswordRequest)
	
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
func (controller *AuthenticationController) Logout(ctx *gin.Context) {
	token := extractTokenFromRequest(ctx)

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty token"})
		return
	}

	err := controller.authenticationService.Logout(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Logout Successfully",
		Data:    nil,
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
package service

import (
	"errors"
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/data/request"
	"gin-boilerplate/exception"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"
	"gin-boilerplate/repository"
	"gin-boilerplate/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

type AuthServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewAuthServiceImpl(userRepository repository.UserRepository, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

// Login
func (service *AuthServiceImpl) Login(user request.LoginRequest) (accessToken string, err error) {
	err_validate := service.Validate.Struct(user)
	helper.ErrorPanic(err_validate)

	new_users, users_err := service.UserRepository.FindByEmail(user.Email)
	if users_err != nil {
		return "", errors.New("email or password is wrong")
	}

	config, _ := config.LoadConfig(".")

	verify_error := utils.VerifyPassword(new_users.Password, user.Password)
	if verify_error != nil {
		return "", errors.New("email or password is wrong")
	}

	// Generate Token
	accessToken, err_access := utils.GenerateToken(config.AccessTokenExpiresIn, new_users.Id, config.AccessTokenSecret)
	helper.ErrorPanic(err_access)

	return accessToken, nil
}

// Register
func (service *AuthServiceImpl) Register(user request.CreateUsersRequest) {
	err := service.Validate.Struct(user)
	helper.ErrorPanic(err)

	hashedPassword, err := utils.HashPassword(user.Password)
	helper.ErrorPanic(err)

	newUser := model.Users{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}
	service.UserRepository.Save(newUser)
}

// ForgotPassword
func (service *AuthServiceImpl) ForgotPassword(user request.ForgotPasswordRequest) (string, error) {
	err := service.Validate.Struct(user)
	helper.ErrorPanic(err)
	existingUser, err := service.UserRepository.FindByEmail(user.Email)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	otp := utils.GenerateOTP(4)
	if err != nil {
		return "", errors.New("failed to generate token otp")
	}

	passwordReset := model.PasswordResets{
		Email:     existingUser.Email,
		Otp:       otp,
		CreatedAt: time.Now().Add(time.Minute * 5),
	}

	service.UserRepository.SaveOtp(passwordReset)

	emailData := utils.EmailData{
		Otp:     otp,
		Email:   existingUser.Email,
		Subject: " Reset Password",
	}

	utils.SendEmail(&existingUser, &emailData, "resetPassword.html")

	return fmt.Sprintf("%d", otp), nil
}

// CheckOtp
func (service *AuthServiceImpl) CheckOtp(reset request.CheckOtpRequest) (string, error) {
	err := service.Validate.Struct(reset)
	helper.ErrorPanic(err)

	existingOtp, err := service.UserRepository.FindByOtp(reset.Otp)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if reset.Otp != existingOtp.Otp {
		return "", errors.New("invalid otp")
	}

	if time.Now().After(existingOtp.CreatedAt) {
		return "", errors.New("otp has expired")
	}

	return "Otp Valid", nil
}

// ResetPassword
func (service *AuthServiceImpl) ResetPassword(otp int, user request.ResetPasswordRequest) (string, error) {
    err := service.Validate.Struct(user)
	helper.ErrorPanic(err)

	existingOtp, err := service.UserRepository.FindByOtp(otp)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if otp != existingOtp.Otp {
		return "", errors.New("invalid otp")
	}

	if time.Now().After(existingOtp.CreatedAt) {
		return "", errors.New("otp has expired")
	}

	fmt.Println("test password:", user.Password)
	hashedPassword, err := utils.HashPassword(user.Password)
	helper.ErrorPanic(err)

	dataset := model.Users{
		Email:    existingOtp.Email,
		Password: hashedPassword,
	}

	fmt.Println("test user:", dataset)
	service.UserRepository.UpdateOtp(dataset)

	service.UserRepository.DeleteOtp(existingOtp.Otp)

	return "", nil
}

// Logout
func (service *AuthServiceImpl) Logout(token string) error {
	if token == "" {
		return errors.New("empty token")
	}

	err := utils.AddToBlacklist(token)
	helper.ErrorPanic(err)
	return nil
}

// RefreshToken
func (service *AuthServiceImpl) RefreshToken(token string) (refreshToken string, err error) {
	if token == "" {
		return "", errors.New("empty token")
	}

	config, _ := config.LoadConfig(".")

	refreshToken, err_refresh := utils.RefreshToken(token, config.RefreshTokenExpiresIn, config.AccessTokenSecret)
	if err_refresh != nil {
		fmt.Println("Failed to refresh token:", err_refresh)
	}

	fmt.Println("Refresh token:", refreshToken)
	return refreshToken, nil
}

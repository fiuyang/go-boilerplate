package service

import (
	"errors"
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/data/request"
	"gin-boilerplate/helper"
	"gin-boilerplate/utils"
	"gin-boilerplate/model"
	"gin-boilerplate/repository"
	"time"
)

type AuthenticationServiceImpl struct {
	UsersRepository repository.UsersRepository
}

func NewAuthenticationServiceImpl(usersRepository repository.UsersRepository) AuthenticationService {
	return &AuthenticationServiceImpl{
		UsersRepository: usersRepository,
	}
}

//Login
func (a *AuthenticationServiceImpl) Login(users request.LoginRequest) (string, error) {	
	new_users, users_err := a.UsersRepository.FindByEmail(users.Email)
	if users_err != nil {
		return "", errors.New("invalid email or Password")
	}
	
	config, _ := config.LoadConfig(".")
	
	verify_error := utils.VerifyPassword(new_users.Password, users.Password)
	if verify_error != nil {
		return "", errors.New("invalid email or Password")
	}
	
	// Generate Token
	token, err_token := utils.GenerateToken(config.TokenExpiresIn, new_users.Id, config.TokenSecret)
	helper.ErrorPanic(err_token)
	return token, nil
	
}

//Register
func (a *AuthenticationServiceImpl) Register(users request.CreateUsersRequest) {
	
	hashedPassword, err := utils.HashPassword(users.Password)
	helper.ErrorPanic(err)
	
	newUser := model.Users{
		Username: users.Username,
		Email:    users.Email,
		Password: hashedPassword,
	}
	a.UsersRepository.Save(newUser)
}

// ForgotPassword
func (a *AuthenticationServiceImpl) ForgotPassword(users request.ForgotPasswordRequest) (string, error) {
	existingUser, err := a.UsersRepository.FindByEmail(users.Email)
	if err != nil {
		return "", errors.New("Email not found")
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

	a.UsersRepository.SaveOtp(passwordReset)
	
	emailData := utils.EmailData{
		Otp: otp,
		Email: existingUser.Email,
		Subject: " Reset Password",
	}
	
	utils.SendEmail(&existingUser, &emailData, "resetPassword.html")
	
	return fmt.Sprintf("%d", otp), nil
}

// CheckOtp
func (a *AuthenticationServiceImpl) CheckOtp(resets request.CheckOtpRequest) (string, error) {
	
	existingOtp, err := a.UsersRepository.FindByOtp(resets.Otp)
	if err != nil {
		return "", errors.New("Otp not found")
	}
	
	if resets.Otp != existingOtp.Otp {
		return "", errors.New("Invalid OTP")
	}
	
	if time.Now().After(existingOtp.CreatedAt) {
		return "", errors.New("OTP has expired")
	}
	
	return "Otp Valid", nil
}

//ResetPassword
func (a *AuthenticationServiceImpl) ResetPassword(otp int, users request.ResetPasswordRequest) (string, error) {
	
	existingOtp, err := a.UsersRepository.FindByOtp(otp)
	if err != nil {
		return "", errors.New("Otp not found")
	}

	if otp != existingOtp.Otp {
		return "", errors.New("Invalid OTP")
	}
	
	if time.Now().After(existingOtp.CreatedAt) {
		return "", errors.New("OTP has expired")
	}

	fmt.Println("test password:", users.Password)
	hashedPassword, err := utils.HashPassword(users.Password)
	helper.ErrorPanic(err)
	
	
	user := model.Users{
		Email: existingOtp.Email,
		Password: hashedPassword,
	}
	
	fmt.Println("test user:", user)
	a.UsersRepository.UpdateOtp(user)

	a.UsersRepository.DeleteOtp(existingOtp.Otp)

	return "", nil
}

// Logout
func (a *AuthenticationServiceImpl) Logout(token string) error {
	if token == "" {
		return errors.New("empty token")
	}

	err := utils.AddToBlacklist(token)
	helper.ErrorPanic(err)
	return nil
}

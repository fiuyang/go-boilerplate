package service

import "gin-boilerplate/data/request"

type AuthService interface {
	Login(user request.LoginRequest) (accessToken string, err error)
	Register(user request.CreateUsersRequest)
	Logout(token string) error
	ForgotPassword(user request.ForgotPasswordRequest) (string, error)
	CheckOtp(reset request.CheckOtpRequest) (string, error)
	ResetPassword(otp int, user request.ResetPasswordRequest) (string, error)
	RefreshToken(token string) (refreshToken string, err error)
}

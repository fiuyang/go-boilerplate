package service

import "gin-boilerplate/data/request"

type AuthenticationService interface {
	Login(users request.LoginRequest) (string, error)
	Register(users request.CreateUsersRequest)
	Logout(token string) (error)
	ForgotPassword(users request.ForgotPasswordRequest) (string, error)
	CheckOtp(resets request.CheckOtpRequest) (string, error)
	ResetPassword(otp int, users request.ResetPasswordRequest) (string, error)
}

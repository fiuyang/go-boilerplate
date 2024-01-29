package request

type CreateUsersRequest struct {
	Username string `validate:"required" json:"username"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8,max=100" json:"password"`
}

type UpdateUsersRequest struct {
	Id       int    `validate:"required"`
	Username string `validate:"required,max=200,min=2" json:"username"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8,max=100" json:"password"`
}

type LoginRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=2,max=100" json:"password"`
}

type ForgotPasswordRequest struct {
	Email              string `validate:"required,email" json:"email"`
}

type CheckOtpRequest struct {
	Otp              int `validate:"required" json:"otp"`
}

type ResetPasswordRequest struct {
	Password             string `validate:"required" json:"password"`
	PasswordConfirmation string `validate:"required" json:"password_confirmation"`
}
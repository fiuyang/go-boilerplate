package repository

import "gin-boilerplate/model"

type UserRepository interface {
	Save(user model.Users)
	Update(user model.Users)
	Delete(userId int)
	BulkDelete(userIds []int)
	FindById(userId int) (model.Users, error)
	FindAll(filters map[string]string) []model.Users
	FindByUsername(username string) (model.Users, error)
	FindByEmail(email string) (model.Users, error)
	UpdateOtp(user model.Users)
	SaveOtp(reset model.PasswordResets)
	DeleteOtp(otp int)
	FindByOtp(otp int) (model.PasswordResets, error)
}

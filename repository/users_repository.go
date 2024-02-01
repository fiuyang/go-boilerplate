package repository

import "gin-boilerplate/model"

type UsersRepository interface {
	Save(users model.Users)
	Update(users model.Users)
	Delete(usersId int)
	BulkDelete(userIds []int)
	FindById(usersId int) (model.Users, error)
	FindAll(filters map[string]string) []model.Users
	FindByUsername(username string) (model.Users, error)
	FindByEmail(email string) (model.Users, error)
	UpdateOtp(users model.Users)
	SaveOtp(resets model.PasswordResets)
	DeleteOtp(otp int)
	FindByOtp(otp int) (model.PasswordResets, error)
}

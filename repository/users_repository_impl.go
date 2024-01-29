package repository

import (
	"errors"
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"

	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

func (u *UsersRepositoryImpl) Delete(usersId int) {
	var users model.Users
	result := u.Db.Where("id = ?", usersId).Delete(&users)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) FindAll(filters map[string]string) []model.Users {
	var users []model.Users
    query := u.Db.Model(&users)

    for field, value := range filters {
		if value != "" {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	results := query.Find(&users)
	helper.ErrorPanic(results.Error)
	return users
}

func (u *UsersRepositoryImpl) FindById(usersId int) (model.Users, error) {
	var users model.Users
	result := u.Db.Find(&users, usersId)
	if result != nil {
		return users, nil
	} else {
		return users, errors.New("users is not found")
	}
}

func (u *UsersRepositoryImpl) Save(users model.Users) {
	result := u.Db.Create(&users)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) Update(users model.Users) {
	var updateUsers = request.UpdateUsersRequest{
		Id:       users.Id,
		Username: users.Username,
		Email:    users.Email,
		Password: users.Password,
	}
	result := u.Db.Model(&users).Updates(updateUsers)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) FindByUsername(username string) (model.Users, error) {
	var users model.Users
	result := u.Db.First(&users, "username = ?", username)

	if result.Error != nil {
		return users, errors.New("invalid username or Password")
	}
	return users, nil
}

func (u *UsersRepositoryImpl) FindByEmail(email string) (model.Users, error) {
	var users model.Users
	result := u.Db.First(&users, "email = ?", email)

	if result.RowsAffected == 0 {
		return users, errors.New("email not found")
	}

	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

func (u *UsersRepositoryImpl) UpdateOtp(users model.Users) {

	updateFields := map[string]interface{}{
        "Password": users.Password,
    }

    tx := u.Db.Begin()

    result := tx.Model(&users).
        Where("email = ?", users.Email).
        Updates(updateFields)

    if result.Error != nil {
		helper.ErrorPanic(result.Error)
        tx.Rollback()
        return
    }

    tx.Commit()
}

func (u *UsersRepositoryImpl) SaveOtp(resets model.PasswordResets) {
	result := u.Db.Create(&resets)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) FindByOtp(Otp int) (model.PasswordResets, error) {
	var resets model.PasswordResets
	result := u.Db.First(&resets, "otp = ?", Otp)

	if result.RowsAffected == 0 {
		return resets, errors.New("otp not found")
	}

	if result.Error != nil {
		return resets, result.Error
	}
	return resets, nil
}

func (u *UsersRepositoryImpl) DeleteOtp(otp int) {
	var resets model.PasswordResets
	result := u.Db.Where("otp = ?", otp).Delete(&resets)
	helper.ErrorPanic(result.Error)
}
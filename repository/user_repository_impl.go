package repository

import (
	"errors"
	"fmt"
	"gin-boilerplate/data/request"
	"gin-boilerplate/helper"
	"gin-boilerplate/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: Db}
}

func (repo *UserRepositoryImpl) Delete(userId int) {
	var user model.Users
	result := repo.Db.Where("id = ?", userId).Delete(&user)
	helper.ErrorPanic(result.Error)
}

func (repo *UserRepositoryImpl) BulkDelete(userIds []int) {
	var user model.Users
	result := repo.Db.Where("id IN (?)", userIds).Delete(&user)
	helper.ErrorPanic(result.Error)
}

func (repo *UserRepositoryImpl) FindAll(filters map[string]string) []model.Users {
	var users []model.Users
    query := repo.Db.Model(&users)

    for field, value := range filters {
		if value != "" {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	results := query.Find(&users)
	helper.ErrorPanic(results.Error)
	return users
}

func (repo *UserRepositoryImpl) FindById(userId int) (model.Users, error) {
	var user model.Users
	result := repo.Db.Find(&user, userId)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

func (repo *UserRepositoryImpl) Save(user model.Users) {
	result := repo.Db.Create(&user)
	helper.ErrorPanic(result.Error)
}

func (repo *UserRepositoryImpl) Update(user model.Users) {
	var dataset = request.UpdateUsersRequest{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	result := repo.Db.Model(&user).Updates(dataset)
	helper.ErrorPanic(result.Error)
}

func (repo *UserRepositoryImpl) FindByUsername(username string) (model.Users, error) {
	var user model.Users
	result := repo.Db.First(&user, "username = ?", username)

	if result.Error != nil {
		return user, errors.New("username or Password is wrong")
	}
	return user, nil
}

func (repo *UserRepositoryImpl) FindByEmail(email string) (model.Users, error) {
	var user model.Users
	result := repo.Db.First(&user, "email = ?", email)

	if result.RowsAffected == 0 {
		return user, errors.New("email not found")
	}

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (repo *UserRepositoryImpl) UpdateOtp(user model.Users) {

	updateFields := map[string]interface{}{
        "Password": user.Password,
    }

    tx := repo.Db.Begin()

    result := tx.Model(&user).
        Where("email = ?", user.Email).
        Updates(updateFields)

    if result.Error != nil {
		helper.ErrorPanic(result.Error)
        tx.Rollback()
        return
    }

    tx.Commit()
}

func (repo *UserRepositoryImpl) SaveOtp(reset model.PasswordResets) {
	result := repo.Db.Create(&reset)
	helper.ErrorPanic(result.Error)
}

func (repo *UserRepositoryImpl) FindByOtp(Otp int) (model.PasswordResets, error) {
	var reset model.PasswordResets
	result := repo.Db.First(&reset, "otp = ?", Otp)

	if result.RowsAffected == 0 {
		return reset, errors.New("otp not found")
	}

	if result.Error != nil {
		return reset, result.Error
	}
	return reset, nil
}

func (repo *UserRepositoryImpl) DeleteOtp(otp int) {
	var reset model.PasswordResets
	result := repo.Db.Where("otp = ?", otp).Delete(&reset)
	helper.ErrorPanic(result.Error)
}

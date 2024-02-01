package service

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"mime/multipart"
)

type UsersService interface {
	FindAll(filters map[string]string) []response.UsersResponse
	Create(users request.CreateUsersRequest)
	Update(users request.UpdateUsersRequest)
	// Delete(userId int)
	BulkDelete(userIds []int)
	FindById(userId int) response.UsersResponse
	Export() (string, error)
	Import(file *multipart.FileHeader) error
}

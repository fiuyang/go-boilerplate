package service

import (
	"gin-boilerplate/data/request"
	"gin-boilerplate/data/response"
	"mime/multipart"
)

type UserService interface {
	FindAll(filters map[string]string) []response.UsersResponse
	Create(user request.CreateUsersRequest)
	Update(user request.UpdateUsersRequest)
	BulkDelete(userIds []int) error 
	FindById(userId int) response.UsersResponse
	Export() (string, error)
	Import(file *multipart.FileHeader) error
}

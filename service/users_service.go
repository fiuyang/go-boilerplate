package service

import "gin-boilerplate/data/response"

type UsersService interface {
	FindAll(filters map[string]string) []response.UsersResponse
}

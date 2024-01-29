package service

import (
	"gin-boilerplate/data/response"
	"gin-boilerplate/repository"
)

type UsersServiceImpl struct {
	UsersRepository repository.UsersRepository
}

func NewUsersServiceImpl(usersRepository repository.UsersRepository) UsersService {
	return &UsersServiceImpl{
		UsersRepository: usersRepository,
	}
}

func (u *UsersServiceImpl) FindAll(filters map[string]string) []response.UsersResponse {
	result := u.UsersRepository.FindAll(filters)

	var users []response.UsersResponse
	for _, value := range result {
		user := response.UsersResponse{
			Id:       value.Id,
			Username: value.Username,
			Email:    value.Email,
		}
		users = append(users, user)
	}

	return users
}

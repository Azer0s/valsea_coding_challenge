package service

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/transactional"
)

type UserService interface {
	CreateUser(request transactional.CreateUserRequest) (*dto.UserDTO, error)
	GetUserById(userId string) (*dto.UserDTO, error)
	GetAllUsers() ([]dto.UserDTO, error)
}

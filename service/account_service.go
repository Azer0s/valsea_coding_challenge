package service

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/transactional"
)

type AccountService interface {
	CreateAccount(request transactional.CreateUserRequest) (*dto.UserDTO, error)
	GetAccountById(userId string) (*dto.UserDTO, error)
	GetAllAccounts() ([]dto.UserDTO, error)
}

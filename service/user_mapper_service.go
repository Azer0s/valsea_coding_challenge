package service

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
)

type UserMapperService interface {
	ToUserDto(userEntity *entity.UserEntity) *dto.UserDTO
}

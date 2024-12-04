package impl

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
)

type UserMapperServiceImpl struct {
}

func (s *UserMapperServiceImpl) ToUserDto(userEntity *entity.UserEntity) (*dto.UserDTO, error) {
	return &dto.UserDTO{
		Id:      userEntity.Id.String(),
		Name:    userEntity.Name,
		Balance: userEntity.Balance.InexactFloat64(),
	}, nil
}

func NewUserMapperServiceImpl() *UserMapperServiceImpl {
	return &UserMapperServiceImpl{}
}

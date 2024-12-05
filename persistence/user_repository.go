package persistence

import (
	"valsea_coding_challenge/domain/entity"
)

type UserRepository interface {
	FindById(id string) (*entity.UserEntity, error)
	FindAll() []*entity.UserEntity
	Save(user *entity.UserEntity) error
	Delete(user *entity.UserEntity) error
	Update(user *entity.UserEntity) error

	Lock(users []*entity.UserEntity) error
	Unlock(users []*entity.UserEntity) error
}

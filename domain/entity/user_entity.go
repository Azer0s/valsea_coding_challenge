package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserEntity struct {
	Id      uuid.UUID
	Name    string
	Balance decimal.Decimal
}

func NewUserEntity(name string, balance decimal.Decimal) *UserEntity {
	return &UserEntity{
		Id:      uuid.New(),
		Name:    name,
		Balance: balance,
	}
}

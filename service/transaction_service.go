package service

import (
	"github.com/shopspring/decimal"
	"valsea_coding_challenge/domain/dto"
)

type TransactionService interface {
	Transfer(fromUserId string, toUserId string, amount decimal.Decimal) ([]dto.TransactionDTO, error)
	Deposit(userId string, amount decimal.Decimal) (*dto.TransactionDTO, error)
	Withdraw(userId string, amount decimal.Decimal) (*dto.TransactionDTO, error)

	GetTransactions(accountID string) ([]dto.TransactionDTO, error)
}

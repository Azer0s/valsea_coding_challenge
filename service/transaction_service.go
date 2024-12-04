package service

import (
	"github.com/shopspring/decimal"
	"valsea_coding_challenge/domain/dto"
)

type TransactionService interface {
	Transfer(fromAccountID string, toAccountID string, amount decimal.Decimal) error
	Deposit(accountID string, amount decimal.Decimal) error
	Withdraw(accountID string, amount decimal.Decimal) error

	GetTransactions(accountID string) ([]dto.TransactionDTO, error)
}

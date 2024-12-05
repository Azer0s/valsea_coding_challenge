package persistence

import "valsea_coding_challenge/domain/entity"

type TransactionRepository interface {
	CreateHistory(userId string) error
	GetAllTransactionsForUser(userId string) ([]*entity.TransactionEntity, error)
	SaveForUser(userId string, transaction *entity.TransactionEntity) error
}

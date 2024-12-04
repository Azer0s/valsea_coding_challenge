package impl

import (
	"errors"
	"valsea_coding_challenge/domain/entity"
	"valsea_coding_challenge/util"
)

type InMemoryTransactionRepository struct {
	transactions *util.TypedSyncMap[string, *util.AtomicLinkedList[*entity.TransactionEntity]]
}

func (i *InMemoryTransactionRepository) CreateHistory(userId string) {
	i.transactions.Store(userId, util.NewAtomicLinkedList[*entity.TransactionEntity]())
}

func (i *InMemoryTransactionRepository) GetAllTransactionsForUser(userId string) ([]*entity.TransactionEntity, error) {
	history, ok := i.transactions.Load(userId)
	if !ok {
		return nil, errors.New("user not found")
	}

	return history.Slice(), nil
}

func (i *InMemoryTransactionRepository) SaveForUser(userId string, transaction *entity.TransactionEntity) error {
	history, ok := i.transactions.Load(userId)
	if !ok {
		return errors.New("user not found")
	}

	history.Add(transaction)
	return nil
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: util.NewTypedSyncMap[string, *util.AtomicLinkedList[*entity.TransactionEntity]](),
	}
}

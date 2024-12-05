package impl

import (
	"errors"
	"go.uber.org/zap"
	"valsea_coding_challenge/domain/entity"
	"valsea_coding_challenge/util"
)

type InMemoryTransactionRepository struct {
	log          *zap.Logger
	transactions *util.TypedSyncMap[string, *util.AtomicLinkedList[*entity.TransactionEntity]]
}

func (i *InMemoryTransactionRepository) CreateHistory(userId string) error {
	if _, ok := i.transactions.Load(userId); ok {
		return errors.New("transaction history already exists for user")
	}

	i.transactions.Store(userId, util.NewAtomicLinkedList[*entity.TransactionEntity]())
	i.log.Debug("Creating transaction history for user", zap.String("userId", userId))

	return nil
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

func (i *InMemoryTransactionRepository) DeleteTransaction(userId string, transactionId string) error {
	history, ok := i.transactions.Load(userId)
	if !ok {
		return errors.New("user not found")
	}

	transaction := *history.Find(func(t *entity.TransactionEntity) bool {
		return t.Id.String() == transactionId
	})
	if transaction == nil {
		return errors.New("transaction not found")
	}

	return history.Remove(transaction)
}

func NewInMemoryTransactionRepository(log *zap.Logger) *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: util.NewTypedSyncMap[string, *util.AtomicLinkedList[*entity.TransactionEntity]](),
		log:          log,
	}
}

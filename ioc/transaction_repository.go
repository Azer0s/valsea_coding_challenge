package ioc

import (
	"go.uber.org/zap"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/persistence/impl"
)

func ProvideTransactionRepository(log *zap.Logger) persistence.TransactionRepository {
	return impl.NewInMemoryTransactionRepository(log)
}

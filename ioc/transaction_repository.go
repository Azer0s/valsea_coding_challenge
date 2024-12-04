package ioc

import (
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/persistence/impl"
)

func ProvideTransactionRepository() persistence.TransactionRepository {
	return impl.NewInMemoryTransactionRepository()
}

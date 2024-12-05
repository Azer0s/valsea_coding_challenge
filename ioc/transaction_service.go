package ioc

import (
	"go.uber.org/zap"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/service"
	"valsea_coding_challenge/service/impl"
)

func ProvideTransactionService(log *zap.Logger, transactionRepository persistence.TransactionRepository, transactionMapperService service.TransactionMapperService, userRepository persistence.UserRepository) service.TransactionService {
	return impl.NewTransactionServiceImpl(log, transactionRepository, transactionMapperService, userRepository)
}

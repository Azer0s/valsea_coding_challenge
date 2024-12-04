package ioc

import (
	"go.uber.org/zap"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/service"
	"valsea_coding_challenge/service/impl"
)

func ProvideAccountService(userRepository persistence.UserRepository, transactionRepository persistence.TransactionRepository, userMapperService service.UserMapperService, log *zap.Logger) service.AccountService {
	return impl.NewAccountServiceImpl(userRepository, transactionRepository, userMapperService, log)
}

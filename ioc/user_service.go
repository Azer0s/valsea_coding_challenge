package ioc

import (
	"go.uber.org/zap"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/service"
	"valsea_coding_challenge/service/impl"
)

func ProvideUserService(userRepository persistence.UserRepository, transactionRepository persistence.TransactionRepository, userMapperService service.UserMapperService, log *zap.Logger) service.UserService {
	return impl.NewUserServiceImpl(userRepository, transactionRepository, userMapperService, log)
}

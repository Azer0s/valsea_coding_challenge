package ioc

import (
	"go.uber.org/zap"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/persistence/impl"
)

func ProvideUserRepository(log *zap.Logger) persistence.UserRepository {
	return impl.NewInMemoryUserRepository(log)
}

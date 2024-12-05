package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"log"
	"valsea_coding_challenge/controller"
	"valsea_coding_challenge/ioc"
	"valsea_coding_challenge/util"
)

func provide(container *dig.Container, constructor interface{}, name string) {
	err := container.Provide(constructor)
	if err != nil {
		log.Fatalf("Failed to provide %s", name)
	}
}

func StartServer(config *util.Config, starting chan<- struct{}) {
	container := dig.New()

	provide(container, func() *util.Config {
		if config == nil {
			config = &util.Config{
				Port:     8080,
				LogLevel: zap.InfoLevel,
			}
		}
		return config
	}, "config")
	provide(container, ioc.ProvideZap, "zap logger")
	provide(container, ioc.ProvideGinGonic, "gin engine")

	provide(container, ioc.ProvideUserMapperService, "user mapper service")
	provide(container, ioc.ProvideUserRepository, "user repository")

	provide(container, ioc.ProvideTransactionMapperService, "transaction mapper service")
	provide(container, ioc.ProvideTransactionRepository, "transaction repository")

	provide(container, ioc.ProvideUserService, "account service")
	provide(container, ioc.ProvideTransactionService, "transaction service")

	err := container.Invoke(controller.AccountController)
	if err != nil {
		log.Fatal("Failed to register routes for account controller", zap.Error(err))
	}

	err = container.Invoke(controller.TransactionController)
	if err != nil {
		log.Fatal("Failed to register routes for transaction controller", zap.Error(err))
	}

	// Run the application
	err = container.Invoke(func(config *util.Config, log *zap.Logger, r *gin.Engine) {
		log.Info(fmt.Sprintf("Starting server on port %d", config.Port))

		if starting != nil {
			starting <- struct{}{}
		}

		err = r.Run(fmt.Sprintf(":%d", config.Port))
		if err != nil {
			log.Fatal("Failed to run gin engine", zap.Error(err))
		}
	})
	if err != nil {
		log.Fatal("Failed to run gin engine", zap.Error(err))
	}
}

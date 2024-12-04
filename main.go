package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"log"
	"valsea_coding_challenge/controller"
	"valsea_coding_challenge/ioc"
)

func provide(container *dig.Container, constructor interface{}, name string) {
	err := container.Provide(constructor)
	if err != nil {
		log.Fatalf("Failed to provide %s", name)
	}
}

func main() {
	container := dig.New()

	provide(container, ioc.ProvideZap, "zap logger")
	provide(container, ioc.ProvideGinGonic, "gin engine")

	provide(container, ioc.ProvideUserMapperService, "user mapper service")
	provide(container, ioc.ProvideUserRepository, "user repository")

	provide(container, ioc.ProvideTransactionRepository, "transaction repository")

	provide(container, ioc.ProvideAccountService, "account service")

	err := container.Invoke(controller.AccountController)
	if err != nil {
		log.Fatal("Failed to register routes for account controller", zap.Error(err))
	}

	// Run the application
	err = container.Invoke(func(log *zap.Logger, r *gin.Engine) {
		log.Info("Starting server on port 8080")

		err = r.Run(":8080")
		if err != nil {
			log.Fatal("Failed to run gin engine", zap.Error(err))
		}
	})
	if err != nil {
		log.Fatal("Failed to run gin engine", zap.Error(err))
	}
}

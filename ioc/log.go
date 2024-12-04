package ioc

import "go.uber.org/zap"

func ProvideZap() *zap.Logger {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return log
}

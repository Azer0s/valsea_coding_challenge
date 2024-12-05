package ioc

import (
	"go.uber.org/zap"
	"valsea_coding_challenge/util"
)

func ProvideZap(config *util.Config) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(config.LogLevel)

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return log
}

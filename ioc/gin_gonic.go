package ioc

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func ProvideGinGonic(log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal("Failed to set trusted proxies", zap.Error(err))
	}

	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(log, true))
	return r
}

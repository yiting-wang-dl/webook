package ioc

import (
	"github.com/spf13/viper"
	"github.com/webook/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.LoggerV1 {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}

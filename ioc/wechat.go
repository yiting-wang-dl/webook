package ioc

import (
	"github.com/webook/internal/service/oauth2/wechat"
	"github.com/webook/pkg/logger"
	"os"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("cannot find environment variable WECHAT_APP_ID")
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("cannot find environment variable WECHAT_APP_SECRET")
	}
	return wechat.NewService(appID, appSecret, l)
}

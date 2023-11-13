//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/webook/internal/repository"
	"github.com/webook/internal/repository/cache"
	"github.com/webook/internal/repository/dao"
	"github.com/webook/internal/service"
	"github.com/webook/internal/web"
	ijwt "github.com/webook/internal/web/jwt"
	"github.com/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third party dependency
		ioc.InitRedis,
		ioc.InitDB,
		ioc.InitLogger,

		// DAO
		dao.NewUserDAO,
		//dao.NewArticleGORMDAO,

		// cache
		cache.NewCodeCache,
		cache.NewUserCache,
		//cache.NewArticleRedisCache,

		// repository
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,
		//repository.NewCachedArticleRepositry,

		// service
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		//service.NewArticleService,

		// handler
		web.NewUserHandler,
		//web.NewArticleHandler,
		ijwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

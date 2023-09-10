//go:build wireinject

package main

import (
	"boot-camp-homework/webook/internal/repository"
	"boot-camp-homework/webook/internal/repository/dao"
	"boot-camp-homework/webook/internal/service"
	"boot-camp-homework/webook/internal/web"
	"boot-camp-homework/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		//初始化dao
		dao.NewUserDao,

		ioc.InitUserCache,
		ioc.InitCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,

		web.NewUserHandler,
		//gin.Default,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

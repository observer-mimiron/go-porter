package svc

import (
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/conf"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/logger"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config conf.Config
	Logger *zap.Logger
	Db     mysql.Repo
	Redis  redis.Repo
}

func NewServiceContext(c conf.Config) *ServiceContext {
	log, err := logger.NewJSONLogger(c.Log)
	if err != nil {
		log.Fatal("new logger err", zap.Error(err))
	}
	mysqlClient, err := mysql.New(conf.Get().MySQL)
	if err != nil {
		log.Fatal("new mysql err", zap.Error(err))
	}

	redisClient, err := redis.New(conf.Get().Redis)
	if err != nil {
		log.Fatal("new redis err", zap.Error(err))
	}

	return &ServiceContext{
		Config: c,
		Logger: log,
		Db:     mysqlClient,
		Redis:  redisClient,
	}
}

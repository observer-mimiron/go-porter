package svc

import (
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/cache/redis"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/conf"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/database/mysql"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/logger"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config conf.Config
	Logger *zap.Logger
	Db     mysql.Repo
	Redis  redis.Repo
}

func NewServiceContext(c conf.Config) (*ServiceContext, error) {
	log, err := logger.NewJSONLogger(c.Log)
	if err != nil {
		return nil, err
	}

	mysqlClient, err := mysql.New(c.MySQL)
	if err != nil {
		log.Error("new mysql err", zap.Error(err))
		return nil, err
	}

	redisClient, err := redis.New(c.Redis)
	if err != nil {
		log.Error("new redis err", zap.Error(err))
		_ = mysqlClient.DbWClose()
		_ = mysqlClient.DbRClose()
		return nil, err
	}

	return &ServiceContext{
		Config: c,
		Logger: log,
		Db:     mysqlClient,
		Redis:  redisClient,
	}, nil
}

package svc

import (
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/conf"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/logger"
	"go-porter/pkg/core/pkg/net/httpx"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config conf.Config
	Logger *zap.Logger
	Mux    httpx.Mux
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

	//开启相关功能组
	mux, err := httpx.New(log,
		//httpx.WithDisablePrometheus(), // 关闭 prometheus
		//httpx.WithDisablePProf(),      // 关闭 WithDisablePProf
		httpx.WithEnableCors(), //跨域
		//httpx.WithEnableRate(),        //限流
		//core.WithAlertNotify(alert.NotifyHandler(logger)),     //报警
		//httpx.WithRecordMetrics(metrics.RecordHandler(log)), //监控
	)
	if err != nil {
		log.Fatal("new mux err", zap.Error(err))
	}

	return &ServiceContext{
		Config: c,
		Logger: log,
		Db:     mysqlClient,
		Redis:  redisClient,
		Mux:    mux,
	}
}

package router

import (
	"go-porter/internal/alert"
	"go-porter/internal/metrics"
	"go-porter/internal/pkg/core"
	"go-porter/internal/repository/mysql"
	"go-porter/internal/repository/redis"
	"go-porter/internal/router/interceptor"
	"go-porter/pkg/errors"
	"go.uber.org/zap"
)

type resource struct {
	mux          core.Mux
	logger       *zap.Logger
	db           mysql.Repo
	cache        redis.Repo
	interceptors interceptor.Interceptor
}

type Server struct {
	Mux   core.Mux
	Db    mysql.Repo
	Cache redis.Repo
}

func NewHTTPServer(logger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	r := new(resource)
	r.logger = logger

	// 初始化 DB
	dbRepo, err := mysql.New()
	if err != nil {
		logger.Fatal("new db err", zap.Error(err))
		panic(err)
	}

	r.db = dbRepo
	// 初始化 Cache
	cacheRepo, err := redis.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
		panic(err)
	}
	r.cache = cacheRepo

	// 初始化 开启相关功能组
	mux, err := core.New(logger,
		core.WithDisablePrometheus(), // 关闭 prometheus
		core.WithDisablePProf(),      // 关闭 WithDisablePProf

		core.WithEnableCors(),                                 //跨域
		core.WithEnableRate(),                                 //限流
		core.WithAlertNotify(alert.NotifyHandler(logger)),     //报警
		core.WithRecordMetrics(metrics.RecordHandler(logger)), //监控
	)

	if err != nil {
		panic(err)
	}

	r.mux = mux
	r.interceptors = interceptor.New(logger, r.cache, r.db)

	// 设置 API 路由
	setApiRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Db = r.db
	s.Cache = r.cache

	return s, nil
}

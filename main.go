package main

import (
	"context"
	"flag"
	"go-porter/configs"
	"go-porter/internal/svc"
	"go-porter/pkg/core/pkg/conf"
	"go-porter/pkg/core/pkg/net/httpx"
	"log"
	"net/http"
	"time"

	"go-porter/internal/interface/router"
	"go-porter/pkg/core/pkg/shutdown"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "/configs/config.toml", "the config file")

func main() {
	flag.Parse()

	conf.Init(configFile)
	c := conf.Get()

	// 初始化服务组件
	svcCtx := svc.NewServiceContext(c)

	//开启相关功能组
	mux, err := httpx.New(svcCtx.Logger,
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

	defer func() {
		// 确保关闭时日志已经刷到磁盘
		_ = svcCtx.Logger.Sync()
	}()

	// 初始化 router
	router.SetApiRouter(svcCtx, mux)

	// 配置 http server
	server := &http.Server{
		Addr:    configs.ProjectPort,
		Handler: mux,
	}

	// 启动 http server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			svcCtx.Logger.Fatal("http server startup err", zap.Error(err))
		}
	}()

	// 关闭服务
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				svcCtx.Logger.Error("server shutdown err", zap.Error(err))
			}
		},

		// 关闭 db
		func() {
			if svcCtx.Db != nil {
				if err := svcCtx.Db.DbWClose(); err != nil {
					svcCtx.Logger.Error("dbw close err", zap.Error(err))
				}

				if err := svcCtx.Db.DbRClose(); err != nil {
					svcCtx.Logger.Error("dbr close err", zap.Error(err))
				}
			}
		},

		// 关闭 redis
		func() {
			if svcCtx.Redis != nil {
				if err := svcCtx.Redis.Close(); err != nil {
					svcCtx.Logger.Error("cache close err", zap.Error(err))
				}
			}
		},
	)
}

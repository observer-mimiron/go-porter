package main

import (
	"context"
	"flag"
	"go-porter/configs"
	"go-porter/internal/svc"
	"go-porter/pkg/core/pkg/conf"
	"net/http"
	"time"

	"go-porter/internal/router"
	"go-porter/pkg/core/pkg/shutdown"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "/configs/config.toml", "the config file")

func main() {
	flag.Parse()

	//
	conf.Init(configFile)
	c := conf.Get()

	// 初始化服务组件
	serviceCtx := svc.NewServiceContext(c)
	defer func() {
		// 确保关闭时日志已经刷到磁盘
		_ = serviceCtx.Logger.Sync()
	}()

	// 初始化 router
	router.SetApiRouter(serviceCtx)

	// 配置 http server
	server := &http.Server{
		Addr:    configs.ProjectPort,
		Handler: serviceCtx.Mux,
	}

	// 启动 http server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serviceCtx.Logger.Fatal("http server startup err", zap.Error(err))
		}
	}()

	// 关闭服务
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				serviceCtx.Logger.Error("server shutdown err", zap.Error(err))
			}
		},

		// 关闭 db
		func() {
			if serviceCtx.Db != nil {
				if err := serviceCtx.Db.DbWClose(); err != nil {
					serviceCtx.Logger.Error("dbw close err", zap.Error(err))
				}

				if err := serviceCtx.Db.DbRClose(); err != nil {
					serviceCtx.Logger.Error("dbr close err", zap.Error(err))
				}
			}
		},

		// 关闭 redis
		func() {
			if serviceCtx.Redis != nil {
				if err := serviceCtx.Redis.Close(); err != nil {
					serviceCtx.Logger.Error("cache close err", zap.Error(err))
				}
			}
		},
	)
}

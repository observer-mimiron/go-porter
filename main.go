package main

import (
	"context"
	"flag"
	"go-porter/configs"
	"go-porter/pkg/core/pkg/conf"
	"net/http"
	"time"

	"go-porter/internal/router"
	"go-porter/pkg/core/pkg/logger"
	"go-porter/pkg/core/pkg/shutdown"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "/configs/config.toml", "the config file")

func main() {
	flag.Parse()

	conf.Init(configFile)
	c := conf.Get()
	// 初始化 access logger
	accessLogger, err := logger.NewJSONLogger(c.Log)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = accessLogger.Sync()
	}()

	// 初始化 HTTP 服务
	s, err := router.NewHTTPServer(accessLogger)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    configs.ProjectPort,
		Handler: s.Mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			accessLogger.Fatal("http server startup err", zap.Error(err))
		}
	}()

	// 优雅关闭
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				accessLogger.Error("server shutdown err", zap.Error(err))
			}
		},

		// 关闭 db
		func() {
			if s.Db != nil {
				if err := s.Db.DbWClose(); err != nil {
					accessLogger.Error("dbw close err", zap.Error(err))
				}

				if err := s.Db.DbRClose(); err != nil {
					accessLogger.Error("dbr close err", zap.Error(err))
				}
			}
		},

		// 关闭 cache
		func() {
			if s.Cache != nil {
				if err := s.Cache.Close(); err != nil {
					accessLogger.Error("cache close err", zap.Error(err))
				}
			}
		},
	)
}

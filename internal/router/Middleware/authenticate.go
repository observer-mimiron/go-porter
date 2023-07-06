package Middleware

import (
	"encoding/json"
	"go-porter/configs"
	"go-porter/internal/app/service/admin"
	"go-porter/pkg/core/pkg/database/mysql"
	"go.uber.org/zap"
	"net/http"

	"go-porter/internal/code"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/errors"
	"go-porter/pkg/core/pkg/net/httpx"
	"go-porter/pkg/core/pkg/proposal"
)

var _ Authenticate = (*authenticate)(nil)

type Authenticate interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx httpx.Context) (info proposal.SessionUserInfo, err httpx.BusinessError)
}

type authenticate struct {
	logger       *zap.Logger
	cache        redis.Repo
	db           mysql.Repo
	adminService admin.Service
}

func New(logger *zap.Logger, cache redis.Repo, db mysql.Repo) Authenticate {
	return &authenticate{
		logger:       logger,
		cache:        cache,
		db:           db,
		adminService: admin.New(db, cache),
	}
}

func (i *authenticate) i() {}

func (i *authenticate) CheckLogin(ctx httpx.Context) (sessionUserInfo proposal.SessionUserInfo, err httpx.BusinessError) {
	token := ctx.GetHeader(configs.HeaderLoginToken)
	if token == "" {
		err = httpx.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Token 参数"))

		return
	}

	if !i.cache.Exists(configs.RedisKeyPrefixLoginUser + token) {
		err = httpx.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("请先登录"))

		return
	}

	cacheData, cacheErr := i.cache.Get(configs.RedisKeyPrefixLoginUser+token, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		err = httpx.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(cacheErr)

		return
	}

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		httpx.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(jsonErr)

		return
	}

	return
}

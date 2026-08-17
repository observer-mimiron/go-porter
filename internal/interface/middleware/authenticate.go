package middleware

import (
	"encoding/json"
	"github.com/observer-mimiron/go-porter/configs"
	"github.com/observer-mimiron/go-porter/internal/errCode"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/cache/redis"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/database/mysql"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/net/httpx"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/proposal"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var _ Authenticate = (*authenticate)(nil)

type Authenticate interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx httpx.Context) (info proposal.SessionUserInfo, err *errCode.ErrCode)
}

type authenticate struct {
	logger *zap.Logger
	cache  redis.Repo
	db     mysql.Repo
}

func New(logger *zap.Logger, cache redis.Repo, db mysql.Repo) Authenticate {
	return &authenticate{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}

func (i *authenticate) i() {}

func (i *authenticate) CheckLogin(ctx httpx.Context) (sessionUserInfo proposal.SessionUserInfo, err *errCode.ErrCode) {
	token := ctx.GetHeader(configs.HeaderLoginToken)
	if token == "" {
		ctx.AbortWithError(errors.Wrap(errCode.ErrAuthorization, "Header 中缺少 Token 参数"))
		return sessionUserInfo, errCode.ErrAuthorization
	}

	if !i.cache.Exists(configs.RedisKeyPrefixLoginUser + token) {
		ctx.AbortWithError(errors.Wrap(errCode.ErrAuthorization, "请先登录"))
		return sessionUserInfo, errCode.ErrAuthorization
	}

	cacheData, cacheErr := i.cache.Get(configs.RedisKeyPrefixLoginUser+token, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		ctx.AbortWithError(errors.Wrapf(errCode.ErrAuthorization, "请先登录"))
		return sessionUserInfo, errCode.ErrAuthorization
	}

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		ctx.AbortWithError(errors.Wrapf(errCode.ErrAuthorization, "请先登录"))
		return sessionUserInfo, errCode.ErrAuthorization
	}

	return
}

package middleware

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go-porter/configs"
	"go-porter/internal/errCode"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/net/httpx"
	"go-porter/pkg/core/pkg/proposal"
	"go.uber.org/zap"
)

var _ Authenticate = (*authenticate)(nil)

type Authenticate interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx httpx.Context) (info proposal.SessionUserInfo, err error)
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

func (i *authenticate) CheckLogin(ctx httpx.Context) (sessionUserInfo proposal.SessionUserInfo, err error) {
	token := ctx.GetHeader(configs.HeaderLoginToken)
	if token == "" {
		//err = errors.Error(
		//	http.StatusUnauthorized,
		//	ecode.ErrAuthorizationError,
		//	ecode.Text(ecode.ErrAuthorizationError)).WithError(errors.New("Header 中缺少 Token 参数"))
		ctx.AbortWithError(errors.Wrap(errCode.ErrAuthorization, "Header 中缺少 Token 参数"))
		return
	}

	if !i.cache.Exists(configs.RedisKeyPrefixLoginUser + token) {
		//err = errors.Error(
		//	http.StatusUnauthorized,
		//	ecode.ErrAuthorizationError,
		//	ecode.Text(ecode.ErrAuthorizationError)).WithError(errors.New("请先登录"))
		ctx.AbortWithError(errors.Wrap(errCode.ErrAuthorization, "请先登录"))
		return
	}

	cacheData, cacheErr := i.cache.Get(configs.RedisKeyPrefixLoginUser+token, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		//err = errors.Error(
		//	http.StatusUnauthorized,
		//	ecode.ErrAuthorizationError,
		//	ecode.Text(ecode.ErrAuthorizationError)).WithError(cacheErr)
		ctx.AbortWithError(errors.Wrapf(errCode.ErrAuthorization, "请先登录"))
		return
	}

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		//errors.Error(
		//	http.StatusUnauthorized,
		//	ecode.ErrAuthorizationError,
		//	ecode.Text(ecode.ErrAuthorizationError)).WithError(jsonErr)
		ctx.AbortWithError(errors.Wrapf(errCode.ErrAuthorization, "请先登录"))
		return
	}

	return
}

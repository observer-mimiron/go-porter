package interceptor

import (
	"encoding/json"
	"go-porter/configs"
	"net/http"

	"go-porter/internal/code"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/errors"
	"go-porter/pkg/core/pkg/net/httpx"
	"go-porter/pkg/core/pkg/proposal"
)

func (i *interceptor) CheckLogin(ctx httpx.Context) (sessionUserInfo proposal.SessionUserInfo, err httpx.BusinessError) {
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

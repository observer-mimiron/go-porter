package admin

import (
	"github.com/pkg/errors"
	"go-porter/configs"
	"go-porter/internal/errCode"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
)

type logoutResponse struct {
	Username string `json:"username"` // 用户账号
}

// Logout 管理员登出
// @Summary 管理员登出
// @Description 管理员登出
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} logoutResponse
// @Failure 400 {object} ecode.Failure
// @Router /hanlder/admin/logout [post]
// @Security LoginToken
func (h *handler) Logout() httpx.HandlerFunc {
	return func(c httpx.Context) {
		res := new(logoutResponse)
		res.Username = c.SessionUserInfo().UserName

		if !h.svcCtx.Redis.Del(configs.RedisKeyPrefixLoginUser+c.GetHeader(configs.HeaderLoginToken), redis.WithTrace(c.Trace())) {
			c.AbortWithError(errors.Wrap(errCode.ErrAdminLogOut, "Logout error"))
			return
		}

		c.Payload(res)
	}
}

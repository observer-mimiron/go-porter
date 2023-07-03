package admin

import (
	"go-porter/configs"
	"go-porter/pkg/core/pkg/net/httpx"
	"net/http"

	"go-porter/internal/app/service/admin"
	"go-porter/internal/code"
	"go-porter/internal/pkg/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/errors"
	"go-porter/pkg/core/pkg/proposal"
)

type loginRequest struct {
	Username string `form:"username"` // 用户名
	Password string `form:"password"` // 密码
}

type loginResponse struct {
	Token string `json:"token"` // 用户身份标识
}

// Login 管理员登录
// @Summary 管理员登录
// @Description 管理员登录
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "MD5后的密码"
// @Success 200 {object} loginResponse
// @Failure 400 {object} code.Failure
// @Router /api/login [post]
// @Security LoginToken
func (h *handler) Login() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(loginRequest)
		res := new(loginResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		searchOneData := new(admin.SearchOneData)
		searchOneData.Username = req.Username
		searchOneData.Password = password.GeneratePassword(req.Password)
		searchOneData.IsUsed = 1
		info, err := h.adminService.Detail(c, searchOneData)
		if err != nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		if info == nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(errors.New("未查询出符合条件的用户")),
			)
			return
		}

		token := password.GenerateLoginToken(info.Id)

		// 用户信息
		sessionUserInfo := &proposal.SessionUserInfo{
			UserID:   info.Id,
			UserName: info.Username,
		}

		// 将用户信息记录到 Redis 中
		err = h.cache.Set(configs.RedisKeyPrefixLoginUser+token, string(sessionUserInfo.Marshal()), configs.LoginSessionTTL, redis.WithTrace(c.Trace()))
		if err != nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		res.Token = token
		c.Payload(res)
	}
}

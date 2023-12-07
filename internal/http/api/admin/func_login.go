package admin

import (
	"github.com/pkg/errors"
	"go-porter/configs"
	"go-porter/internal/ecode"
	"go-porter/internal/service/admin"
	"go-porter/internal/util/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
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
// @Failure 400 {object} ecode.Failure
// @Router /api/login [post]
// @Security LoginToken
func (h *handler) Login() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(loginRequest)
		res := new(loginResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrParamBind, "Create error %+v", err))
			return
		}

		searchOneData := new(admin.SearchOneData)
		searchOneData.Username = req.Username
		searchOneData.Password = password.GeneratePassword(req.Password)
		searchOneData.IsUsed = 1

		adminService := admin.New(h.svcCtx)
		info, err := adminService.Detail(c, searchOneData)
		if err != nil {
			c.AbortWithError(errors.WithMessage(err, "Login error"))
			return
		}

		if info == nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrUserNotFound, "Login error %+v", searchOneData))
			return
		}

		token := password.GenerateLoginToken(info.Id)

		// 用户信息
		sessionUserInfo := &proposal.SessionUserInfo{
			UserID:   info.Id,
			UserName: info.Username,
		}

		// 将用户信息记录到 Redis 中
		err = h.svcCtx.Redis.Set(configs.RedisKeyPrefixLoginUser+token, string(sessionUserInfo.Marshal()), configs.LoginSessionTTL, redis.WithTrace(c.Trace()))
		if err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrAdminLogin, "Login error %+v", err))
			return
		}

		res.Token = token
		c.Payload(res)
	}
}

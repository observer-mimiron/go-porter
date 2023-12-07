package admin

import (
	"github.com/pkg/errors"
	"go-porter/configs"
	"go-porter/internal/ecode"
	"go-porter/internal/util/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
)

type offlineRequest struct {
	Id string `form:"id"` // 主键ID
}

type offlineResponse struct {
	Id int32 `json:"id"` // 主键ID
}

// Offline 下线管理员
// @Summary 下线管理员
// @Description 下线管理员
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id formData string true "Hashid"
// @Success 200 {object} offlineResponse
// @Failure 400 {object} ecode.Failure
// @Router /api/admin/offline [patch]
// @Security LoginToken
func (h *handler) Offline() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(offlineRequest)
		res := new(offlineResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrParamBind, "Offline error %+v", err))
			return
		}

		ids, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrHashIdsDxerror, "Offline error %+v", err))
			return
		}

		id := int32(ids[0])

		b := h.svcCtx.Redis.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(c.Trace()))
		if !b {
			c.AbortWithError(errors.Wrapf(ecode.ErrAdminOffline, "Offline error %+v", err))
			return
		}

		res.Id = id
		c.Payload(res)
	}
}

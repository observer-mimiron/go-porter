package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/errCode"
	"go-porter/internal/service/admin"
	"go-porter/pkg/core/pkg/net/httpx"
)

type resetPasswordRequest struct {
	Id string `uri:"id"` // HashID
}

type resetPasswordResponse struct {
	Id int32 `json:"id"` // 主键ID
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 重置密码
// @Tags API.admin
// @Accept json
// @Produce json
// @Param id path string true "hashId"
// @Success 200 {object} resetPasswordResponse
// @Failure 400 {object} ecode.Failure
// @Router /hanlder/admin/reset_password/{id} [patch]
// @Security LoginToken
func (h *handler) ResetPassword() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(resetPasswordRequest)
		res := new(resetPasswordResponse)
		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrParamBind, "ResetPassword error %+v", err))
			return
		}

		ids, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrServer, "ResetPassword  error %+v", err))
			return
		}

		id := int32(ids[0])

		adminService := admin.New(h.svcCtx)
		err = adminService.ResetPassword(c, id)
		if err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrAdminResetPassword, "ResetPassword  error %+v", err))
			return
		}

		res.Id = id
		c.Payload(res)
	}
}

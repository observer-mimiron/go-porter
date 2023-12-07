package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/ecode"
	"go-porter/internal/service/admin"
	"go-porter/internal/util/password"
	"go-porter/pkg/core/pkg/net/httpx"
)

type modifyPasswordRequest struct {
	OldPassword string `form:"old_password"` // 旧密码
	NewPassword string `form:"new_password"` // 新密码
}

type modifyPasswordResponse struct {
	Username string `json:"username"` // 用户账号
}

// ModifyPassword 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param old_password formData string true "旧密码"
// @Param new_password formData string true "新密码"
// @Success 200 {object} modifyPasswordResponse
// @Failure 400 {object} ecode.Failure
// @Router /api/admin/modify_password [patch]
// @Security LoginToken
func (h *handler) ModifyPassword() httpx.HandlerFunc {
	return func(ctx httpx.Context) {
		req := new(modifyPasswordRequest)
		res := new(modifyPasswordResponse)
		if err := ctx.ShouldBindForm(req); err != nil {
			ctx.AbortWithError(errors.Wrapf(ecode.ErrParamBind, "ModifyPassword error %+v", err))
			return
		}

		searchOneData := new(admin.SearchOneData)
		searchOneData.Id = ctx.SessionUserInfo().UserID
		searchOneData.Password = password.GeneratePassword(req.OldPassword)
		searchOneData.IsUsed = 1

		adminService := admin.New(h.svcCtx)
		info, err := adminService.Detail(ctx, searchOneData)
		if err != nil || info == nil {
			ctx.AbortWithError(err)
			return
		}

		if err := adminService.ModifyPassword(ctx, ctx.SessionUserInfo().UserID, req.NewPassword); err != nil {
			ctx.AbortWithError(errors.Wrapf(ecode.ErrAdminModifyPassword, "ModifyPassword error %+v", err))
			return
		}

		res.Username = ctx.SessionUserInfo().UserName
		ctx.Payload(res)
	}
}

package admin

import (
	"go-porter/pkg/core/pkg/net/httpx"
	"net/http"

	"go-porter/internal/http/code"
	"go-porter/internal/service/admin"
)

type modifyPersonalInfoRequest struct {
	Nickname string `form:"nickname"` // 昵称
	Mobile   string `form:"mobile"`   // 手机号
}

type modifyPersonalInfoResponse struct {
	Username string `json:"username"` // 用户账号
}

// ModifyPersonalInfo 修改个人信息
// @Summary 修改个人信息
// @Description 修改个人信息
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param nickname formData string true "昵称"
// @Param mobile formData string true "手机号"
// @Success 200 {object} modifyPersonalInfoResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/modify_personal_info [patch]
// @Security LoginToken
func (h *handler) ModifyPersonalInfo() httpx.HandlerFunc {
	return func(ctx httpx.Context) {
		req := new(modifyPersonalInfoRequest)
		res := new(modifyPersonalInfoResponse)
		if err := ctx.ShouldBindForm(req); err != nil {
			ctx.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		modifyData := new(admin.ModifyData)
		modifyData.Nickname = req.Nickname
		modifyData.Mobile = req.Mobile

		adminService := admin.New(h.svcCtx)
		if err := adminService.ModifyPersonalInfo(ctx, ctx.SessionUserInfo().UserID, modifyData); err != nil {
			ctx.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.AdminModifyPersonalInfoError,
				code.Text(code.AdminModifyPersonalInfoError)).WithError(err),
			)
			return
		}

		res.Username = ctx.SessionUserInfo().UserName
		ctx.Payload(res)
	}
}
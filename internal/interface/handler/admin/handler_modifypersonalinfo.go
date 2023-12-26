package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/errCode"
	"go-porter/internal/service/admin"
	"go-porter/pkg/core/pkg/net/httpx"
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
// @Failure 400 {object} ecode.Failure
// @Router /hanlder/admin/modify_personal_info [patch]
// @Security LoginToken
func (h *handler) ModifyPersonalInfo() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(modifyPersonalInfoRequest)
		res := new(modifyPersonalInfoResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrParamBind, "ModifyPersonalInfo error %+v", err))
			return
		}

		modifyData := new(admin.ModifyData)
		modifyData.Nickname = req.Nickname
		modifyData.Mobile = req.Mobile

		adminService := admin.New(h.svcCtx)
		if err := adminService.ModifyPersonalInfo(c, c.SessionUserInfo().UserID, modifyData); err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrParamBind, "ModifyPersonalInfo error %+v", err))
			return
		}

		res.Username = c.SessionUserInfo().UserName
		c.Payload(res)
	}
}

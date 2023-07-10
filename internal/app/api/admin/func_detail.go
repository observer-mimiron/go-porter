package admin

import (
	"go-porter/pkg/core/pkg/net/httpx"
	"net/http"

	"go-porter/internal/app/service/admin"
	"go-porter/internal/code"
)

type detailResponse struct {
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Mobile   string `json:"mobile"`   // 手机号
}

// Detail 管理员详情
// @Summary 管理员详情
// @Description 管理员详情
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/info [get]
// @Security LoginToken
func (h *handler) Detail() httpx.HandlerFunc {
	return func(ctx httpx.Context) {
		res := new(detailResponse)

		searchOneData := new(admin.SearchOneData)
		searchOneData.Id = ctx.SessionUserInfo().UserID
		searchOneData.IsUsed = 1

		adminService := admin.New(h.svcCtx)
		info, err := adminService.Detail(ctx, searchOneData)
		if err != nil {
			ctx.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		res.Username = info.Username
		res.Nickname = info.Nickname
		res.Mobile = info.Mobile
		ctx.Payload(res)
	}
}

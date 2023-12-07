package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/ecode"
	"go-porter/internal/service/admin"
	"go-porter/pkg/core/pkg/net/httpx"
)

type createRequest struct {
	Username string `form:"username" binding:"required"` // 用户名
	Nickname string `form:"nickname" binding:"required"` // 昵称
	Mobile   string `form:"mobile" binding:"required"`   // 手机号
	Password string `form:"password" binding:"required"` // MD5后的密码
}

type createResponse struct {
	Id int32 `json:"id"` // 主键ID
}

// Create 新增管理员
// @Summary 新增管理员
// @Description 新增管理员
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "用户名"
// @Param nickname formData string true "昵称"
// @Param mobile formData string true "手机号"
// @Param password formData string true "MD5后的密码"
// @Success 200 {object} createResponse
// @Failure 400 {object} ecode.Failure
// @Router /api/admin [post]
// @Security LoginToken
func (h *handler) Create() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(createRequest)
		res := new(createResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrParamBind, "Create error %+v", err))
			return
		}

		createData := new(admin.CreateAdminData)
		createData.Nickname = req.Nickname
		createData.Username = req.Username
		createData.Mobile = req.Mobile
		createData.Password = req.Password
		adminService := admin.New(h.svcCtx)
		id, err := adminService.Create(c, createData)
		if err != nil {
			c.AbortWithError(errors.WithMessage(err, "Create error"))
			return
		}

		res.Id = id
		c.Payload(res)
	}
}

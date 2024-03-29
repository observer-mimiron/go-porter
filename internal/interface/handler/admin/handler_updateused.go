package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/errCode"
	"go-porter/internal/service/admin"
	"go-porter/pkg/core/pkg/net/httpx"
)

type updateUsedRequest struct {
	Id   string `form:"id"`   // 主键ID
	Used int32  `form:"used"` // 是否启用 1:是 -1:否
}

type updateUsedResponse struct {
	Id int32 `json:"id"` // 主键ID
}

// UpdateUsed 更新管理员为启用/禁用
// @Summary 更新管理员为启用/禁用
// @Description 更新管理员为启用/禁用
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id formData string true "Hashid"
// @Param used formData int true "是否启用 1:是 -1:否"
// @Success 200 {object} updateUsedResponse
// @Failure 400 {object} ecode.Failure
// @Router /hanlder/admin/used [patch]
// @Security LoginToken
func (h *handler) UpdateUsed() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(updateUsedRequest)
		res := new(updateUsedResponse)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrParamBind, "UpdateUsed error %+v", err))
			return
		}

		ids, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrServer, "UpdateUsed error %+v", err))
			return
		}

		id := int32(ids[0])
		adminService := admin.New(h.svcCtx)
		err = adminService.UpdateUsed(c, id, req.Used)
		if err != nil {
			c.AbortWithError(errors.Wrapf(errCode.ErrServer, "UpdateUsed error %+v", err))
			return
		}

		res.Id = id
		c.Payload(res)
	}
}

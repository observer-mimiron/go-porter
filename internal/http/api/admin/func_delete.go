package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/ecode"
	"go-porter/internal/service/admin"
	"go-porter/pkg/core/pkg/net/httpx"
)

type deleteRequest struct {
	Id string `uri:"id"` // HashID
}

type deleteResponse struct {
	Id int32 `json:"id"` // 主键ID
}

// Delete 删除管理员
// @Summary 删除管理员
// @Description 删除管理员
// @Tags API.admin
// @Accept json
// @Produce json
// @Param id path string true "hashId"
// @Success 200 {object} deleteResponse
// @Failure 400 {object} ecode.Failure
// @Router /api/admin/{id} [delete]
// @Security LoginToken
func (h *handler) Delete() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(deleteRequest)
		res := new(deleteResponse)
		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrParamBind, "Delete error %+v", err))
			return
		}

		ids, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrHashIdsDxerror, "Delete  error %+v", err))
			return
		}

		id := int32(ids[0])
		adminService := admin.New(h.svcCtx)
		err = adminService.Delete(c, id)
		if err != nil {
			c.AbortWithError(errors.Wrapf(ecode.ErrHashIdsDxerror, "Delete  error %+v", err))
			return
		}

		res.Id = id
		c.Payload(res)
	}
}

package admin

import (
	"go-porter/internal/service/admin"
	"go-porter/pkg/core/pkg/net/httpx"
	"net/http"

	"go-porter/internal/http/code"
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
// @Failure 400 {object} code.Failure
// @Router /api/admin/{id} [delete]
// @Security LoginToken
func (h *handler) Delete() httpx.HandlerFunc {
	return func(c httpx.Context) {
		req := new(deleteRequest)
		res := new(deleteResponse)
		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		ids, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.HashIdsDecodeError,
				code.Text(code.HashIdsDecodeError)).WithError(err),
			)
			return
		}

		id := int32(ids[0])
		adminService := admin.New(h.svcCtx)
		err = adminService.Delete(c, id)
		if err != nil {
			c.AbortWithError(httpx.Error(
				http.StatusBadRequest,
				code.AdminDeleteError,
				code.Text(code.AdminDeleteError)).WithError(err),
			)
			return
		}

		res.Id = id
		c.Payload(res)
	}
}
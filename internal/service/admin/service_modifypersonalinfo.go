package admin

import (
	"go-porter/internal/model"
	"go-porter/pkg/core/pkg/net/httpx"
)

type ModifyData struct {
	Nickname string // 昵称
	Mobile   string // 手机号
}

func (s *service) ModifyPersonalInfo(ctx httpx.Context, id int32, modifyData *ModifyData) (err error) {
	data := map[string]interface{}{
		"nickname":     modifyData.Nickname,
		"mobile":       modifyData.Mobile,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	err = s.svc.Db.GetDbW().WithContext(ctx.RequestContext()).Model(&model.Admin{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}

	return
}

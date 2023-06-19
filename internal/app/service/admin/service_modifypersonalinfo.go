package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/core"
)

type ModifyData struct {
	Nickname string // 昵称
	Mobile   string // 手机号
}

func (s *service) ModifyPersonalInfo(ctx core.Context, id int32, modifyData *ModifyData) (err error) {
	data := map[string]interface{}{
		"nickname":     modifyData.Nickname,
		"mobile":       modifyData.Mobile,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := model.NewQueryBuilder()
	qb.WhereId("=", id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	return
}

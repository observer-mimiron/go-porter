package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/core"
)

type SearchOneData struct {
	Id       int32  // 用户ID
	Username string // 用户名
	Nickname string // 昵称
	Mobile   string // 手机号
	Password string // 密码
	IsUsed   int32  // 是否启用 1:是  -1:否
}

func (s *service) Detail(ctx core.Context, searchOneData *SearchOneData) (info *model.Admin, err error) {

	qb := model.NewQueryBuilder()
	qb.WhereIsDeleted("=", -1)

	if searchOneData.Id != 0 {
		qb.WhereId("=", searchOneData.Id)
	}

	if searchOneData.Username != "" {
		qb.WhereUsername("=", searchOneData.Username)
	}

	if searchOneData.Nickname != "" {
		qb.WhereNickname("=", searchOneData.Nickname)
	}

	if searchOneData.Mobile != "" {
		qb.WhereMobile("=", searchOneData.Mobile)
	}

	if searchOneData.Password != "" {
		qb.WherePassword("=", searchOneData.Password)
	}

	if searchOneData.IsUsed != 0 {
		qb.WhereIsUsed("=", searchOneData.IsUsed)
	}

	info, err = qb.QueryOne(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}

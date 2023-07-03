package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/net/httpx"
)

type SearchOneData struct {
	Id       int32  // 用户ID
	Username string // 用户名
	Nickname string // 昵称
	Mobile   string // 手机号
	Password string // 密码
	IsUsed   int32  // 是否启用 1:是  -1:否
}

func (s *service) Detail(ctx httpx.Context, searchOneData *SearchOneData) (info *model.Admin, err error) {
	qb := s.db.GetDbR().WithContext(ctx.RequestContext()).Model(&model.Admin{})
	qb.Where("is_deleted = ?", -1)

	if searchOneData.Id != 0 {
		qb.Where("id = ?", searchOneData.Id)
	}

	if searchOneData.Username != "" {
		qb.Where("username = ?", searchOneData.Username)
	}

	if searchOneData.Nickname != "" {
		qb.Where("nickname = ?", searchOneData.Nickname)
	}

	if searchOneData.Mobile != "" {
		qb.Where("mobile = ?", searchOneData.Mobile)
	}

	if searchOneData.Password != "" {
		qb.Where("password = ?", searchOneData.Password)
	}

	if searchOneData.IsUsed != 0 {
		qb.Where("is_used = ?", searchOneData.IsUsed)
	}

	err = qb.First(&info).Error
	if err != nil {
		return nil, err
	}

	return
}

package admin

import (
	"go-porter/internal/app/model"
	"go-porter/internal/pkg/core"
	"go-porter/internal/pkg/mysql"
)

type SearchData struct {
	Page     int    // 第几页
	PageSize int    // 每页显示条数
	Username string // 用户名
	Nickname string // 昵称
	Mobile   string // 手机号
}

func (s *service) PageList(ctx core.Context, searchData *SearchData) (listData []*model.Admin, err error) {

	page := searchData.Page
	if page == 0 {
		page = 1
	}

	pageSize := searchData.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	qb := model.NewQueryBuilder()
	qb.WhereIsDeleted(mysql.EqualPredicate, -1)

	if searchData.Username != "" {
		qb.WhereUsername(mysql.EqualPredicate, searchData.Username)
	}

	if searchData.Nickname != "" {
		qb.WhereNickname(mysql.EqualPredicate, searchData.Nickname)
	}

	if searchData.Mobile != "" {
		qb.WhereMobile(mysql.EqualPredicate, searchData.Mobile)
	}

	listData, err = qb.
		Limit(pageSize).
		Offset(offset).
		OrderById(false).
		QueryAll(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}

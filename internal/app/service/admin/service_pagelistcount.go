package admin

import (
	"go-porter/internal/app/model"
	"go-porter/internal/pkg/core"
	"go-porter/internal/pkg/mysql"
)

func (s *service) PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error) {

	qb := model.NewQueryBuilder()
	qb = qb.WhereIsDeleted(mysql.EqualPredicate, -1)

	if searchData.Username != "" {
		qb.WhereUsername(mysql.EqualPredicate, searchData.Username)
	}

	if searchData.Nickname != "" {
		qb.WhereNickname(mysql.EqualPredicate, searchData.Nickname)
	}

	if searchData.Mobile != "" {
		qb.WhereMobile(mysql.EqualPredicate, searchData.Mobile)
	}

	total, err = qb.Count(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return 0, err
	}

	return
}

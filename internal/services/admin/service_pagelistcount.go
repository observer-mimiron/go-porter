package admin

import (
	"go-porter/internal/pkg/core"
	"go-porter/internal/repository/mysql"
	"go-porter/internal/repository/mysql/admin"
)

func (s *service) PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error) {

	qb := admin.NewQueryBuilder()
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

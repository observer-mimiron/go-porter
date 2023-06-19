package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/core"
)

func (s *service) PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error) {

	qb := model.NewQueryBuilder()
	qb = qb.WhereIsDeleted("=", -1)

	if searchData.Username != "" {
		qb.WhereUsername("=", searchData.Username)
	}

	if searchData.Nickname != "" {
		qb.WhereNickname("=", searchData.Nickname)
	}

	if searchData.Mobile != "" {
		qb.WhereMobile("=", searchData.Mobile)
	}

	total, err = qb.Count(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return 0, err
	}

	return
}

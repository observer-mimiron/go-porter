package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/net/httpx"
)

type SearchData struct {
	Page     int    // 第几页
	PageSize int    // 每页显示条数
	Username string // 用户名
	Nickname string // 昵称
	Mobile   string // 手机号
}

func (s *service) PageList(ctx httpx.Context, searchData *SearchData) (listData []*model.Admin, count int64, err error) {

	page := searchData.Page
	if page == 0 {
		page = 1
	}

	pageSize := searchData.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	qb := s.db.GetDbR().WithContext(ctx.RequestContext()).Model(&model.Admin{})
	if searchData.Username != "" {
		qb.Where("username like ?", "%"+searchData.Username+"%")
	}

	if searchData.Nickname != "" {
		qb.Where("nickname like ?", "%"+searchData.Nickname+"%")
	}

	if searchData.Mobile != "" {
		qb.Where("mobile = ?", searchData.Mobile)
	}

	qb.Count(&count)
	qb.Limit(pageSize).Offset(offset)

	err = qb.Find(&listData).Error
	if err != nil {
		return nil, 0, err
	}

	return listData, count, nil
}

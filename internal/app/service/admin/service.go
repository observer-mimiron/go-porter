package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/core"
	"go-porter/pkg/core/pkg/database/mysql"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx core.Context, adminData *CreateAdminData) (id int32, err error)
	PageList(ctx core.Context, searchData *SearchData) (listData []*model.Admin, err error)
	PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error)
	UpdateUsed(ctx core.Context, id int32, used int32) (err error)
	Delete(ctx core.Context, id int32) (err error)
	Detail(ctx core.Context, searchOneData *SearchOneData) (info *model.Admin, err error)
	ResetPassword(ctx core.Context, id int32) (err error)
	ModifyPassword(ctx core.Context, id int32, newPassword string) (err error)
	ModifyPersonalInfo(ctx core.Context, id int32, modifyData *ModifyData) (err error)
}

type service struct {
	db    mysql.Repo
	cache redis.Repo
}

func New(db mysql.Repo, cache redis.Repo) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) i() {}

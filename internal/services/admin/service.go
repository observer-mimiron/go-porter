package admin

import (
	"go-porter/internal/pkg/core"
	"go-porter/internal/repository/mysql"
	"go-porter/internal/repository/mysql/admin"
	"go-porter/internal/repository/redis"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	Create(ctx core.Context, adminData *CreateAdminData) (id int32, err error)
	PageList(ctx core.Context, searchData *SearchData) (listData []*admin.Admin, err error)
	PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error)
	UpdateUsed(ctx core.Context, id int32, used int32) (err error)
	Delete(ctx core.Context, id int32) (err error)
	Detail(ctx core.Context, searchOneData *SearchOneData) (info *admin.Admin, err error)
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

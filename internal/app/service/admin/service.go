package admin

import (
	"go-porter/internal/app/model"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/net/httpx"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(ctx httpx.Context, adminData *CreateAdminData) (id int32, err error)
	PageList(ctx httpx.Context, searchData *SearchData) (listData []*model.Admin, count int64, err error)
	UpdateUsed(ctx httpx.Context, id int32, used int32) (err error)
	Delete(ctx httpx.Context, id int32) (err error)
	Detail(ctx httpx.Context, searchOneData *SearchOneData) (info *model.Admin, err error)
	ResetPassword(ctx httpx.Context, id int32) (err error)
	ModifyPassword(ctx httpx.Context, id int32, newPassword string) (err error)
	ModifyPersonalInfo(ctx httpx.Context, id int32, modifyData *ModifyData) (err error)
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

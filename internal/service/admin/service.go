package admin

import (
	"go-porter/internal/model"
	"go-porter/internal/svc"
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
	svc *svc.ServiceContext
}

func New(svc *svc.ServiceContext) Service {
	return &service{
		svc: svc,
	}
}

func (s *service) i() {}

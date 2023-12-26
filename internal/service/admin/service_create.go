package admin

import (
	"github.com/pkg/errors"
	"go-porter/internal/errCode"
	"go-porter/internal/model"
	"go-porter/internal/util/password"
	"go-porter/pkg/core/pkg/net/httpx"
)

type CreateAdminData struct {
	Username string // 用户名
	Nickname string // 昵称
	Mobile   string // 手机号
	Password string // 密码
}

func (s *service) Create(ctx httpx.Context, adminData *CreateAdminData) (id int32, err error) {
	admin := new(model.Admin)
	admin.Username = adminData.Username
	admin.Password = password.GeneratePassword(adminData.Password)
	admin.Nickname = adminData.Nickname
	admin.Mobile = adminData.Mobile
	admin.CreatedUser = ctx.SessionUserInfo().UserName
	admin.IsUsed = 1
	admin.IsDeleted = -1

	err = s.svc.Db.GetDbW().Create(admin).Error
	if err != nil {
		return 0, errors.Wrapf(errCode.ErrServer, "创建管理员失败: %v", err)
	}
	return admin.Id, nil
}

package admin

import (
	"github.com/observer-mimiron/go-porter/configs"
	"github.com/observer-mimiron/go-porter/internal/model"
	"github.com/observer-mimiron/go-porter/internal/util/password"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/cache/redis"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/net/httpx"
)

func (s *service) ModifyPassword(ctx httpx.Context, id int32, newPassword string) (err error) {
	data := map[string]interface{}{
		"password":     password.GeneratePassword(newPassword),
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	err = s.svc.Db.GetDbW().WithContext(ctx.RequestContext()).Model(&model.Admin{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}

	s.svc.Redis.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

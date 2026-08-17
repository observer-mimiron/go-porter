package admin

import (
	"github.com/observer-mimiron/go-porter/configs"
	"github.com/observer-mimiron/go-porter/internal/model"
	"github.com/observer-mimiron/go-porter/internal/util/password"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/cache/redis"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/net/httpx"
)

func (s *service) ResetPassword(ctx httpx.Context, id int32) (err error) {
	data := map[string]interface{}{
		"password":     password.ResetPassword(),
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	s.svc.Db.GetDbW().Model(&model.Admin{}).Where("id = ?", id).Updates(data)
	s.svc.Redis.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

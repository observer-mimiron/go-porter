package admin

import (
	"go-porter/configs"
	"go-porter/internal/model"
	"go-porter/internal/util/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
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

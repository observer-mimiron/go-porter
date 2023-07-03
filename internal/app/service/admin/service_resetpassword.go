package admin

import (
	"go-porter/configs"
	"go-porter/internal/app/model"
	"go-porter/internal/pkg/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
)

func (s *service) ResetPassword(ctx httpx.Context, id int32) (err error) {
	data := map[string]interface{}{
		"password":     password.ResetPassword(),
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	s.db.GetDbW().Model(&model.Admin{}).Where("id = ?", id).Updates(data)
	s.cache.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

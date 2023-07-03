package admin

import (
	"go-porter/configs"
	"go-porter/internal/app/model"
	"go-porter/internal/pkg/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
)

func (s *service) ModifyPassword(ctx httpx.Context, id int32, newPassword string) (err error) {
	data := map[string]interface{}{
		"password":     password.GeneratePassword(newPassword),
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	err = s.db.GetDbW().WithContext(ctx.RequestContext()).Model(&model.Admin{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}

	s.cache.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

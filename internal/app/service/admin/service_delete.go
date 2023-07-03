package admin

import (
	"go-porter/configs"
	"go-porter/internal/app/model"
	"go-porter/internal/pkg/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/net/httpx"
)

func (s *service) Delete(ctx httpx.Context, id int32) (err error) {
	data := map[string]interface{}{
		"is_deleted":   1,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	err = s.db.GetDbW().WithContext(ctx.RequestContext()).Model(&model.Admin{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}
	s.cache.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

package admin

import (
	"go-porter/configs"
	"go-porter/internal/app/model"
	"go-porter/internal/pkg/password"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/core"
)

func (s *service) ModifyPassword(ctx core.Context, id int32, newPassword string) (err error) {
	data := map[string]interface{}{
		"password":     password.GeneratePassword(newPassword),
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := model.NewQueryBuilder()
	qb.WhereId("=", id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	s.cache.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

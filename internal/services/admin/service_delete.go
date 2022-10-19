package admin

import (
	"go-porter/configs"
	"go-porter/internal/pkg/core"
	"go-porter/internal/pkg/password"
	"go-porter/internal/repository/mysql"
	"go-porter/internal/repository/mysql/admin"
	"go-porter/internal/repository/redis"
)

func (s *service) Delete(ctx core.Context, id int32) (err error) {
	data := map[string]interface{}{
		"is_deleted":   1,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := admin.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	s.cache.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}

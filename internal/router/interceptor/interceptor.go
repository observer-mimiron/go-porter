package interceptor

import (
	"go-porter/internal/app/service/admin"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/core"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/proposal"

	"go.uber.org/zap"
)

var _ Interceptor = (*interceptor)(nil)

type Interceptor interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx core.Context) (info proposal.SessionUserInfo, err core.BusinessError)
}

type interceptor struct {
	logger       *zap.Logger
	cache        redis.Repo
	db           mysql.Repo
	adminService admin.Service
}

func New(logger *zap.Logger, cache redis.Repo, db mysql.Repo) Interceptor {
	return &interceptor{
		logger:       logger,
		cache:        cache,
		db:           db,
		adminService: admin.New(db, cache),
	}
}

func (i *interceptor) i() {}

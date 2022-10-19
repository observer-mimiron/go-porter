package interceptor

import (
	"go-porter/internal/pkg/core"
	"go-porter/internal/proposal"
	"go-porter/internal/repository/mysql"
	"go-porter/internal/repository/redis"
	"go-porter/internal/services/admin"

	"go.uber.org/zap"
)

var _ Interceptor = (*interceptor)(nil)

type Interceptor interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx core.Context) (info proposal.SessionUserInfo, err core.BusinessError)
	// i 为了避免被其他包实现
	i()
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

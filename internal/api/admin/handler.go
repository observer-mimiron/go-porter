package admin

import (
	"go-porter/configs"
	"go-porter/internal/pkg/core"
	"go-porter/internal/repository/mysql"
	"go-porter/internal/repository/redis"
	"go-porter/internal/services/admin"
	"go-porter/pkg/hash"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Login 管理员登录
	// @Tags API.admin
	// @Router /api/login [post]
	Login() core.HandlerFunc

	// Logout 管理员登出
	// @Tags API.admin
	// @Router /api/admin/logout [post]
	Logout() core.HandlerFunc

	// ModifyPassword 修改密码
	// @Tags API.admin
	// @Router /api/admin/modify_password [patch]
	ModifyPassword() core.HandlerFunc

	// Detail 个人信息
	// @Tags API.admin
	// @Router /api/admin/info [get]
	Detail() core.HandlerFunc

	// ModifyPersonalInfo 修改个人信息
	// @Tags API.admin
	// @Router /api/admin/modify_personal_info [patch]
	ModifyPersonalInfo() core.HandlerFunc

	// Create 新增管理员
	// @Tags API.admin
	// @Router /api/admin [post]
	Create() core.HandlerFunc

	// List 管理员列表
	// @Tags API.admin
	// @Router /api/admin [get]
	List() core.HandlerFunc

	// Delete 删除管理员
	// @Tags API.admin
	// @Router /api/admin/{id} [delete]
	Delete() core.HandlerFunc

	// Offline 下线管理员
	// @Tags API.admin
	// @Router /api/admin/offline [patch]
	Offline() core.HandlerFunc

	// UpdateUsed 更新管理员为启用/禁用
	// @Tags API.admin
	// @Router /api/admin/used [patch]
	UpdateUsed() core.HandlerFunc

	// ResetPassword 重置密码
	// @Tags API.admin
	// @Router /api/admin/reset_password/{id} [patch]
	ResetPassword() core.HandlerFunc
}

type handler struct {
	logger       *zap.Logger
	cache        redis.Repo
	hashids      hash.Hash
	adminService admin.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:       logger,
		cache:        cache,
		hashids:      hash.New(configs.Get().HashIds.Secret, configs.Get().HashIds.Length),
		adminService: admin.New(db, cache),
	}
}

func (h *handler) i() {}

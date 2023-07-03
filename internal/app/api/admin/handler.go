package admin

import (
	"go-porter/internal/app/service/admin"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/conf"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/net/httpx"
	"go-porter/pkg/cryptor/hash"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	// Login 管理员登录
	// @Tags API.admin
	// @Router /api/login [post]
	Login() httpx.HandlerFunc

	// Logout 管理员登出
	// @Tags API.admin
	// @Router /api/admin/logout [post]
	Logout() httpx.HandlerFunc

	// ModifyPassword 修改密码
	// @Tags API.admin
	// @Router /api/admin/modify_password [patch]
	ModifyPassword() httpx.HandlerFunc

	// Detail 个人信息
	// @Tags API.admin
	// @Router /api/admin/info [get]
	Detail() httpx.HandlerFunc

	// ModifyPersonalInfo 修改个人信息
	// @Tags API.admin
	// @Router /api/admin/modify_personal_info [patch]
	ModifyPersonalInfo() httpx.HandlerFunc

	// Create 新增管理员
	// @Tags API.admin
	// @Router /api/admin [post]
	Create() httpx.HandlerFunc

	// List 管理员列表
	// @Tags API.admin
	// @Router /api/admin [get]
	List() httpx.HandlerFunc

	// Delete 删除管理员
	// @Tags API.admin
	// @Router /api/admin/{id} [delete]
	Delete() httpx.HandlerFunc

	// Offline 下线管理员
	// @Tags API.admin
	// @Router /api/admin/offline [patch]
	Offline() httpx.HandlerFunc

	// UpdateUsed 更新管理员为启用/禁用
	// @Tags API.admin
	// @Router /api/admin/used [patch]
	UpdateUsed() httpx.HandlerFunc

	// ResetPassword 重置密码
	// @Tags API.admin
	// @Router /api/admin/reset_password/{id} [patch]
	ResetPassword() httpx.HandlerFunc
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
		hashids:      hash.New(conf.Get().HashIds.Secret, conf.Get().HashIds.Length),
		adminService: admin.New(db, cache),
	}
}

func (h *handler) i() {}

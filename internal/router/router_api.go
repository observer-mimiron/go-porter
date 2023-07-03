package router

import (
	"go-porter/internal/app/api/admin"
	"go-porter/pkg/core/pkg/net/httpx"
)

//AliasForRecordMetrics 别名 用于记录 metrics
//WrapAuthHandler 权限验证
func setApiRouter(r *resource) {
	// admin
	adminHandler := admin.New(r.logger, r.db, r.cache)

	// 需要签名验证，无需登录验证，
	login := r.mux.Group("/api")
	{
		login.POST("/login", adminHandler.Login())
	}

	// 需要签名验证、登录验证
	notRBAC := r.mux.Group("/api", httpx.WrapAuthHandler(r.interceptors.CheckLogin))
	{
		notRBAC.POST("/admin/logout", adminHandler.Logout())
		notRBAC.PATCH("/admin/modify_password", adminHandler.ModifyPassword())
		notRBAC.GET("/admin/info", adminHandler.Detail())
		notRBAC.PATCH("/admin/modify_personal_info", adminHandler.ModifyPersonalInfo())
	}
}

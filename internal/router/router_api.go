package router

import (
	"go-porter/internal/app/api/admin"
	"go-porter/internal/svc"
	"go-porter/pkg/core/pkg/net/httpx"
)

/**
	AliasForRecordMetrics 别名 用于记录 metrics
	WrapAuthHandler 权限验证
**/
func SetApiRouter(svcCtx *svc.ServiceContext, mux httpx.Mux) {
	// admin
	adminHandler := admin.New(svcCtx)

	// 需要签名验证，无需登录验证，
	login := mux.Group("/api")
	{
		login.POST("/login", adminHandler.Login())
	}

	// 需要签名验证、登录验证
	//	notRBAC := r.mux.Group("/api", httpx.WrapAuthHandler(r.authenticates))
	notRBAC := mux.Group("/api")
	{
		notRBAC.POST("/admin/logout", adminHandler.Logout())
		notRBAC.PATCH("/admin/modify_password", adminHandler.ModifyPassword())
		notRBAC.GET("/admin/info", adminHandler.Detail())
		notRBAC.PATCH("/admin/modify_personal_info", adminHandler.ModifyPersonalInfo())
	}
}

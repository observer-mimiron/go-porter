package router

import (
	"go-porter/internal/app/api/admin"
	"go-porter/internal/svc"
)

/**
	AliasForRecordMetrics 别名 用于记录 metrics
	WrapAuthHandler 权限验证
**/
func SetApiRouter(svc *svc.ServiceContext) {
	// admin
	adminHandler := admin.New(svc.Logger, svc.Db, svc.Redis)

	// 需要签名验证，无需登录验证，
	login := svc.Mux.Group("/api")
	{
		login.POST("/login", adminHandler.Login())
	}

	// 需要签名验证、登录验证
	//	notRBAC := r.mux.Group("/api", httpx.WrapAuthHandler(r.authenticates))
	notRBAC := svc.Mux.Group("/api")
	{
		notRBAC.POST("/admin/logout", adminHandler.Logout())
		notRBAC.PATCH("/admin/modify_password", adminHandler.ModifyPassword())
		notRBAC.GET("/admin/info", adminHandler.Detail())
		notRBAC.PATCH("/admin/modify_personal_info", adminHandler.ModifyPersonalInfo())
	}
}

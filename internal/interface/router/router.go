package router

import (
	"github.com/observer-mimiron/go-porter/internal/interface/handler/admin"
	"github.com/observer-mimiron/go-porter/internal/interface/middleware"
	"github.com/observer-mimiron/go-porter/internal/svc"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/net/httpx"
)

/*
*

	AliasForRecordMetrics 别名 用于记录 metrics
	WrapAuthHandler 权限验证

*
*/
func SetApiRouter(svcCtx *svc.ServiceContext, mux httpx.Mux) {
	// admin
	adminHandler := admin.New(svcCtx)
	auth := middleware.New(svcCtx.Logger, svcCtx.Redis, svcCtx.Db)

	// 需要签名验证，无需登录验证，
	login := mux.Group("/api")
	{
		login.POST("/login", adminHandler.Login())
	}

	// 需要签名验证、登录验证
	notRBAC := mux.Group("/api", httpx.WrapAuthHandler(auth.CheckLogin))
	{
		notRBAC.POST("/admin/logout", adminHandler.Logout())
		notRBAC.PATCH("/admin/modify_password", adminHandler.ModifyPassword())
		notRBAC.GET("/admin/detail", adminHandler.Detail())
		notRBAC.PATCH("/admin/modify_personal_info", adminHandler.ModifyPersonalInfo())
	}
}

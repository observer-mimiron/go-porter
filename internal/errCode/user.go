package errCode

var (
	// client errors
	ErrValidation        = NewErrCode(400001, "验证失败")
	ErrTooManyRequests   = NewErrCode(400002, "请求过多")
	ErrParamBind         = NewErrCode(400003, "参数信息错误")
	ErrUserNotFound      = NewErrCode(402004, "用户不存在")
	ErrRecordNotFound    = NewErrCode(400004, "记录不存在")
	ErrIPLimited         = NewErrCode(402001, "IP受限")
	ErrAuthorization     = NewErrCode(402000, "签名信息错误")
	ErrTokenInvalid      = NewErrCode(402001, "用户验证失败")
	ErrTokenExpired      = NewErrCode(402002, "用户验证过期")
	ErrPasswordIncorrect = NewErrCode(402003, "密码错误")

	// business errors:
	//通用错误 需要warp错误信息

	ErrUrlSign       = NewErrCode(200105, "参数签名错误")
	ErrCacheSet      = NewErrCode(200105, "设置缓存失败")
	ErrCacheGet      = NewErrCode(200105, "获取缓存失败")
	ErrCacheDel      = NewErrCode(200105, "删除缓存失败")
	ErrCacheNotExist = NewErrCode(200105, "缓存不存在")
	ErrResubmit      = NewErrCode(200105, "请勿重复提交")

	ErrWriteConfig   = NewErrCode(200105, "写入配置文件失败")
	ErrSendEmail     = NewErrCode(200105, "发送邮件失败")
	ErrMySQLExec     = NewErrCode(200105, "SQL 执行失败")
	ErrSocketConnect = NewErrCode(200105, "Socket 未连接")
	ErrSocketSend    = NewErrCode(200105, "Socket 消息发送失败")

	ErrAdminResetPassword      = NewErrCode(200105, "重置密码失败")
	ErrAdminLogin              = NewErrCode(200105, "登录失败")
	ErrAdminLogOut             = NewErrCode(200105, "退出失败")
	ErrAdminModifyPassword     = NewErrCode(200105, "修改密码失败")
	ErrAdminModifyPersonalInfo = NewErrCode(200105, "修改个人信息失败")
	ErrAdminMenuList           = NewErrCode(200105, "获取管理员菜单授权列表失败")
	ErrAdminMenuCreate         = NewErrCode(200105, "管理员菜单授权失败")
	ErrAdminOffline            = NewErrCode(200105, "下线管理员失败")
	ErrAdminDetail             = NewErrCode(200105, "获取个人信息失败")
)

package ecode

import "go-porter/pkg/core/pkg/xerror"

//#### 错误码为 6 位数
//| 10 | 01 | 01 |
//| :------ | :------ | :------ |
//| 服务级错误码 | 模块级错误码 | 具体错误码 |
//- 服务级错误码：1 位数进行表示，比如 10 为系统级错误；20 为普通错误，通常是由用户非法操作引起。
//- 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为订单模块。
//- 具体的错误码：2 位数进行表示，比如 01 为手机号不合法；02 为验证码输入错误。

//异常
var (
	// server errors
	ErrInternalServer = xerror.Error(500001, "Server error.")
	ErrDatabase       = xerror.Error(500002, "Database error.")
	ErrJobExpired     = xerror.Error(500004, "Job expired.")
	ErrRecordNotFound = xerror.Error(500004, "record not found.")
	ErrRedis          = xerror.Error(500005, "redis error.")
)

//客户端异常
var (
	// client errors
	ErrValidation        = xerror.Error(400001, "Validation failed.")
	ErrThrottler         = xerror.Error(400002, "Requests are too frequent, please request later. ")
	ErrIllegalRequest    = xerror.Error(400003, "Illegal request. ")
	ErrTooManyRequests   = xerror.Error(400002, "请求过多")
	ErrParamBind         = xerror.Error(400003, "参数信息错误")
	ErrIPLimited         = xerror.Error(402001, "IP limited.")
	ErrTokenInvalid      = xerror.Error(402001, "The token was invalid.")
	ErrTokenExpired      = xerror.Error(402002, "The token was expired.")
	ErrPasswordIncorrect = xerror.Error(402003, "password was incorrect.")
	ErrUserNotFound      = xerror.Error(402004, "User was not existed.")
	ErrNoScore           = xerror.Error(402005, "You haven't scored yet .")
	ErrFoodNotFound      = xerror.Error(402005, "Food was not existed.")
	ErrParams            = xerror.Error(402006, "params error.")
	ErrForbidden         = xerror.Error(402007, "Forbidden")
	ErrFileNotFound      = xerror.Error(404001, "File was not existed.")
)

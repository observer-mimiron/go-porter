package code

import (
	_ "embed"
)

//#### 错误码为 5 位数
//| 1 | 01 | 01 |
//| :------ | :------ | :------ |
//| 服务级错误码 | 模块级错误码 | 具体错误码 |
//- 服务级错误码：1 位数进行表示，比如 1 为系统级错误；2 为普通错误，通常是由用户非法操作引起。
//- 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为订单模块。
//- 具体的错误码：2 位数进行表示，比如 01 为手机号不合法；02 为验证码输入错误。

// Failure 错误时返回结构
type Failure struct {
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
}

const (
	ServerError        = 10101
	TooManyRequests    = 10102
	ParamBindError     = 10103
	AuthorizationError = 10104
	UrlSignError       = 10105
	CacheSetError      = 10106
	CacheGetError      = 10107
	CacheDelError      = 10108
	CacheNotExist      = 10109
	ResubmitError      = 10110
	HashIdsEncodeError = 10111
	HashIdsDecodeError = 10112
	RBACError          = 10113
	RedisConnectError  = 10114
	MySQLConnectError  = 10115
	WriteConfigError   = 10116
	SendEmailError     = 10117
	MySQLExecError     = 10118
	GoVersionError     = 10119
	SocketConnectError = 10120
	SocketSendError    = 10121

	AuthorizedCreateError    = 20101
	AuthorizedListError      = 20102
	AuthorizedDeleteError    = 20103
	AuthorizedUpdateError    = 20104
	AuthorizedDetailError    = 20105
	AuthorizedCreateAPIError = 20106
	AuthorizedListAPIError   = 20107
	AuthorizedDeleteAPIError = 20108

	AdminCreateError             = 20201
	AdminListError               = 20202
	AdminDeleteError             = 20203
	AdminUpdateError             = 20204
	AdminResetPasswordError      = 20205
	AdminLoginError              = 20206
	AdminLogOutError             = 20207
	AdminModifyPasswordError     = 20208
	AdminModifyPersonalInfoError = 20209
	AdminMenuListError           = 20210
	AdminMenuCreateError         = 20211
	AdminOfflineError            = 20212
	AdminDetailError             = 20213

	MenuCreateError       = 20301
	MenuUpdateError       = 20302
	MenuListError         = 20303
	MenuDeleteError       = 20304
	MenuDetailError       = 20305
	MenuCreateActionError = 20306
	MenuListActionError   = 20307
	MenuDeleteActionError = 20308

	CronCreateError  = 20401
	CronUpdateError  = 20402
	CronListError    = 20403
	CronDetailError  = 20404
	CronExecuteError = 20405
)

var CodeText = map[int]string{
	ServerError:        "内部服务器错误",
	TooManyRequests:    "请求过多",
	ParamBindError:     "参数信息错误",
	AuthorizationError: "签名信息错误",
	UrlSignError:       "参数签名错误",
	CacheSetError:      "设置缓存失败",
	CacheGetError:      "获取缓存失败",
	CacheDelError:      "删除缓存失败",
	CacheNotExist:      "缓存不存在",
	ResubmitError:      "请勿重复提交",
	HashIdsEncodeError: "HashID 加密失败",
	HashIdsDecodeError: "HashID 解密失败",
	RBACError:          "暂无访问权限",
	RedisConnectError:  "Redis 连接失败",
	MySQLConnectError:  "MySQL 连接失败",
	WriteConfigError:   "写入配置文件失败",
	SendEmailError:     "发送邮件失败",
	MySQLExecError:     "SQL 执行失败",
	GoVersionError:     "Go 版本不满足要求",
	SocketConnectError: "Socket 未连接",
	SocketSendError:    "Socket 消息发送失败",

	AuthorizedCreateError:    "创建调用方失败",
	AuthorizedListError:      "获取调用方列表失败",
	AuthorizedDeleteError:    "删除调用方失败",
	AuthorizedUpdateError:    "更新调用方失败",
	AuthorizedDetailError:    "获取调用方详情失败",
	AuthorizedCreateAPIError: "创建调用方 API 地址失败",
	AuthorizedListAPIError:   "获取调用方 API 地址列表失败",
	AuthorizedDeleteAPIError: "删除调用方 API 地址失败",

	AdminCreateError:             "创建管理员失败",
	AdminListError:               "获取管理员列表失败",
	AdminDeleteError:             "删除管理员失败",
	AdminUpdateError:             "更新管理员失败",
	AdminResetPasswordError:      "重置密码失败",
	AdminLoginError:              "登录失败",
	AdminLogOutError:             "退出失败",
	AdminModifyPasswordError:     "修改密码失败",
	AdminModifyPersonalInfoError: "修改个人信息失败",
	AdminMenuListError:           "获取管理员菜单授权列表失败",
	AdminMenuCreateError:         "管理员菜单授权失败",
	AdminOfflineError:            "下线管理员失败",
	AdminDetailError:             "获取个人信息失败",

	MenuCreateError:       "创建菜单失败",
	MenuUpdateError:       "更新菜单失败",
	MenuDeleteError:       "删除菜单失败",
	MenuListError:         "获取菜单列表失败",
	MenuDetailError:       "获取菜单详情失败",
	MenuCreateActionError: "创建菜单栏功能权限失败",
	MenuListActionError:   "获取菜单栏功能权限列表失败",
	MenuDeleteActionError: "删除菜单栏功能权限失败",

	CronCreateError:  "创建后台任务失败",
	CronUpdateError:  "更新后台任务失败",
	CronListError:    "获取定时任务列表失败",
	CronDetailError:  "获取定时任务详情失败",
	CronExecuteError: "手动执行定时任务失败",
}

func Text(code int) string {
	return CodeText[code]
}

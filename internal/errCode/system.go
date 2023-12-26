package errCode

//#### 错误码为 6 位数
//| 10 | 01 | 01 |
//| :------ | :------ | :------ |
//| 服务级错误码 | 模块级错误码 | 具体错误码 |
//- 服务级错误码：1 位数进行表示，比如 10 为系统级错误；20 为普通错误，通常是由用户非法操作引起。
//- 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为订单模块。
//- 具体的错误码：2 位数进行表示，比如 01 为手机号不合法；02 为验证码输入错误。

const Alert = 1

//异常
var (
	// server errors
	ErrServer         = NewErrCode(500001, "服务器异常").Alert()
	ErrDatabase       = NewErrCode(500002, "数据库异常").Alert()
	ErrRedis          = NewErrCode(500005, "Redis异常").Alert()
	ErrForbidden      = NewErrCode(500006, "非法访问").Alert()
	ErrIllegalRequest = NewErrCode(500007, "非法请求").Alert()
)

package configs

import "time"

const (
	// ProjectName 项目名称
	ProjectName = "porter"

	// ProjectPort 项目端口
	ProjectPort = ":9999"

	// HeaderLoginToken 登录验证 Token，Header 中传递的参数
	HeaderLoginToken = "Token"

	// HeaderSignToken 签名验证 Authorization，Header 中传递的参数
	HeaderSignToken = "Authorization"

	// HeaderSignTokenDate 签名验证 Date，Header 中传递的参数
	HeaderSignTokenDate = "Authorization-Date"

	// HeaderSignTokenTimeout 签名有效期为 2 分钟
	HeaderSignTokenTimeout = time.Minute * 2

	// RedisKeyPrefixLoginUser Redis Key 前缀 - 登录用户信息
	RedisKeyPrefixLoginUser = ProjectName + ":login-user:"

	// RedisKeyPrefixSignature Redis Key 前缀 - 签名验证信息
	RedisKeyPrefixSignature = ProjectName + ":signature:"

	// MaxRequestsPerSecond 每秒最大请求量
	MaxRequestsPerSecond = 10000

	// LoginSessionTTL 登录有效期为 24 小时
	LoginSessionTTL = time.Hour * 24
)

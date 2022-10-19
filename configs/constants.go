package configs

import "time"

const (
	// ProjectName 项目名称
	ProjectName = "porter"

	// ProjectDomain 项目域名
	ProjectDomain = "http://127.0.0.1"

	// ProjectPort 项目端口
	ProjectPort = ":9999"

	// ProjectAccessLogFile 项目访问日志存放文件
	ProjectAccessLogFile = "./logs/" + ProjectName + "-access.log"

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

	// ZhCN 简体中文 - 中国
	ZhCN = "zh-cn"

	// MaxRequestsPerSecond 每秒最大请求量
	MaxRequestsPerSecond = 10000

	// LoginSessionTTL 登录有效期为 24 小时
	LoginSessionTTL = time.Hour * 24
)

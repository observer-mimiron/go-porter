package model

import (
	"time"
)

type Admin struct {
	Id          int32     // 主键
	Username    string    // 用户名
	Password    string    // 密码
	Nickname    string    // 昵称
	Mobile      string    // 手机号
	IsUsed      int32     // 是否启用 1:是  -1:否
	IsDeleted   int32     // 是否删除 1:是  -1:否
	CreatedUser string    // 创建人
	UpdatedUser string    // 更新人
	CreatedAt   time.Time `gorm:"time"` // 创建时间
	UpdatedAt   time.Time `gorm:"time"` // 更新时间
}

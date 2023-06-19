package model

import (
	"fmt"
	"github.com/pkg/errors"
	"go-porter/pkg/core/pkg/database/mysql"
	"gorm.io/gorm"
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

func NewModel() *Admin {
	return new(Admin)
}

func NewQueryBuilder() *adminQueryBuilder {
	return new(adminQueryBuilder)
}

func (t *Admin) Create(db *gorm.DB) (id int32, err error) {
	if err = db.Create(t).Error; err != nil {
		return 0, errors.Wrap(err, "create err")
	}
	return t.Id, nil
}

type adminQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func (qb *adminQueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
	ret := db
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	for _, order := range qb.order {
		ret = ret.Order(order)
	}
	ret = ret.Limit(qb.limit).Offset(qb.offset)
	return ret
}

func (qb *adminQueryBuilder) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	db = db.Model(&Admin{})

	for _, where := range qb.where {
		db.Where(where.prefix, where.value)
	}

	if err = db.Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *adminQueryBuilder) Delete(db *gorm.DB) (err error) {
	for _, where := range qb.where {
		db = db.Where(where.prefix, where.value)
	}

	if err = db.Delete(&Admin{}).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}

func (qb *adminQueryBuilder) Count(db *gorm.DB) (int64, error) {
	var c int64
	res := qb.buildQuery(db).Model(&Admin{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (qb *adminQueryBuilder) First(db *gorm.DB) (*Admin, error) {
	ret := &Admin{}
	res := qb.buildQuery(db).First(ret)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		ret = nil
	}
	return ret, res.Error
}

func (qb *adminQueryBuilder) QueryOne(db *gorm.DB) (*Admin, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return ret[0], err
	}
	return nil, err
}

func (qb *adminQueryBuilder) QueryAll(db *gorm.DB) ([]*Admin, error) {
	var ret []*Admin
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}

func (qb *adminQueryBuilder) Limit(limit int) *adminQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *adminQueryBuilder) Offset(offset int) *adminQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *adminQueryBuilder) WhereId(p mysql.Predicate, value int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereIdIn(value []int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereIdNotIn(value []int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderById(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "id "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereUsername(p mysql.Predicate, value string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "username", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereUsernameIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "username", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereUsernameNotIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "username", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByUsername(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "username "+order)
	return qb
}

func (qb *adminQueryBuilder) WherePassword(p mysql.Predicate, value string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WherePasswordIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WherePasswordNotIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "password", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByPassword(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "password "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereNickname(p mysql.Predicate, value string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "nickname", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereNicknameIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "nickname", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereNicknameNotIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "nickname", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByNickname(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "nickname "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereMobile(p mysql.Predicate, value string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "mobile", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereMobileIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "mobile", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereMobileNotIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "mobile", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByMobile(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "mobile "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereIsUsed(p mysql.Predicate, value int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_used", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereIsUsedIn(value []int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_used", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereIsUsedNotIn(value []int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_used", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByIsUsed(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "is_used "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereIsDeleted(p mysql.Predicate, value int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereIsDeletedIn(value []int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereIsDeletedNotIn(value []int32) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByIsDeleted(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "is_deleted "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereCreatedAt(p mysql.Predicate, value time.Time) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_at", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereCreatedAtIn(value []time.Time) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_at", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereCreatedAtNotIn(value []time.Time) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_at", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByCreatedAt(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "created_at "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereCreatedUser(p mysql.Predicate, value string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_user", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereCreatedUserIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_user", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereCreatedUserNotIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_user", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByCreatedUser(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "created_user "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereUpdatedAt(p mysql.Predicate, value time.Time) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_at", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereUpdatedAtIn(value []time.Time) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_at", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereUpdatedAtNotIn(value []time.Time) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_at", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByUpdatedAt(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "updated_at "+order)
	return qb
}

func (qb *adminQueryBuilder) WhereUpdatedUser(p mysql.Predicate, value string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_user", p),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereUpdatedUserIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_user", "IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) WhereUpdatedUserNotIn(value []string) *adminQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_user", "NOT IN"),
		value,
	})
	return qb
}

func (qb *adminQueryBuilder) OrderByUpdatedUser(asc bool) *adminQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "updated_user "+order)
	return qb
}

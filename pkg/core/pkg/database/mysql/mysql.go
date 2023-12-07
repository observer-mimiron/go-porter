package mysql

import (
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"

	"github.com/pkg/errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Conf struct {
	Base `toml:"base"`
	Read struct {
		Addr string `toml:"addr"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		Name string `toml:"name"`
	} `toml:"read"`
	Write struct {
		Addr string `toml:"addr"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		Name string `toml:"name"`
	} `toml:"write"`
}

type Base struct {
	MaxOpenConn     int           `toml:"maxOpenConn"`
	MaxIdleConn     int           `toml:"maxIdleConn"`
	ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
}

// Predicate is a string that acts as a condition in the where clause
type Predicate string

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDbR() *gorm.DB
	GetDbW() *gorm.DB
	DbRClose() error
	DbWClose() error
}

type dbRepo struct {
	DbR *gorm.DB
	DbW *gorm.DB
}

func New(cfg Conf) (Repo, error) {
	dbr, err := dbConnect(cfg.Base, cfg.Read.User, cfg.Read.Pass, cfg.Read.Addr, cfg.Read.Name)
	if err != nil {
		panic(err)
		return nil, err
	}

	dbw, err := dbConnect(cfg.Base, cfg.Write.User, cfg.Write.Pass, cfg.Write.Addr, cfg.Write.Name)
	if err != nil {
		panic(err)
		return nil, err
	}
	return &dbRepo{
		DbR: dbr,
		DbW: dbw,
	}, nil
}

func (d *dbRepo) i() {}

func (d *dbRepo) GetDbR() *gorm.DB {
	return d.DbR
}

func (d *dbRepo) GetDbW() *gorm.DB {
	return d.DbW
}

func (d *dbRepo) DbRClose() error {
	sqlDB, err := d.DbR.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *dbRepo) DbWClose() error {
	sqlDB, err := d.DbW.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func dbConnect(base Base, user, pass, addr, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		user,
		pass,
		addr,
		dbName,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info), // 日志配置
	})

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", dbName))
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(base.MaxOpenConn)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(base.MaxIdleConn)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Minute * base.ConnMaxLifeTime)

	// 使用插件
	err = db.Use(&TracePlugin{})
	if err != nil {
		return nil, err
	}

	//err = db.Use(prometheus.New(prometheus.Config{
	//	DBName:          "db1",                       // 使用 `DBName` 作为指标 label
	//	RefreshInterval: 15,                          // 指标刷新频率（默认为 15 秒）
	//	PushAddr:        "prometheus pusher address", // 如果配置了 `PushAddr`，则推送指标
	//	StartServer:     true,                        // 启用一个 http 服务来暴露指标
	//	HTTPServerPort:  8080,                        // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
	//	MetricsCollector: []prometheus.MetricsCollector{
	//		&prometheus.MySQL{
	//			VariableNames: []string{"Threads_running"},
	//		},
	//	}, // 用户自定义指标
	//}))
	//if err != nil {
	//	return nil, err
	//}
	return db, nil
}

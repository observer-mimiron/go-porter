package conf

import (
	_ "embed"
	"go-porter/pkg/core/pkg/logger"
	"os"
	"path/filepath"
	"time"

	"go-porter/pkg/file"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	Name string `toml:"projectName"`
	Host string `toml:"host"`
	Port string `toml:"port"`

	MySQL `toml:"mysql"`

	Redis `toml:"redis"`

	HashIds `toml:"hashids"`

	Log logger.LogConf `toml:"log"`
}

type MySQL struct {
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
	Base struct {
		MaxOpenConn     int           `toml:"maxOpenConn"`
		MaxIdleConn     int           `toml:"maxIdleConn"`
		ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
	} `toml:"base"`
}

type Redis struct {
	Addr         string `toml:"addr"`
	Pass         string `toml:"pass"`
	Db           int    `toml:"db"`
	MaxRetries   int    `toml:"maxRetries"`
	PoolSize     int    `toml:"poolSize"`
	MinIdleConns int    `toml:"minIdleConns"`
}

type HashIds struct {
	Secret string `toml:"secret"`
	Length int    `toml:"length"`
}

//func init() {
//	var r io.Reader
//
//	r = bytes.NewReader(Configs)
//
//	viper.SetConfigType("toml")
//
//	if err := viper.ReadConfig(r); err != nil {
//		panic(err)
//	}
//
//	if err := viper.Unmarshal(config); err != nil {
//		panic(err)
//	}
//
//	viper.SetConfigName(env.Active().Value() + "_configs")
//	viper.AddConfigPath("./configs")
//
//	configFile := "./configs/" + env.Active().Value() + "_configs.toml"
//	_, ok := file.IsExists(configFile)
//	if !ok {
//		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
//			panic(err)
//		}
//
//		f, err := os.Create(configFile)
//		if err != nil {
//			panic(err)
//		}
//		defer f.Close()
//
//		if err := viper.WriteConfig(); err != nil {
//			panic(err)
//		}
//	}
//
//	viper.WatchConfig()
//	viper.OnConfigChange(func(e fsnotify.Event) {
//		if err := viper.Unmarshal(config); err != nil {
//			panic(err)
//		}
//	})
//}

func Get() Config {
	return *config
}

func Init(configFile *string) {

	if !filepath.IsAbs(*configFile) {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		*configFile = filepath.Join(dir, *configFile)
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(*configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	_, ok := file.IsExists(*configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(*configFile), 0766); err != nil {
			panic(err)
		}

		f, err := os.Create(*configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 解析配置文件
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

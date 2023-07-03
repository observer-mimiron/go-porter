package conf

import (
	_ "embed"
	"go-porter/pkg/core/pkg/cache/redis"
	"go-porter/pkg/core/pkg/database/mysql"
	"go-porter/pkg/core/pkg/logger"
	"go-porter/pkg/cryptor/hash"
	"go-porter/pkg/file"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	Name string `toml:"projectName"`
	Host string `toml:"host"`
	Port string `toml:"port"`

	HashIds hash.Conf `toml:"hashids"`

	MySQL mysql.Conf `toml:"mysql"`

	Redis redis.Conf `toml:"redis"`

	Log logger.Conf `toml:"log"`
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

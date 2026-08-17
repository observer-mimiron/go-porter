package conf

import (
	"fmt"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/cache/redis"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/database/mysql"
	"github.com/observer-mimiron/go-porter/pkg/core/pkg/logger"
	"github.com/observer-mimiron/go-porter/pkg/cryptor/hash"
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

// Get returns the latest loaded configuration snapshot.
func Get() Config {
	return *config
}

// Init loads a TOML configuration file and watches it for changes.
func Init(configFile string) error {
	if !filepath.IsAbs(configFile) {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		configFile = filepath.Join(dir, configFile)
	}
	if _, err := os.Stat(configFile); err != nil {
		return fmt.Errorf("stat config file %q: %w", configFile, err)
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 解析配置文件
		if err := viper.Unmarshal(config); err != nil {
			return
		}
	})
	return nil
}

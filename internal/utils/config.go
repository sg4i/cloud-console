package utils

import (
	"strings"
	"sync"
	"time"

	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/spf13/viper"
)

var (
	once sync.Once
)

// LoadConfig 加载配置文件，configFile 为配置文件的完整路径（包含文件名）
func LoadConfig(configFile string) {
	// sync.Once 来确保 LoadConfig 函数中的初始化代码只会执行一次
	once.Do(func() {
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			viper.SetConfigName("config")
			viper.SetConfigType("yml")
			viper.AddConfigPath(".")
		}

		if err := viper.ReadInConfig(); err != nil {
			logger.Log.WithError(err).Error("读取配置文件失败")
		}

		viper.AutomaticEnv()
		// 使用 viper.SetEnvKeyReplacer 来设置环境变量名称的转换规则。
		// 例如，配置文件中的 storage.mongo.uri 对应的环境变量名称将是 STORAGE_MONGO_URI
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	})
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetUint64(key string) uint64 {
	return viper.GetUint64(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

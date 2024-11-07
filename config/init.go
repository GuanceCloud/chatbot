package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigType("yml")  // 设置配置文件类型
	viper.AddConfigPath("./config/")    // 设置配置文件搜索路径
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}
}

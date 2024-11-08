package main

import (
	"github.com/GuanceCloud/chatbot/config"
	"github.com/GuanceCloud/chatbot/initialize"
	"github.com/GuanceCloud/chatbot/utils"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// Redis初始化
	utils.InitRedis()
	defer utils.CloseRedis()

	r := initialize.SetupRouter()
	serverPort := viper.GetString("server.port")
	r.Run(serverPort)
}

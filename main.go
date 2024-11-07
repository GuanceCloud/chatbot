package main

import (
	"github.com/GuanceCloud/chatbot/config"
	"github.com/GuanceCloud/chatbot/initialize"
	"github.com/GuanceCloud/chatbot/utils"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// Redis初始化
	utils.InitRedis()
	defer utils.CloseRedis()

	r := initialize.SetupRouter()
	r.Run(":80")
}

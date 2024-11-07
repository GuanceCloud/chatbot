package main

import (
	"github.com/GuanceCloud/chatbot/config"
	"github.com/GuanceCloud/chatbot/initialize"
)

func main() {
	// err := initialize.InitMySQL()
	// if err != nil {
	// 	panic(err)
	// }
	// defer initialize.Close()

	// 初始化配置
	config.InitConfig()
	r := initialize.SetupRouter()
	r.Run(":80")
}
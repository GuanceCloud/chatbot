package main

import (
	"github.com/GuanceCloud/chatbot/initialize"
)

func main() {
	// err := initialize.InitMySQL()
	// if err != nil {
	// 	panic(err)
	// }
	// defer initialize.Close()

	r := initialize.SetupRouter()
	r.Run(":80")
}

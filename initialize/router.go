/*
 * @Author: 安俊霖
 * @Date: 2024-11-06 20:10:09
 * @Description:
 */
package initialize

import (
	"github.com/GuanceCloud/chatbot/middleware"
	"github.com/GuanceCloud/chatbot/router"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.Auth())
	Group := r.Group("")
	{
		router.InitUserRouter(Group)
		router.InitChatRouter(Group)
	}
	return r
}

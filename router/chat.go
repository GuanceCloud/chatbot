package router

import (
	api "github.com/GuanceCloud/chatbot/api"
	"github.com/gin-gonic/gin"
)

func InitChatRouter(Router *gin.RouterGroup) {
	ChatRouter := Router.Group("")
	{
		ChatRouter.POST("/smart_query_stream", api.SmartQueryStream)
	}
}

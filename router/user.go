/*
 * @Author: 安俊霖
 * @Date: 2024-11-06 20:13:42
 * @Description:
 */
package router

import (
	api "github.com/GuanceCloud/chatbot/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("")
	{
		UserRouter.POST("/get_token", api.GetToken)
	}
}

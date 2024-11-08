/* Copyright 2024 GuanceCloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/GuanceCloud/chatbot/pkg/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否应该跳过中间件
		if c.Request.URL.Path == "/get_token" {
			c.Next()
			return
		}

		// 从请求头中获取 Authorization 字段
		authHeader := c.GetHeader("Authorization")
		// 检查 Authorization 字段是否存在且以 "Bearer" 开头
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"retcode": -40001,
				"message": "Authorization header is required",
				"data":    gin.H{},
			})
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 验证 token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"retcode": -40002,
				"message": "Invalid or expired token",
				"data":    gin.H{},
			})
			c.Abort()
			return
		}

		// 如果验证成功，将用户信息添加到上下文（如果需要）
		c.Set("user_id", claims["user_id"])
		c.Next() // 继续处理请求
	}
}

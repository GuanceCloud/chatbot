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

package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GuanceCloud/chatbot/pkg/utils"
)

// @description 获取token
// @Tags user
// @Param user_id formData string true "用户id"
// @Router /get_token [post]
func GetToken(c *gin.Context) {
	var data struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "user_id is required",
			"data":    gin.H{},
		})
		return
	}

	userID := data.UserID
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "user_id is required",
			"data":    gin.H{},
		})
		return
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	log.Printf("Generate token: '%s' with user_id: '%s'", token, userID)

	c.JSON(http.StatusOK, gin.H{
		"retcode": 0,
		"message": "success",
		"data": gin.H{
			"token": token,
		},
	})
}

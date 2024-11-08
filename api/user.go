package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GuanceCloud/chatbot/utils"
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

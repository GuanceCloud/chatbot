/*
 * @Author: 安俊霖
 * @Date: 2024-11-06 20:20:58
 * @Description:
 */
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type DifyRequest struct {
	Inputs         map[string]interface{} `json:"inputs"`
	Query          string                 `json:"query"`
	ResponseMode   string                 `json:"response_mode"`
	ConversationID string                 `json:"conversation_id"`
	User           string                 `json:"user"`
	Files          []map[string]string    `json:"files"`
}

// @Router /smart_query_stream
func SmartQueryStream(c *gin.Context) {
	var requestBody struct {
		Query  string `json:"query"`
		UserId string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "Invalid request",
			"data":    gin.H{},
		})
		return
	}

	// 这里是从token获取的，从请求体里获取也可以
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "user_id is required",
			"data":    gin.H{},
		})
		return
	}

	query := requestBody.Query
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "query is required",
			"data":    gin.H{},
		})
		return
	}

	// Dify配置
	var (
		difyBaseURL = viper.GetString("dify.baseURL")
		apiKey      = viper.GetString("dify.apiKey")
	)
	// dify 请求
	difyReq := DifyRequest{
		Inputs:         map[string]interface{}{},
		Query:          query,
		ResponseMode:   "streaming",
		ConversationID: "",
		User:           userID.(string),
	}

	// 将请求编码为 JSON
	reqBody, err := json.Marshal(difyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"retcode": -30000,
			"message": "Failed to marshal request body",
			"data":    gin.H{},
		})
		return
	}

	// 发起请求到 dify 服务
	client := &http.Client{}
	req, err := http.NewRequest("POST", difyBaseURL+"chat-messages", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"retcode": -30000,
			"message": "Failed to create request to dify service",
			"data":    gin.H{},
		})
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("======================")
	fmt.Println(req.Header.Get("Authorization"))
	fmt.Println("======================")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"retcode": -30000,
			"message": "Failed to send request to dify service",
			"data":    gin.H{},
		})
		return
	}
	defer resp.Body.Close()

	// 检查 dify 服务响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{
			"retcode": -30000,
			"message": "Dify service returned an error",
			"data":    string(body),
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 流式传输响应
	done := make(chan bool)
	go func() {
		defer close(done)
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if n == 0 || err != nil {
				return
			}
			c.SSEvent("data", string(buf[:n]))
		}
	}()

	// 等待流式传输完成
	<-done
	c.SSEvent("event", "message_end")
}

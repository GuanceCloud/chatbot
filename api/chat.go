/*
 * @Author: 安俊霖
 * @Date: 2024-11-06 20:20:58
 * @Description:
 */
package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DifyServiceURL 是 dify 服务的基础 URL
const DifyServiceURL = "http://localhost:8000/v1"

// DifyRequest 是发送给 dify 服务的请求结构体
type DifyRequest struct {
	Inputs         map[string]interface{} `json:"inputs"`
	Query          string                 `json:"query"`
	ResponseMode   string                 `json:"response_mode"`
	ConversationID string                 `json:"conversation_id"`
	User           string                 `json:"user"`
	Files          []map[string]string    `json:"files"`
}

// 假设以下是您的其他相关函数和类型
type Document struct {
	PageContent string `json:"page_content"`
	Metadata    struct {
		Source string `json:"source"`
	} `json:"metadata"`
}

func SmartQueryStream(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "user_id is required",
			"data":    gin.H{},
		})
		return
	}

	query := c.PostForm("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"retcode": -20000,
			"message": "query is required",
			"data":    gin.H{},
		})
		return
	}

	// 设置 dify 请求
	difyReq := DifyRequest{
		Inputs:         map[string]interface{}{},
		Query:          query,
		ResponseMode:   "streaming",
		ConversationID: "",
		User:           userID.(string),
		Files: []map[string]string{
			{
				"type":            "image",
				"transfer_method": "remote_url",
				"url":             "https://cloud.dify.ai/logo/logo-site.png",
			},
		},
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
	resp, err := http.Post(DifyServiceURL+"/chat-messages", "application/json", bytes.NewBuffer(reqBody))
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
			if err != nil {
				return
			}
		}
	}()

	// 等待流式传输完成
	<-done
	c.SSEvent("event", "message_end")
}

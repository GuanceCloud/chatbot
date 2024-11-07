/*
 * @Author: 安俊霖
 * @Date: 2024-11-06 20:20:58
 * @Description:
 */
package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/GuanceCloud/chatbot/utils"
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

// dify 响应结构
type DifyResponseMessage struct {
	Event          string `json:"event"`
	TaskID         string `json:"task_id"`
	MessageID      string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Answer         string `json:"answer"`
	CreatedAt      int    `json:"created_at"`
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

	// 从 Redis 获取 conversation_id
	conversationId := utils.GetRedis(userID.(string))

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
		ConversationID: conversationId,
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

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			c.String(http.StatusInternalServerError, "Error reading SSE stream")
			return
		}

		parseConversationId(line, userID.(string))

		// 处理每一行数据
		_, err = c.Writer.Write(line)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error writing to client")
			return
		}
		c.Writer.Flush()
	}
}

func parseConversationId(line []byte, userID string) {
	// 去掉行尾的换行符
	line = bytes.TrimSuffix(line, []byte{'\n'})

	// 检查是否是数据行
	if bytes.HasPrefix(line, []byte("data:")) {
		// 提取JSON数据
		jsonData := line[5:] // 跳过"data: "

		// 定义一个用于解析JSON的变量
		var msg DifyResponseMessage

		// 解析JSON
		err := json.Unmarshal(jsonData, &msg)
		if err == nil {
			// 处理解析后的数据
			if msg.Event == "message" && msg.ConversationID != "" {
				// fmt.Println("********************")
				// fmt.Println("conversation_id: ", msg.ConversationID)
				// fmt.Println("********************")
				// 放到redis中
				if utils.GetRedis(userID) == "" {
					utils.SetRedis(userID, msg.ConversationID, 3600)
				}
			}
		}
	}
}

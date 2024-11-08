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

package utils

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/GuanceCloud/chatbot/internal/config"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis(cfg config.RedisConfig) {
	redisAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	redisDB, err := strconv.Atoi(cfg.DB)
	if err != nil {
		panic(fmt.Sprintf("invalid redis db param value: %s", cfg.DB))
	}

	// 初始化 Redis 客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.Password,
		DB:       redisDB,
		Username: cfg.Username,
	})

	// 验证连接
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("Redis 连接失败：" + err.Error())
	}
}

// CloseRedis 关闭 Redis 客户端
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

func keyWithPrefix(key string) string {
	return "chatbot_api:" + key
}

func SetRedis(key string, value string, t int64) bool {
	key = keyWithPrefix(key)
	expire := time.Duration(t) * time.Second
	if err := redisClient.Set(ctx, key, value, expire).Err(); err != nil {
		return false
	}
	return true
}

func GetRedis(key string) string {
	key = keyWithPrefix(key)
	result, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

func DelRedis(key string) bool {
	key = keyWithPrefix(key)
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ExpireRedis(key string, t int64) bool {
	key = keyWithPrefix(key)
	// 延长过期时间
	expire := time.Duration(t) * time.Second
	if err := redisClient.Expire(ctx, key, expire).Err(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

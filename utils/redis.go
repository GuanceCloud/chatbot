/*
 * @Author: 安俊霖
 * @Date: 2024-11-07 15:48:39
 * @Description:
 */
/*
 * @Author: 安俊霖
 * @Date: 2024-11-07 15:48:39
 * @Description:
 */
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis() {
	// 从 viper 中读取 Redis 配置
	redisAddr := viper.GetString("redis.addr")
	redisPassword := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.db")

	// 初始化 Redis 客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
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

func SetRedis(key string, value string, t int64) bool {
	expire := time.Duration(t) * time.Second
	if err := redisClient.Set(ctx, key, value, expire).Err(); err != nil {
		return false
	}
	return true
}

func GetRedis(key string) string {
	result, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return result
}

func DelRedis(key string) bool {
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ExpireRedis(key string, t int64) bool {
	// 延长过期时间
	expire := time.Duration(t) * time.Second
	if err := redisClient.Expire(ctx, key, expire).Err(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

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

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       string `yaml:"db"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DifyConfig struct {
	BaseURL string `yaml:"baseURL"`
	APIKey  string `yaml:"apiKey"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Redis  RedisConfig  `yaml:"redis"`
	Dify   DifyConfig   `yaml:"dify"`

	GuanceSecret string `yaml:"guanceSecret"`
}

func InitConfig() Config {
	viper.SetConfigType("yml")        // 设置配置文件类型
	viper.AddConfigPath("./.config/") // 设置配置文件搜索路径

	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")

	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("redis.username", "REDIS_USERNAME")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")

	viper.BindEnv("dify.baseURL", "DIFY_BASE_URL")
	viper.BindEnv("dify.apiKey", "DIFY_API_KEY")

	viper.BindEnv("guanceSecret", "GUANCE_SECRET")

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		fmt.Printf("load config file: %s\n", err)
	}

	var config Config
	viper.Unmarshal(&config)
	return config
}

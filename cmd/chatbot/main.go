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

package main

import (
	"github.com/GuanceCloud/chatbot/internal/config"
	"github.com/GuanceCloud/chatbot/internal/initialize"
	"github.com/GuanceCloud/chatbot/pkg/utils"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// Redis初始化
	utils.InitRedis()
	defer utils.CloseRedis()

	r := initialize.SetupRouter()
	serverPort := viper.GetString("server.port")
	r.Run(serverPort)
}

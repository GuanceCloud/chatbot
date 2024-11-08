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

func InitConfig() {
	viper.SetConfigType("yml")       // 设置配置文件类型
	viper.AddConfigPath("./config/") // 设置配置文件搜索路径
	err := viper.ReadInConfig()      // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}
}
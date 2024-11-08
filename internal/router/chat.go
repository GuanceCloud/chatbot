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

package router

import (
	api "github.com/GuanceCloud/chatbot/internal/api"
	"github.com/gin-gonic/gin"
)

func InitChatRouter(Router *gin.RouterGroup) {
	ChatRouter := Router.Group("/open_kf_api/queries")
	{
		ChatRouter.POST("/smart_query_stream", api.SmartQueryStream)
	}
}

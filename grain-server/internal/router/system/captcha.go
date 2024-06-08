// Copyright © 2023 Grain. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/config"
	handler "github.com/go-grain/grain/internal/handler/system"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/middleware"
	redisx "github.com/go-grain/grain/pkg/redis"
)

type CaptchaRouter struct {
	api     *handler.CaptchaHandle
	public  gin.IRoutes
	private gin.IRoutes
}

func NewCaptchaRouter(routerGroup *gin.RouterGroup, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *CaptchaRouter {
	sv := service.NewCaptcha(rdb, conf, logger)
	return &CaptchaRouter{
		api:    handler.NewCaptchaHandle(sv),
		public: routerGroup.Group("captcha"),
		private: routerGroup.Group("captcha").Use(
			middleware.JwtAuth(rdb),
		),
	}
}

func (r *CaptchaRouter) InitRouters() {
	// 发送手机号验证码
	r.public.POST("sendMobileCaptcha", r.api.SendMobileCaptcha)
	// 发送 email 验证码
	r.public.POST("sendEmailCaptcha", r.api.SendEmailCaptcha)
	// 发送用户手机号
	r.private.POST("sendUserMobileCaptcha", r.api.SendUserMobileCaptcha)
	// 发送用户 email 验证码
	r.private.POST("sendUserEmailCaptcha", r.api.SendUserEmailCaptcha)
}

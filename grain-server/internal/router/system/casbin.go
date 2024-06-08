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
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/config"
	handler "github.com/go-grain/grain/internal/handler/system"
	repo "github.com/go-grain/grain/internal/repo/system"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/middleware"
	redisx "github.com/go-grain/grain/pkg/redis"
)

type CasbinRouter struct {
	api     *handler.CasbinHandle
	public  gin.IRoutes
	private gin.IRoutes
}

func NewCasbinRouter(routerGroup *gin.RouterGroup, rdb redisx.IRedis, conf *config.Config, logger log.Logger, enforcer *casbin.CachedEnforcer) *CasbinRouter {
	data := repo.NewCasbinRepo()
	sv := service.NewCasbinService(data, conf, logger, enforcer)
	return &CasbinRouter{
		api:    handler.NewCasbinHandle(sv),
		public: routerGroup.Group(""),
		private: routerGroup.Group("").Use(
			middleware.JwtAuth(rdb),
			middleware.Casbin(enforcer),
		),
	}
}

func (r *CasbinRouter) InitRouters() *CasbinRouter {
	r.private.PUT("casbin", r.api.Update)
	// 获取某角色可访问的api接口列表
	r.private.GET("casbin/authApiList", r.api.AuthApiList)
	return r
}

func (r *CasbinRouter) InitCasbin() *CasbinRouter {
	_ = r.api.InitCasbinHandle()
	return r
}

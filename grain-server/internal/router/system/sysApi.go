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

type ApiRouter struct {
	api     *handler.ApiHandle
	public  gin.IRoutes
	private gin.IRoutes
}

func NewApiRouter(routerGroup *gin.RouterGroup, rdb redisx.IRedis, conf *config.Config, logger log.Logger, enforcer *casbin.CachedEnforcer) *ApiRouter {
	data := repo.NewApiRepo(rdb)
	sv := service.NewApiService(data, rdb, conf, logger)
	return &ApiRouter{
		api:    handler.NewApiHandle(sv),
		public: routerGroup.Group("sysApi"),
		private: routerGroup.
			Group("sysApi").
			Use(
				middleware.JwtAuth(rdb),
				middleware.Casbin(enforcer),
			),
	}
}

func (r *ApiRouter) InitRouters() *ApiRouter {
	r.private.PUT("", r.api.UpdateApi)
	r.private.POST("", r.api.CreateApi)
	r.private.GET("list", r.api.GetApiList)
	r.private.DELETE("", r.api.DeleteApiById)
	r.private.GET("apiGroups", r.api.GetApiGroup)
	r.private.DELETE("deleteApiByIds", r.api.DeleteApiByIds)
	r.private.GET("apiAndPermissions", r.api.GetApiAndPermissions)
	return r
}

func (r *ApiRouter) InitApi() *ApiRouter {
	_ = r.api.InitApi()
	return r
}

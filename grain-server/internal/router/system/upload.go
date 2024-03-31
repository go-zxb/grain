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
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	handler "github.com/go-grain/grain/internal/handler/system"
	repo "github.com/go-grain/grain/internal/repo/system"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/middleware"
)

type UploadRouter struct {
	engine          *gin.Engine
	public          gin.IRoutes
	private         gin.IRoutes
	privateRoleAuth gin.IRoutes
	api             *handler.UploadHandle
}

func NewUploadRouter(routerGroup *gin.RouterGroup, engine *gin.Engine, rdb redis.IRedis, conf *config.Config, logger log.Logger, enforcer *casbin.CachedEnforcer) *UploadRouter {
	data := repo.NewUploadRepo(rdb)
	sv := service.NewUploadService(data, rdb, conf, logger)
	return &UploadRouter{
		engine: engine,
		public: routerGroup.Group("upload"),
		api:    handler.NewUploadHandle(sv),
		private: routerGroup.Group("upload").Use(
			middleware.JwtAuth(rdb),
			middleware.Casbin(enforcer),
		),
	}
}

func (r *UploadRouter) InitRouters() {
	r.engine.Static("uploads", "uploads")
	r.private.POST("", r.api.UploadFile)
	r.private.GET("list", r.api.GetUploadList)
	r.private.DELETE("", r.api.DeleteUploadById)
	r.private.DELETE("deleteUploadByIds", r.api.DeleteUploadByIds)
}

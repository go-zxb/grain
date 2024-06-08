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

type RoleRouter struct {
	api     *handler.RoleHandle
	public  gin.IRoutes
	private gin.IRoutes
}

func NewRoleRouter(routerGroup *gin.RouterGroup, rdb redisx.IRedis, conf *config.Config, logger log.Logger, enforcer *casbin.CachedEnforcer) *RoleRouter {
	data := repo.NewRoleRepo(rdb)
	sv := service.NewRoleService(data, rdb, conf, logger)
	return &RoleRouter{
		api:    handler.NewRoleHandle(sv),
		public: routerGroup.Group("sysRole"),
		private: routerGroup.Group("sysRole").Use(
			middleware.JwtAuth(rdb),
			middleware.Casbin(enforcer),
		),
	}
}

func (r *RoleRouter) InitRouters() *RoleRouter {
	r.private.POST("", r.api.CreateRole)
	r.private.PUT("", r.api.UpdateRole)
	r.private.GET("list", r.api.GetRoleList)
	r.private.DELETE("", r.api.DeleteRoleById)
	r.private.DELETE("deleteRoleByIds", r.api.DeleteRoleByIds)
	return r
}

func (r *RoleRouter) InitRole() *RoleRouter {
	_ = r.api.InitRole()
	return r
}

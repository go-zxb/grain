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

type MenuRouter struct {
	public  gin.IRoutes
	private gin.IRoutes
	api     *handler.MenuHandle
}

func NewMenuRouter(routerGroup *gin.RouterGroup, rdb redis.IRedis, conf *config.Config, logger *log.Logger, enforcer *casbin.CachedEnforcer) *MenuRouter {
	data := repo.NewMenuRepo(rdb)
	sv := service.NewMenuService(data, rdb, conf, logger)
	return &MenuRouter{
		public: routerGroup.Group("sysMenu"),
		api:    handler.NewMenuHandle(sv),
		private: routerGroup.Group("").Use(
			middleware.JwtAuth(rdb),
			middleware.Casbin(enforcer),
		),
	}
}

func (r *MenuRouter) InitRouters() *MenuRouter {
	r.private.PUT("sysMenu", r.api.UpdateMenu)
	r.private.POST("sysMenu", r.api.CreateMenu)
	r.private.GET("sysMenu/list", r.api.GetMenuList)
	r.private.GET("sysMenu/userMenu", r.api.GetUserMenu)
	r.private.DELETE("sysMenu", r.api.DeleteMenuById)
	r.private.DELETE("sysMenu/deleteMenuByIds", r.api.DeleteMenuByIds)
	return r
}

func (r *MenuRouter) InitMenu() *MenuRouter {
	_ = r.api.InitMenu()
	return r
}

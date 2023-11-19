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
	"gorm.io/gorm"
)

type CodeFactoryRouter struct {
	public  gin.IRoutes
	private gin.IRoutes
	api     *handler.CodeAssistantHandle
}

func NewCodeAssistantRouter(routerGroup *gin.RouterGroup, db *gorm.DB, rdb redis.IRedis, conf *config.Config, logger *log.Logger, enforcer *casbin.CachedEnforcer) *CodeFactoryRouter {
	data := repo.NewCodeAssistantRepo(db, rdb)
	sv := service.NewCodeAssistantService(data, rdb, conf, logger)
	return &CodeFactoryRouter{
		public: routerGroup.Group("codeAssistant"),
		api:    handler.NewCodeAssistantHandle(sv),
		private: routerGroup.Group("codeAssistant").Use(
			middleware.JwtAuth(rdb),
			middleware.Casbin(enforcer),
		),
	}
}

func (r *CodeFactoryRouter) InitRouters() {
	r.private.POST("projects", r.api.CreateProject)
	r.private.PUT("projects", r.api.UpdateProject)
	r.private.GET("projects/list", r.api.GetProjectList)
	r.private.DELETE("projects", r.api.DeleteProjectById)

	r.private.POST("fields", r.api.CreateField)
	r.private.PUT("fields", r.api.UpdateField)
	r.private.DELETE("fields", r.api.DeleteFieldById)
	r.private.GET("fields/list", r.api.GetFieldList)

	r.private.POST("models", r.api.CreateModel)
	r.private.PUT("models", r.api.UpdateModel)
	r.private.DELETE("models", r.api.DeleteModelById)
	r.private.GET("models/list", r.api.GetModelList)

	r.private.GET("viewCode", r.api.ViewCode)
	r.private.POST("generateCode", r.api.GenerateCode)
}

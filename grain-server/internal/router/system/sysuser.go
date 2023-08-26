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

type SysUserRouter struct {
	// redis 实例对象
	rdb redis.IRedis
	// api handle
	api *handler.SysUserHandle
	// gin engine 干干净净的 engine
	engine *gin.Engine
	// 公开接口使用
	public gin.IRoutes
	// 需要登录才能使用的接口使用
	private gin.IRoutes
	// 需要鉴权的接口使用
	privateRoleAuth gin.IRoutes
}

func NewSysUserRouter(engine *gin.Engine, routerGroup *gin.RouterGroup, rdb redis.IRedis, conf *config.Config, enforcer *casbin.CachedEnforcer, logger *log.Logger) *SysUserRouter {
	data := repo.NewSysUserRepo(rdb)
	sv := service.NewSysUserService(data, rdb, conf, logger)
	return &SysUserRouter{
		rdb:    rdb,
		api:    handler.NewSysUserHandle(sv),
		engine: engine,
		public: routerGroup.Group("sysUser"),
		private: routerGroup.Group("sysUser").Use(
			middleware.JwtAuth(rdb)),
		privateRoleAuth: routerGroup.Group("sysUser").Use(
			middleware.JwtAuth(rdb),
			middleware.Casbin(enforcer),
		),
	}
}

func (r *SysUserRouter) InitRouters() {
	//登录接口
	r.public.POST("login", r.api.Login)
	//退出接口
	r.private.POST("logout", r.api.LogOut)
	//更新个人信息接口
	r.private.PUT("update", r.api.UpdateSysUser)
	//修改邮箱接口
	r.private.PUT("modifyEmail", r.api.ModifyEmail)
	//修改手机号接口
	r.private.PUT("modifyMobile", r.api.ModifyMobile)
	//修改手机号接口
	r.private.PUT("modifyPassword", r.api.ModifyPassword)
	//修改头像接口
	r.private.POST("modifyAvatar", r.api.UploadAvatar)
	//获取个人信息接口
	r.privateRoleAuth.GET("info", r.api.GetLoginUserInfo)
	//创建用户接口
	r.privateRoleAuth.POST("create", r.api.CreateSysUser)
	//根据用户id获取用户个人信息接口
	r.privateRoleAuth.GET("", r.api.GetSysUserById)
	//获取用户列表数据接口
	r.privateRoleAuth.GET("list", r.api.GetSysUserList)
	//切换角色接口
	r.private.POST("switchRole", middleware.SwitchRole(r.rdb), r.api.SwitchRole)
	//设置默认角色接口
	r.privateRoleAuth.PUT("setDefaultRole", r.api.SetDefaultRole)
	//编辑用户信息接口
	r.privateRoleAuth.PUT("editUserInfo", r.api.EditSysUser)
	// 确认修改邮箱接口
	r.engine.GET("confirmModifyEmail", r.api.ConfirmModifyEmail)
	//根据用户Id删除用户
	r.privateRoleAuth.DELETE("", r.api.DeleteSysUserById)
	//根据Id批量删除用户
	r.privateRoleAuth.DELETE("deleteSysUserByIdList", r.api.DeleteSysUserByIdList)
}

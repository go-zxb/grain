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

package service

import (
	"errors"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	"github.com/go-pay/gopay/pkg/xlog"
	"gorm.io/gorm"
)

const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

type ICasbinRepo interface {
	Update(roles []*model.CasbinRule) error
	AuthApiList(role string) ([]*model.CasbinRule, error)
}

type CasbinService struct {
	repo     ICasbinRepo
	enforcer *casbin.CachedEnforcer
	conf     *config.Config
	log      *log.Helper
}

func NewCasbinService(repo ICasbinRepo, conf *config.Config, logger log.Logger, e *casbin.CachedEnforcer) *CasbinService {
	return &CasbinService{repo: repo, enforcer: e, log: log.NewHelper(logger), conf: conf}
}

// InitCasbinRoleRule 初始化角色默认权限规则
func (s *CasbinService) InitCasbinRoleRule() error {
	defaultAdminRole := s.conf.System.DefaultAdminRole
	defaultRole := s.conf.System.DefaultRole
	casbinRule := []*model.CasbinRule{

		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/casbin", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/casbin/list", V2: "GET"},

		// 系统用户
		{Ptype: "p", V0: defaultRole, V1: "/api/v1/sysUser/info", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/info", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/update", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/avatar", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/create", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/editUserInfo", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/setDefaultRole", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysUser/deleteSysUserByIdList", V2: "DELETE"},

		// 系统角色
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysRole", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysRole", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysRole", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysRole/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysRole/deleteRoleList", V2: "DELETE"},

		// 系统API
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi/list", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi/apiGroups", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysApi/apiAndPermissions", V2: "GET"},

		//系统菜单
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysMenu", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysMenu", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysMenu", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysMenu/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysMenu/userMenu", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysMenu/deleteSysMenuList", V2: "DELETE"},

		//代码助手
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/fields", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/models", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/projects", V2: "POST"},

		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/models/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/fields/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/projects/list", V2: "GET"},

		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/models", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/fields", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/projects", V2: "PUT"},

		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/models", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/fields", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/projects", V2: "DELETE"},

		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/codeAssistant/viewCode", V2: "GET"},

		// 系统日志
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysLog", V2: "DELETE"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/sysLog/list", V2: "GET"},

		//组织管理
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/organize", V2: "PUT"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/organize", V2: "POST"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/organize/list", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/organize/listGroup", V2: "GET"},
		{Ptype: "p", V0: defaultAdminRole, V1: "/api/v1/organize/organizeById", V2: "DELETE"},
	}

	q := query.Q.CasbinRule

	count, err := q.Count()
	if err != nil {
		return err
	}

	// 有数据则退出,没有则添加默认数据
	if count != 0 {
		return nil
	}

	_, err = q.Where(q.ID.Gt(0)).Delete()
	if err != nil {
		return err
	}

	return q.Create(casbinRule...)
}

// NewCasbin 创建一个casbin实例对象
func NewCasbin(db *gorm.DB) *casbin.CachedEnforcer {
	// 使用 GORM 适配器创建 Casbin 的 enforcer 对象
	a, _ := gormadapter.NewAdapterByDB(db)
	newModelFromString, err := casbinModel.NewModelFromString(modelText)
	if err != nil {
		xlog.Info(err)
		return nil
	}

	enforcer, _ := casbin.NewCachedEnforcer(newModelFromString, a)

	// 将策略规则从数据库加载到 Casbin 中
	if err := enforcer.LoadPolicy(); err != nil {
		xlog.Info(err)
		return nil
	}

	return enforcer
}

// ReLoadPolicy 重新加载权限数据
func (s *CasbinService) ReLoadPolicy() error {
	// 将策略规则从数据库加载到 Casbin 中
	if err := s.enforcer.LoadPolicy(); err != nil {
		xlog.Info(err)
		return err
	}
	return nil
}

// RemoveFilteredPolicy 移除xx角色已分配的权限
func (s *CasbinService) RemoveFilteredPolicy(role string) error {
	_, err := s.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新
func (s *CasbinService) Update(roles *model.CasbinReq, ctx *gin.Context) error {

	if len(roles.Data) == 0 {
		return errors.New("至少要选择一个吧")
	}

	apis, err := query.Q.SysApi.Where(query.SysApi.ID.In(roles.Data...)).Find()
	if err != nil {
		return err
	}

	var c = make([]*model.CasbinRule, len(apis))
	for i, val := range apis {
		c[i] = &model.CasbinRule{}
		c[i].Ptype = "p"
		c[i].V0 = roles.Role
		c[i].V1 = val.Path
		c[i].V2 = val.Method
	}

	oldList, _ := s.AuthApiList(roles.Role)
	if err := s.RemoveFilteredPolicy(roles.Role); err != nil {
		return err
	}

	if err := s.repo.Update(c); err != nil {
		if err = s.repo.Update(oldList); err != nil {
			s.log.Errorw("errMsg", "更新角色权限失败", "err", err.Error())
			return errors.New("更新失败,完犊子了,我一点补救的办法都没有 我能怎么办 你说我能怎么办 ^*^*^")
		}
	}

	if err := s.ReLoadPolicy(); err != nil {
		return err
	}
	s.log.Infow("errMsg", "更新角色权限")
	return nil
}

// AuthApiList 获取已分配的角色资源列表
func (s *CasbinService) AuthApiList(role string) (list []*model.CasbinRule, err error) {
	list, err = s.repo.AuthApiList(role)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("无数据")
	}
	return list, nil
}

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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	redisx "github.com/go-grain/grain/pkg/redis"
)

type IMenuRepo interface {
	CreateMenu(menu *model.SysMenu) error
	CreateUserMenu(menu []*model.SysUserMenu) error
	GetMenuById(id uint) (*model.SysMenu, error)
	GetUserMenu(role string, parentId uint) (u []*model.SysMenu, err error)
	GetMenuList() (list []*model.SysMenu, err error)
	GetMenuListByParentId(req *model.SysMenuReq, parentId uint) ([]*model.SysMenu, error)
	GetUserMenuByRoleAndID(role string, pid uint) (list []*model.SysUserMenu, err error)
	GetUserMenuByRole(role string) (list []*model.SysUserMenu, err error)
	UpdateMenu(menu *model.SysMenu) error
	UpdateMenus(menu []*model.SysMenu) error
	DeleteMenuById(menuId uint) error
	DeleteMenuByIds(ids []uint) error
	DeleteUserMenuByRole(role string) error
}

type MenuService struct {
	repo IMenuRepo
	rdb  redisx.IRedis
	conf *config.Config
	log  *log.Helper
}

func NewMenuService(repo IMenuRepo, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *MenuService {
	return &MenuService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  log.NewHelper(logger),
	}
}

// InitMenu 默认菜单
func (s *MenuService) InitMenu() error {

	q := query.Q.SysMenu

	count, err := q.Count()
	if err != nil {
		return err
	}

	//  有数据则退出,没有则添加默认数据
	if count != 0 {
		return nil
	}

	menu := &model.SysMenu{
		ParentId: 0,
		Path:     "manage",
		Name:     "manage",
		CnName:   "菜单管理",
		Meta: &model.Meta{
			I18n:         "menu.sysManage",
			RequiresAuth: false,
			Icon:         "icon-command",
			Order:        2,
		},
	}

	if err = q.Create(menu); err != nil {
		return err
	}
	Children := []*model.SysMenu{
		{
			ParentId: menu.ID,
			Path:     "sysLog",
			Name:     "sysLog",
			CnName:   "系统日志",
			Meta: &model.Meta{
				I18n:         "menu.sysLog",
				RequiresAuth: false,
				Icon:         "",
				Order:        6,
				Roles:        []string{s.conf.System.DefaultAdminRole},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysMenu",
			Name:     "sysMenu",
			CnName:   "系统菜单",
			Meta: &model.Meta{
				I18n:         "menu.sysMenu",
				RequiresAuth: false,
				Icon:         "",
				Order:        5,
				Roles:        []string{s.conf.System.DefaultAdminRole},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysOrganize",
			Name:     "sysOrganize",
			CnName:   "组织管理",
			Meta: &model.Meta{
				I18n:         "menu.organize",
				RequiresAuth: true,
				Order:        4,
				Roles:        []string{s.conf.System.DefaultAdminRole},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysUser",
			Name:     "sysUser",
			CnName:   "系统用户",
			Meta: &model.Meta{
				I18n:         "menu.sysUser",
				RequiresAuth: false,
				Icon:         "",
				Order:        3,
				Roles:        []string{s.conf.System.DefaultAdminRole},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysApi",
			Name:     "sysApi",
			CnName:   "系统Api",
			Meta: &model.Meta{
				I18n:         "menu.sysApi",
				RequiresAuth: false,
				Icon:         "",
				Order:        2,
				Roles:        []string{s.conf.System.DefaultAdminRole},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysRole",
			Name:     "sysRole",
			CnName:   "用户角色",
			Meta: &model.Meta{
				I18n:         "menu.sysRole",
				RequiresAuth: false,
				Icon:         "",
				Order:        1,
				Roles:        []string{s.conf.System.DefaultAdminRole},
			},
		},
	}

	if err = q.Create(Children...); err != nil {
		return err
	}

	menu = &model.SysMenu{
		ParentId: 0,
		Path:     "dashboard",
		Name:     "dashboard",
		CnName:   "仪表盘",
		Meta: &model.Meta{
			I18n:         "menu.dashboard",
			RequiresAuth: false,
			Icon:         "icon-dashboard",
			Order:        0,
		},
	}

	if err = q.Create(menu); err != nil {
		return err
	}

	menu = &model.SysMenu{
		ParentId: menu.ID,
		Path:     "workplace",
		Name:     "Workplace",
		CnName:   "工作台",
		Meta: &model.Meta{
			I18n:         "menu.dashboard.workplace",
			RequiresAuth: false,
			Order:        0,
			Roles:        []string{s.conf.System.DefaultAdminRole},
		},
	}
	if err = q.Create(menu); err != nil {
		return err
	}

	menu = &model.SysMenu{
		ParentId: 0,
		Path:     "codeAssistant",
		Name:     "codeAssistant",
		CnName:   "代码助手",
		Meta: &model.Meta{
			I18n:         "menu.codeAssistant",
			RequiresAuth: false,
			Icon:         "icon-code",
			Order:        1,
		},
	}
	if err = q.Create(menu); err != nil {
		return err
	}
	menu = &model.SysMenu{
		ParentId: menu.ID,
		Path:     "generateCode",
		Name:     "generateCode",
		CnName:   "代码生成器",
		Meta: &model.Meta{
			I18n:         "menu.generateCode",
			RequiresAuth: true,
			Order:        0,
			Roles:        []string{s.conf.System.DefaultAdminRole},
		},
	}
	if err = q.Create(menu); err != nil {
		return err
	}

	menu = &model.SysMenu{
		ParentId: 0,
		Path:     "attachments",
		Name:     "attachments",
		CnName:   "附件管理",
		Meta: &model.Meta{
			I18n:         "menu.attachments",
			RequiresAuth: false,
			Icon:         "icon-upload",
			Order:        3,
		},
	}
	if err = q.Create(menu); err != nil {
		return err
	}
	menu = &model.SysMenu{
		ParentId: menu.ID,
		Path:     "sysFile",
		Name:     "sysFile",
		CnName:   "系统附件",
		Meta: &model.Meta{
			I18n:         "menu.sysFile",
			RequiresAuth: true,
			Order:        0,
			Roles:        []string{s.conf.System.DefaultAdminRole},
		},
	}
	if err = q.Create(menu); err != nil {
		return err
	}

	list, err := s.repo.GetMenuList()
	if err != nil {
		return err
	}

	var newList []*model.SysUserMenu
	for _, menu := range list {
		t := &model.SysUserMenu{
			MID:      menu.ID,
			ParentId: menu.ParentId,
			Role:     s.conf.System.DefaultAdminRole,
			CnName:   menu.CnName,
			Name:     menu.Name,
			Path:     menu.Path,
			Meta:     menu.Meta,
		}
		newList = append(newList, t)
	}

	if err := s.repo.CreateUserMenu(newList); err != nil {
		return err
	}
	return nil
}

func (s *MenuService) CreateMenu(menu *model.SysMenu, ctx *gin.Context) error {
	if err := s.repo.CreateMenu(menu); err != nil {
		s.log.Errorw("errMsg", "创建菜单", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "创建菜单")
	return nil
}

func (s *MenuService) GetUserMenu(role string, ctx *gin.Context) (menu []*model.SysMenu, err error) {
	menuAll, err := s.repo.GetUserMenu(role, 0)
	if err != nil {
		return nil, err
	}

	for _, m := range menuAll {
		t := &model.SysMenu{
			Model:    m.Model,
			ParentId: m.ParentId,
			CnName:   m.CnName,
			Name:     m.Name,
			Path:     m.Path,
			Meta:     m.Meta,
		}
		list, err := s.repo.GetUserMenuByRoleAndID(role, m.ID)
		if err != nil || len(list) == 0 {
			continue
		}
		var m2 []*model.SysMenu
		for _, userMenu := range list {
			m2 = append(m2, &model.SysMenu{
				Model:    userMenu.Model,
				ParentId: userMenu.ParentId,
				CnName:   userMenu.CnName,
				Name:     userMenu.Name,
				Path:     userMenu.Path,
				Meta:     userMenu.Meta,
			})
		}
		t.Children = m2
		menu = append(menu, t)
	}
	return menu, err
}

func (s *MenuService) GetMenuAndPermission(role string, ctx *gin.Context) (menu any, selectKeys []uint, err error) {
	req := &model.SysMenuReq{}
	menuAll, err := s.repo.GetMenuListByParentId(req, 0)
	if err != nil {
		return nil, nil, err
	}
	type Menu struct {
		Key      uint    `form:"key" json:"key"`
		Title    string  `form:"title" json:"title"`
		Children []*Menu `form:"children" json:"children"`
	}
	var menuList []*Menu
	for _, sysMenu := range menuAll {
		t := Menu{
			Key:      sysMenu.ID,
			Title:    sysMenu.CnName,
			Children: nil,
		}
		menuAll2, err := s.repo.GetMenuListByParentId(req, sysMenu.ID)
		if err != nil {
			continue
		}

		for _, m2 := range menuAll2 {
			t.Children = append(t.Children, &Menu{
				Key:      m2.ID,
				Title:    m2.CnName,
				Children: nil,
			})
		}
		menuList = append(menuList, &t)
	}

	byRole, err := s.repo.GetUserMenuByRole(role)
	if err != nil {
		return nil, nil, err
	}

	for _, userMenu := range byRole {
		selectKeys = append(selectKeys, userMenu.MID)
	}

	return menuList, selectKeys, nil
}

func (s *MenuService) GetMenuList(req *model.SysMenuReq, ctx *gin.Context) ([]*model.SysMenu, error) {
	list, err := s.repo.GetMenuListByParentId(req, 0)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}

	for i, menu := range list {
		child, _ := s.repo.GetMenuListByParentId(req, menu.ID)
		if len(child) > 0 {
			list[i].Children = child
			for j, child_ := range list[i].Children {
				child_child, _ := s.repo.GetMenuListByParentId(req, child_.ID)
				if len(child_child) > 0 {
					list[i].Children[j].Children = child_child
				}
			}
		}
	}
	return list, err
}

func (s *MenuService) SetMenuAndPermission(keys []uint, role string) error {
	fmt.Println(keys)
	if len(keys) == 0 {
		return errors.New("参数不能为空")
	}

	if err := s.repo.DeleteUserMenuByRole(role); err != nil {
		return err
	}

	list, err := s.repo.GetMenuList()
	if err != nil {
		return err
	}

	var newList []*model.SysUserMenu
	for _, key := range keys {
		//拿出选择的菜单
		for _, menu := range list {
			t := &model.SysUserMenu{
				MID:      menu.ID,
				ParentId: menu.ParentId,
				Role:     role,
				CnName:   menu.CnName,
				Name:     menu.Name,
				Path:     menu.Path,
				Meta:     menu.Meta,
			}
			if menu.ID == key {
				newList = append(newList, t)
			}
		}
	}
	return s.repo.CreateUserMenu(newList)
}

func (s *MenuService) UpdateMenu(menu *model.SysMenu, ctx *gin.Context) error {
	if menu.ParentId != 0 {
		m, err := s.repo.GetMenuById(menu.ParentId)
		if err != nil {
			return err
		}
		if m.ID == 0 {
			return errors.New("父ID不存在")
		}
	}

	if err := s.repo.UpdateMenu(menu); err != nil {
		s.log.Errorw("errMsg", "更新菜单", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "更新菜单")
	return nil
}

func (s *MenuService) DeleteMenuById(id uint, ctx *gin.Context) error {
	if err := s.repo.DeleteMenuById(id); err != nil {
		s.log.Errorw("errMsg", "删除菜单", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "删除菜单")
	return nil
}

func (s *MenuService) DeleteMenuByIds(ids []uint, ctx *gin.Context) error {
	if err := s.repo.DeleteMenuByIds(ids); err != nil {
		s.log.Errorw("errMsg", "批量删除菜单", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "批量删除菜单")
	return nil
}

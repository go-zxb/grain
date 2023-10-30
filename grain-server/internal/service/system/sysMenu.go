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
	"github.com/gin-gonic/gin"
	xjson "github.com/go-grain/go-utils/json"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	"gorm.io/gorm"
)

type IMenuRepo interface {
	CreateMenu(menu *model.SysMenu) error
	GetMenuById(id uint) (*model.SysMenu, error)
	GetUserMenu(role string, parentId uint) (u []*model.SysMenu, err error)
	GetMenuList(req *model.SysMenuReq, parentId uint) ([]*model.SysMenu, error)
	UpdateMenu(menu *model.SysMenu) error
	DeleteMenuById(menuId uint) error
	DeleteMenuByIds(ids []uint) error
}

type MenuService struct {
	repo IMenuRepo
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Logger
}

func NewMenuService(repo IMenuRepo, rdb redis.IRedis, conf *config.Config, logger *log.Logger) *MenuService {
	return &MenuService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  logger,
	}
}

// InitMenu 默认菜单
func InitMenu(db *gorm.DB) error {

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
		Meta: &model.Meta{
			I18n:         "menu.admin",
			RequiresAuth: false,
			Icon:         "icon-command",
			Order:        2,
		},
	}
	db.Create(menu)
	Children := []*model.SysMenu{
		{
			ParentId: menu.ID,
			Path:     "operationLog",
			Name:     "operationLog",
			Meta: &model.Meta{
				I18n:         "menu.operationLog",
				RequiresAuth: false,
				Icon:         "",
				Order:        5,
				Roles:        []string{"2023"},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysMenu",
			Name:     "sysMenu",
			Meta: &model.Meta{
				I18n:         "menu.sysMenu",
				RequiresAuth: false,
				Icon:         "",
				Order:        4,
				Roles:        []string{"2023"},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysUser",
			Name:     "sysUser",
			Meta: &model.Meta{
				I18n:         "menu.sysUser",
				RequiresAuth: false,
				Icon:         "",
				Order:        3,
				Roles:        []string{"2023"},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysApi",
			Name:     "sysApi",
			Meta: &model.Meta{
				I18n:         "menu.sysApi",
				RequiresAuth: false,
				Icon:         "",
				Order:        2,
				Roles:        []string{"2023"},
			},
		},
		{
			ParentId: menu.ID,
			Path:     "sysRole",
			Name:     "sysRole",
			Meta: &model.Meta{
				I18n:         "menu.sysRole",
				RequiresAuth: false,
				Icon:         "",
				Order:        1,
				Roles:        []string{"2023"},
			},
		},
	}
	db.Create(&Children)

	menu = &model.SysMenu{
		ParentId: 0,
		Path:     "dashboard",
		Name:     "dashboard",
		Meta: &model.Meta{
			I18n:         "menu.dashboard",
			RequiresAuth: false,
			Icon:         "icon-dashboard",
			Order:        0,
		},
	}
	db.Create(menu)
	menu = &model.SysMenu{
		ParentId: menu.ID,
		Path:     "workplace",
		Name:     "Workplace",
		Meta: &model.Meta{
			I18n:         "menu.dashboard.workplace",
			RequiresAuth: false,
			Order:        0,
			Roles:        []string{"2023"},
		},
	}
	db.Create(menu)

	menu = &model.SysMenu{
		ParentId: 0,
		Path:     "codeFactory",
		Name:     "codeFactory",
		Meta: &model.Meta{
			I18n:         "menu.codeFactory",
			RequiresAuth: false,
			Icon:         "icon-code",
			Order:        1,
		},
	}
	db.Create(menu)
	menu = &model.SysMenu{
		ParentId: menu.ID,
		Path:     "generateCode",
		Name:     "generateCode",
		Meta: &model.Meta{
			I18n:         "menu.codeFactory.codeFactory",
			RequiresAuth: true,
			Order:        0,
			Roles:        []string{"2023"},
		},
	}
	db.Create(menu)

	menu = &model.SysMenu{
		ParentId: 0,
		Path:     "uploads",
		Name:     "uploads",
		Meta: &model.Meta{
			I18n:         "menu.upload",
			RequiresAuth: false,
			Icon:         "icon-upload",
			Order:        3,
		},
	}
	db.Create(menu)
	menu = &model.SysMenu{
		ParentId: menu.ID,
		Path:     "upload",
		Name:     "upload",
		Meta: &model.Meta{
			I18n:         "menu.upload",
			RequiresAuth: true,
			Order:        0,
			Roles:        []string{"2023"},
		},
	}
	db.Create(menu)
	return nil
}

func (s *MenuService) CreateMenu(menu *model.SysMenu, ctx *gin.Context) error {
	err := s.repo.CreateMenu(menu)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "创建菜单", menu, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "创建菜单", menu, ctx))
	return nil
}

func (s *MenuService) GetUserMenu(role string, ctx *gin.Context) (menu []*model.SysMenu, err error) {
	menuAll, err := s.repo.GetUserMenu(role, 0)
	if err != nil {
		return nil, err
	}
	ParentAll := make(map[uint]*model.SysMenu)
	authNode := make(map[uint]*model.SysMenu)
	childNode := make(map[uint]*model.SysMenu)
	for i, m := range menuAll {
		if m.ParentId == 0 {
			ParentAll[m.ID] = menuAll[i]
		} else {
			for _, s2 := range m.Meta.Roles {
				if s2 == role {
					m.Meta.RequiresAuth = true
					m.Meta.Roles = []string{role}
					childNode[m.ID] = m
				}
			}
		}
	}

	for key, val := range childNode {
		_, ok := authNode[val.ParentId]
		if ok {
			authNode[val.ParentId].Children = append(authNode[val.ParentId].Children, childNode[key])
		} else {
			authNode[val.ParentId] = ParentAll[val.ParentId]
			authNode[val.ParentId].Children = append(authNode[val.ParentId].Children, childNode[key])
		}

	}

	for _, m := range authNode {
		menu = append(menu, m)
	}
	return menu, err
}

func (s *MenuService) GetMenuList(req *model.SysMenuReq, ctx *gin.Context) ([]*model.SysMenu, error) {
	list, err := s.repo.GetMenuList(req, 0)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}

	for i, menu := range list {
		child, _ := s.repo.GetMenuList(req, menu.ID)
		if len(child) > 0 {
			list[i].Children = child
			for j, child_ := range list[i].Children {
				child_child, _ := s.repo.GetMenuList(req, child_.ID)
				if len(child_child) > 0 {
					list[i].Children[j].Children = child_child
				}
			}
		}
	}
	return list, err
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

	err := s.repo.UpdateMenu(menu)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新菜单", menu, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新菜单", menu, ctx))
	return nil
}

func (s *MenuService) DeleteMenuById(id uint, ctx *gin.Context) error {
	err := s.repo.DeleteMenuById(id)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除菜单", xjson.G{"id": id}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除菜单", xjson.G{"id": id}, ctx))
	return nil
}

func (s *MenuService) DeleteMenuByIds(ids []uint, ctx *gin.Context) error {
	err := s.repo.DeleteMenuByIds(ids)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除菜单", xjson.G{"id": ids}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除菜单", xjson.G{"id": ids}, ctx))
	return nil
}

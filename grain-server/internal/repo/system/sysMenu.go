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

package repo

import (
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/internal/repo/system/query"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/model/system"
)

type MenuRepo struct {
	rdb   redis.IRedis
	query *query.Query
}

func NewMenuRepo(rdb redis.IRedis) service.IMenuRepo {
	return &MenuRepo{
		rdb:   rdb,
		query: query.Q,
	}
}

func (r *MenuRepo) CreateMenu(menu *model.SysMenu) error {
	return r.query.SysMenu.Create(menu)
}

func (r *MenuRepo) GetMenuById(parentId uint) (menu *model.SysMenu, err error) {
	if menu, err = r.query.SysMenu.Where(r.query.SysMenu.ID.Eq(parentId)).First(); err != nil {
		return nil, err
	}
	return
}

func (r *MenuRepo) GetUserMenu(role string, parentId uint) (list []*model.SysMenu, err error) {
	if list, err = r.query.SysMenu.Where(r.query.SysMenu.ParentId.Eq(0)).Find(); err != nil {
		return nil, err
	}
	return
}

func (r *MenuRepo) GetMenuList() (list []*model.SysMenu, err error) {
	if list, err = r.query.SysMenu.Where(r.query.SysMenu.ParentId.Neq(0)).Find(); err != nil {
		return nil, err
	}
	return
}

func (r *MenuRepo) GetMenuListByParentId(req *model.SysMenuReq, parentId uint) (list []*model.SysMenu, err error) {
	count, err := r.query.SysMenu.Count()
	if err != nil {
		return nil, err
	}
	req.Total = count
	if list, err = r.query.SysMenu.Where(r.query.SysMenu.ParentId.Eq(parentId)).Find(); err != nil {
		return nil, err
	}

	return
}

func (r *MenuRepo) UpdateMenus(menu []*model.SysMenu) error {
	if _, err := r.query.SysMenu.Updates(&menu); err != nil {
		return err
	}
	return nil
}

func (r *MenuRepo) UpdateMenu(menu *model.SysMenu) error {
	if _, err := r.query.SysMenu.Updates(menu); err != nil {
		return err
	}
	return nil
}

func (r *MenuRepo) DeleteMenuById(menuId uint) error {
	if _, err := r.query.SysMenu.Where(r.query.SysMenu.ID.Eq(menuId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *MenuRepo) DeleteMenuByIds(ids []uint) error {
	if _, err := r.query.SysMenu.Where(r.query.SysMenu.ID.In(ids...)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *MenuRepo) GetUserMenuByRole(role string) (list []*model.SysUserMenu, err error) {
	if list, err = r.query.SysUserMenu.Where(r.query.SysUserMenu.Role.Eq(role)).Find(); err != nil {
		return
	}
	return
}

func (r *MenuRepo) GetUserMenuByRoleAndID(role string, pid uint) (list []*model.SysUserMenu, err error) {
	if list, err = r.query.SysUserMenu.Where(r.query.SysUserMenu.Role.Eq(role)).Where(r.query.SysUserMenu.ParentId.Eq(pid)).Find(); err != nil {
		return
	}
	return
}

func (r *MenuRepo) DeleteUserMenuByRole(role string) error {
	if _, err := r.query.SysUserMenu.Where(r.query.SysUserMenu.Role.Eq(role)).Unscoped().Delete(); err != nil {
		return err
	}
	return nil
}

func (r *MenuRepo) CreateUserMenu(menu []*model.SysUserMenu) error {
	return r.query.SysUserMenu.Create(menu...)
}

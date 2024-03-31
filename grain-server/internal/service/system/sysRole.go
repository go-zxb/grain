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
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/utils/const"
	"strings"
)

type IRoleRepo interface {
	CreateRole(user *model.SysRole) error
	GetRoleList(req *model.SysRoleQueryPage) ([]*model.SysRole, error)
	UpdateRole(user *model.SysRole) error
	DeleteRoleById(roleId uint) error
	DeleteRoleByIds(userIds []uint) error
}

type RoleService struct {
	repo IRoleRepo
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Helper
}

func NewRoleService(repo IRoleRepo, rdb redis.IRedis, conf *config.Config, logger log.Logger) *RoleService {
	return &RoleService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  log.NewHelper(logger),
	}
}

func (s *RoleService) InitRole() error {
	roles := []*model.SysRole{
		{
			Model:    model.Model{},
			Role:     "2023",
			RoleName: "超级管理员",
		}, {
			Model:    model.Model{},
			Role:     "2024",
			RoleName: "普通成员",
		},
	}
	q := query.Q.SysRole
	count, err := q.Count()
	if err != nil {
		return err
	}

	// 有数据则退出 否则就添加
	if count != 0 {
		return nil
	}

	return q.Create(roles...)
}

func (s *RoleService) CreateRole(role *model.CreateSysRole, ctx *gin.Context) error {
	_role := model.SysRole{
		Role:     role.Role,
		RoleName: role.RoleName,
	}

	if err := s.repo.CreateRole(&_role); err != nil {
		s.log.Errorw("errMsg", "批量删除菜单", "err", err.Error())
		if strings.Contains(err.Error(), "duplicated key not allowed") {
			return errors.New("提交的参数重复")
		}
		return err
	}
	s.log.Errorw("errMsg", "批量删除菜单")
	return nil
}

func (s *RoleService) GetRoleList(req *model.SysRoleQueryPage, ctx *gin.Context) ([]*model.SysRole, error) {
	list, err := s.repo.GetRoleList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New(consts.ErrMsg(consts.NotRoleList))
	}
	return list, nil
}

func (s *RoleService) UpdateRole(role *model.SysRole, ctx *gin.Context) error {
	if err := s.repo.UpdateRole(role); err != nil {
		s.log.Errorw("errMsg", "更新角色", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "更新角色")
	return nil
}

func (s *RoleService) DeleteRoleByIds(roles []uint, ctx *gin.Context) error {
	if err := s.repo.DeleteRoleByIds(roles); err != nil {
		s.log.Errorw("errMsg", "删除角色", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "删除角色")
	return nil
}

func (s *RoleService) DeleteRoleById(roleId uint, ctx *gin.Context) error {
	if err := s.repo.DeleteRoleById(roleId); err != nil {
		s.log.Errorw("errMsg", "删除角色", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "删除角色")
	return nil
}

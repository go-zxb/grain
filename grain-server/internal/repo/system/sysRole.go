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

type RoleRepo struct {
	rdb   redis.IRedis
	query *query.Query
}

func NewRoleRepo(rdb redis.IRedis) service.IRoleRepo {
	return &RoleRepo{
		rdb:   rdb,
		query: query.Q,
	}
}

func (r *RoleRepo) CreateRole(role *model.SysRole) error {
	return r.query.SysRole.Create(role)
}

func (r *RoleRepo) GetRoleList(req *model.SysRoleQueryPage) (list []*model.SysRole, err error) {
	if list, err = r.query.SysRole.Find(); err != nil {
		return nil, err
	}
	return
}

func (r *RoleRepo) UpdateRole(role *model.SysRole) error {
	if _, err := r.query.SysRole.Updates(role); err != nil {
		return err
	}
	return nil
}

func (r *RoleRepo) DeleteRoleById(roleId uint) error {
	if _, err := r.query.SysRole.Where(r.query.SysRole.ID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *RoleRepo) DeleteRoleByIds(roles []uint) error {
	if _, err := r.query.SysRole.Where(r.query.SysRole.ID.In(roles...)).Delete(); err != nil {
		return err
	}
	return nil
}

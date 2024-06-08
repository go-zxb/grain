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
	"github.com/go-grain/grain/internal/repo/system/query"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/model/system"
	redisx "github.com/go-grain/grain/pkg/redis"
)

type SysUserRepo struct {
	rdb   redisx.IRedis
	query *query.Query
}

func NewSysUserRepo(rdb redisx.IRedis) service.ISysUserRepo {
	return &SysUserRepo{
		rdb:   rdb,
		query: query.Q,
	}
}

func (r *SysUserRepo) Login(user *model.LoginReq) (*model.SysUser, error) {
	userinfo, err := r.query.SysUser.Where(r.query.SysUser.Username.Eq(user.Username)).First()
	if err != nil {
		return nil, err
	}
	return userinfo, err
}

func (r *SysUserRepo) CreateSysUser(sysUser *model.SysUser) error {
	return r.query.SysUser.Create(sysUser)
}

func (r *SysUserRepo) GetSysUserById(sysUserId uint) (*model.SysUser, error) {
	return r.query.SysUser.Where(r.query.SysUser.ID.Eq(sysUserId)).First()
}

func (r *SysUserRepo) GetSysUserByUId(uid string) (*model.SysUser, error) {
	return r.query.SysUser.Where(r.query.SysUser.UID.Eq(uid)).First()
}

func (r *SysUserRepo) GetSysUserList(req *model.SysUserReq) (list []*model.SysUser, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 || req.PageSize >= 100 {
		req.PageSize = 20
	}

	q := r.query.SysUser.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	if req.Mobile != "" {
		q.Where(r.query.SysUser.Mobile.Eq(req.Mobile))
	}
	if req.Email != "" {
		q.Where(r.query.SysUser.Email.Eq(req.Email))
	}
	if req.Username != "" {
		q.Where(r.query.SysUser.Username.Eq(req.Username))
	}
	if req.Organize != "" {
		q.Where(r.query.SysUser.Organize.Eq(req.Organize))
	}
	if req.Department != "" {
		q.Where(r.query.SysUser.Department.Eq(req.Department))
	}
	if req.Position != "" {
		q.Where(r.query.SysUser.Position.Eq(req.Position))
	}
	count, err := r.query.SysUser.Count()
	if err != nil {
		return nil, err
	}
	req.Total = count

	list, err = q.Find()
	if err != nil {
		return nil, err
	}
	return
}

func (r *SysUserRepo) UpdateSysUser(sysUser *model.UpdateUserInfo) error {
	q := r.query.SysUser
	if _, err := q.Where(q.UID.Eq(sysUser.UID)).Updates(sysUser); err != nil {
		return err
	}
	return nil
}

func (r *SysUserRepo) EditSysUser(sysUser *model.SysUser) error {
	if _, err := r.query.SysUser.Updates(sysUser); err != nil {
		return err
	}
	return nil
}

func (r *SysUserRepo) SetDefaultRole(sysUser *model.SysUser) error {
	if _, err := r.query.SysUser.Updates(sysUser); err != nil {
		return err
	}
	return nil
}

func (r *SysUserRepo) DeleteSysUserById(sysUserId uint) error {
	if _, err := r.query.SysUser.Where(r.query.SysUser.ID.Eq(sysUserId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *SysUserRepo) DeleteSysUserByIds(sysUserIds []uint) error {
	if _, err := r.query.SysUser.Where(r.query.SysUser.ID.In(sysUserIds...)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *SysUserRepo) UploadAvatar(avatar *model.Upload, uid string) error {
	q := r.query.SysUser
	_, err := q.Where(q.UID.Eq(uid)).Update(q.Avatar, avatar.FileUrl)
	if err != nil {
		return err
	}
	_ = r.query.Upload.Create(avatar)
	return nil
}

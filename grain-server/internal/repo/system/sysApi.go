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

type ApiRepo struct {
	rdb    redis.IRedis
	casbin service.ICasbinRepo
	query  *query.Query
}

func NewApiRepo(rdb redis.IRedis) service.IApiRepo {
	return &ApiRepo{
		rdb:    rdb,
		casbin: NewCasbinRepo(),
		query:  query.Q,
	}
}

func (r *ApiRepo) CreateApi(api *model.SysApi) error {
	return r.query.SysApi.Create(api)
}

func (r *ApiRepo) GetApiList(req *model.SysApiReq) (list []*model.SysApi, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 || req.PageSize >= 100 {
		req.PageSize = 20
	}

	q := r.query.SysApi.Where()

	if req.Path != "" {
		q = q.Where(r.query.SysApi.ApiGroup.Eq(req.Path))
	}
	if req.Group != "" {
		q = q.Where(r.query.SysApi.ApiGroup.Eq(req.Group))
	}
	if req.Method != "" {
		q = q.Where(r.query.SysApi.Method.Eq(req.Method))
	}

	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	req.Total = count
	q = q.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)
	list, err = q.Find()
	if err != nil {
		return nil, err
	}

	return
}

func (r *ApiRepo) GetAllApi() (list []*model.SysApi, err error) {
	if list, err = r.query.SysApi.Find(); err != nil {
		return nil, err
	}
	return
}

func (r *ApiRepo) UpdateApi(api *model.SysApi) error {
	if _, err := r.query.SysApi.Updates(api); err != nil {
		return err
	}
	return nil
}

func (r *ApiRepo) DeleteApiById(id uint) error {
	if _, err := r.query.SysApi.Where(r.query.SysApi.ID.Eq(id)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *ApiRepo) DeleteApiByIds(ids []uint) error {
	if _, err := r.query.SysApi.Where(r.query.SysApi.ID.In(ids...)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *ApiRepo) AuthApiList(role string) (list []*model.CasbinRule, err error) {
	return r.casbin.AuthApiList(role)
}

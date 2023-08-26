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
	"fmt"
	utils "github.com/go-grain/go-utils"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/internal/repo/system/query"
	service "github.com/go-grain/grain/internal/service/system"
	model "github.com/go-grain/grain/model/system"
	"strings"
)

type UploadRepo struct {
	rdb   redis.IRedis
	query *query.Query
}

func NewUploadRepo(rdb redis.IRedis) service.IUploadRepo {
	return &UploadRepo{
		rdb:   rdb,
		query: query.Q,
	}
}

func (r *UploadRepo) CreateUpload(upload *model.Upload) error {
	return r.query.Upload.Create(upload)
}

func (r *UploadRepo) GetUploadList(req *model.UploadReq) (list []*model.Upload, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 || req.PageSize >= 100 {
		req.PageSize = 20
	}

	q := r.query.Upload.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)

	if req.QueryTime != "" {
		t := strings.Split(req.QueryTime, ",")
		if len(t) == 2 {
			s := utils.GetStringToDate(t[0], utils.YMD)
			e := utils.GetStringToDate(t[1], utils.YMD)
			q.Where(r.query.Upload.CreatedAt.Between(s, e))
		}
	}

	if req.FileName != "" {
		q.Where(r.query.Upload.FileName.Like(fmt.Sprintf("%s%s%s", "%", req.FileName, "%")))
	}

	count, err := r.query.Upload.Count()
	if err != nil {
		return nil, err
	}
	req.Total = count

	list, err = q.Order(query.Upload.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *UploadRepo) DeleteUploadById(id uint, uid string) error {
	if _, err := r.query.Upload.Where(r.query.Upload.UID.Eq(uid)).Where(r.query.Upload.ID.Eq(id)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *UploadRepo) DeleteUploadByIds(ids []uint, uid string) error {
	if _, err := r.query.Upload.Where(r.query.Upload.UID.Eq(uid)).Where(r.query.Upload.ID.In(ids...)).Delete(); err != nil {
		return err
	}
	return nil
}

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
	"github.com/go-grain/grain/internal/repo/data"
	"github.com/go-grain/grain/internal/repo/system/query"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/model/system"
	"gorm.io/gorm"
)

type CodeAssistantRepo struct {
	db      *data.DB
	rdb     redis.IRedis
	ApiRepo service.IApiRepo
	query   *query.Query
}

func NewCodeAssistantRepo(db *gorm.DB, rdb redis.IRedis) service.ICodeAssistantRepo {
	return &CodeAssistantRepo{
		db:      &data.DB{DB: db},
		rdb:     rdb,
		ApiRepo: NewApiRepo(rdb),
		query:   query.Q,
	}
}

func (r *CodeAssistantRepo) CreateProject(p *model.Project) error {
	return r.query.Project.Create(p)
}

func (r *CodeAssistantRepo) UpdateProject(p *model.Project) error {
	_, err := r.query.Project.Where(r.query.Project.ID.Eq(p.ID)).Updates(p)
	if err != nil {
		return err
	}
	return nil
}

func (r *CodeAssistantRepo) CreateModel(m *model.Models) error {
	return r.query.Models.Save(m)
}

func (r *CodeAssistantRepo) UpdateModel(m *model.Models) error {
	_, err := r.query.Models.Where(r.query.Models.ID.Eq(m.ID)).Updates(m)
	if err != nil {
		return err
	}
	return nil
}

func (r *CodeAssistantRepo) CreateField(f *model.Fields) error {
	return r.query.Fields.Save(f)
}

func (r *CodeAssistantRepo) UpdateField(f *model.Fields) error {
	_, err := r.query.Fields.Where(r.query.Fields.ID.Eq(f.ID)).Updates(f)
	if err != nil {
		return err
	}
	return nil
}

func (r *CodeAssistantRepo) DeleteProjectById(pid uint) error {
	q := r.query.Project
	if _, err := q.Where(q.ID.Eq(pid)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *CodeAssistantRepo) DeleteModelById(mid uint) error {
	q := r.query.Models
	if _, err := q.Where(q.ID.Eq(mid)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *CodeAssistantRepo) DeleteFieldById(fid uint) error {
	q := r.query.Fields
	if _, err := q.Where(q.ID.Eq(fid)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *CodeAssistantRepo) GetProject(pid uint) (list *model.Project, err error) {
	q := r.query.Project
	if list, err = q.Where(q.ID.Eq(pid)).First(); err != nil {
		return nil, err
	}
	return
}

func (r *CodeAssistantRepo) GetProjectList() (list []*model.Project, err error) {
	if list, err = r.query.Project.Find(); err != nil {
		return nil, err
	}
	return
}

// GetModel 根据模块ID获取模块
func (r *CodeAssistantRepo) GetModel(mid uint) (list *model.Models, err error) {
	q := r.query.Models
	list, err = r.query.Models.Where(q.ID.Eq(mid)).First()
	if err != nil {
		return nil, err
	}
	return
}

// GetModels 根据父ID(项目ID)获取所有模块
func (r *CodeAssistantRepo) GetModels(parentId uint) (list []*model.Models, err error) {
	q := r.query.Models
	list, err = q.Where(q.ParentId.Eq(parentId)).Find()
	if err != nil {
		return nil, err
	}
	return
}

func (r *CodeAssistantRepo) GetFields(parentId uint) (list []*model.Fields, err error) {
	q := r.query.Fields
	list, err = q.Where(q.ParentId.Eq(parentId)).Find()
	if err != nil {
		return nil, err
	}
	return
}

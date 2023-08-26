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
)

type CasbinRepo struct {
	query *query.Query
}

func NewCasbinRepo() service.ICasbinRepo {
	return &CasbinRepo{
		query: query.Q,
	}
}

func (r *CasbinRepo) Update(roles []*model.CasbinRule) error {
	if err := r.query.CasbinRule.Save(roles...); err != nil {
		return err
	}
	return nil
}

func (r *CasbinRepo) AuthApiList(role string) (list []*model.CasbinRule, err error) {
	if list, err = r.query.CasbinRule.Where(r.query.CasbinRule.V0.Eq(role)).Find(); err != nil {
		return nil, err
	}
	return list, err
}

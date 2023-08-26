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

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"time"
)

// ErrorRes 返回响应错误, 实际用不到这个,只是作为生成swag文档使用
type ErrorRes struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

type MySqlModel struct {
	// 自增ID
	ID uint `json:"id" xml:"id" gorm:"primarykey"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt" xml:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" xml:"updatedAt"`
	// 删除时间
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" xml:"deletedAt" gorm:"index"`
}

type MongoModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt" xml:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" xml:"updatedAt"`
	// 删除时间
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" xml:"deletedAt"`
}

// BeforeCreate 钩子函数： 创建前Gorm会调用
func (m *MySqlModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeSave 钩子函数： 保存前Gorm会调用
func (m *MySqlModel) BeforeSave(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate 钩子函数： 更新前Gorm会调用
func (m *MySqlModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

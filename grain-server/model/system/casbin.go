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

// CasbinRule casbin 结构体
type CasbinRule struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	// 类型 p 或 g
	Ptype string `gorm:"size:100" json:"ptype,omitempty"`
	// 角色
	V0 string `gorm:"size:100" json:"v0,omitempty"`
	// 资源
	V1 string `gorm:"size:100" json:"v1,omitempty"`
	// 方法
	V2 string `gorm:"size:100" json:"v2,omitempty"`
	V3 string `gorm:"size:100" json:"v3,omitempty"`
	V4 string `gorm:"size:100" json:"v4,omitempty"`
	V5 string `gorm:"size:100" json:"v5,omitempty"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}

// CasbinReq 用于返回xx角色能操作的的所有资源
type CasbinReq struct {
	Role string `json:"role"`
	Data []uint `json:"data"`
}

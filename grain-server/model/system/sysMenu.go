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
// 菜单的数据模型是根据ArcoPro所需的字段定义的...

package model

import (
	"database/sql/driver"
	"encoding/json"
)

// Meta 菜单数据
type Meta struct {
	// I18n 用来存储该条菜单的国际化信息
	I18n string `form:"i18n" json:"i18n" xml:"i18n" gorm:"comment:国际化"`
	// 是否需要权限才能访问
	RequiresAuth bool `form:"requiresAuth" json:"requiresAuth" xml:"requiresAuth" gorm:"default:true;comment:是否需要授权访问"`
	// icon图标
	Icon string `form:"icon" json:"icon" xml:"icon" gorm:"comment:图标"`
	// 排序
	Order uint `form:"order" json:"order" xml:"order" gorm:"comment:排序"`
	// 那些角色可以访问
	Roles []string `json:"roles,omitempty"`
}

// SysMenu 用来管理动态菜单的结构体
type SysMenu struct {
	Model
	ParentId uint       ` form:"parentId" json:"parentId" xml:"parentId"  gorm:"comment:父ID"`
	Path     string     `form:"path" json:"path" xml:"path" gorm:"comment:路径"`
	Name     string     `form:"name" json:"name" xml:"name" gorm:"comment:名称"`
	Meta     *Meta      `form:"meta" json:"meta" xml:"meta"  gorm:"type:json;comment:"`
	Children []*SysMenu `form:"children" json:"children" xml:"children" gorm:"-"`
}

func (SysMenu) TableName() string {
	return "sys_menus"
}

type SysMenuReq struct {
	PageReq
}

// Value 实现gorm value, scan接口,对Meta解析支持
func (i *Meta) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}

func (i *Meta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), i)
}

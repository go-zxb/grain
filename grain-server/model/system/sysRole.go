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

// SysRole 角色结构体
type SysRole struct {
	Model
	Role     string `form:"role" json:"role" xml:"role" gorm:"unique;not null;comment:角色ID" binding:"required"`
	RoleName string `form:"roleName" json:"roleName" xml:"roleName" gorm:"unique;not null;comment:角色名称" binding:"required"`
}

func (SysRole) TableName() string {
	return "sys_roles"
}

type CreateSysRole struct {
	Role     string `json:"role" binding:"required"`
	RoleName string `json:"roleName" binding:"required"`
}

type SysRoleQueryPage struct {
	PageReq
	Role     string `json:"role" form:"role"`
	RoleName string `json:"RoleName" form:"RoleName"`
}

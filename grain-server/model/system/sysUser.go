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
	"database/sql/driver"
	"encoding/json"
)

type RoleStr struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Roles []string

// SysUser 用户结构体
type SysUser struct {
	//基础字段
	Model
	//用户uid
	UID string `form:"uid" json:"uid" xml:"uid"  gorm:"unique;not null;comment:用户唯一标识符"`
	// 用户名
	Username string `form:"username" json:"username" xml:"username" gorm:"unique;not null;comment:用户名称"`
	// 用户密码
	Password string `json:"password,omitempty" xml:"password"  gorm:"comment:密码"`
	// 昵称
	Nickname string `form:"nickname" json:"nickname" xml:"nickname" gorm:"comment:用户昵称"`
	//email
	Email string `form:"email" json:"email" xml:"email" gorm:"comment:邮箱"`
	//手机号
	Mobile string `form:"mobile" json:"mobile" xml:"mobile" gorm:"comment:手机号"`
	//头像
	Avatar string `json:"avatar"  gorm:"comment:头像"`
	// 个人介绍
	Introduction string `json:"introduction"  gorm:"comment:介绍"`
	// 个人博客网站
	PersonalWebsite string `json:"personalWebsite"  gorm:"comment:个人网站"`
	//认证标志
	Certification int `json:"certification"  gorm:"comment:认证"`
	// 多角色字段
	Roles *Roles `form:"roles" json:"roles" xml:"roles" gorm:"comment:多角色"`
	//现有的默认角色
	Role string `form:"role" json:"role" xml:"role"  gorm:"comment:用户角色"`
	//账号状态,一般 yes 正常,no 异常封禁
	Status string `form:"status" json:"status" xml:"status" gorm:"comment:账号状态;default:no"`
	// 存放 角色ID,角色名称等角色信息,方便前端展示使用
	RoleStr []RoleStr `json:"roleStr,omitempty" gorm:"-"`
	//组织
	Organize string `form:"organize" json:"organize"`
	//部门
	Department string `form:"department" json:"department"`
	//职位
	Position string `form:"position" json:"position"`
}

// CreateSysUser 创建用户时使用这个结构体接收前端提交的数据,
// 使用 CreateSysUser 而不使用 SysUser 结构体是为了避免前端传入了,不该传入的参数
type CreateSysUser struct {
	Username   string `form:"username" json:"username" xml:"username" gorm:"unique;not null;comment:用户名称"`
	Password   string `json:"password,omitempty" xml:"password"  gorm:"comment:密码"`
	Nickname   string `form:"nickname" json:"nickname" xml:"nickname" gorm:"comment:用户昵称"`
	Email      string `form:"email" json:"email" xml:"email" gorm:"comment:邮箱"`
	Mobile     string `form:"mobile" json:"mobile" xml:"mobile" gorm:"comment:手机号"`
	Roles      *Roles `form:"roles" json:"roles" xml:"roles" gorm:"comment:多角色"`
	Role       string `form:"role" json:"role" xml:"role"  gorm:"comment:用户角色"`
	Organize   string `form:"organize" json:"organize"`
	Department string `form:"department" json:"department"`
	Position   string `form:"position" json:"position"`
	Status     string `form:"status" json:"status" xml:"status" gorm:"comment:账号状态;default:no"`
}

type DefaultRole struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

type LoginReq struct {
	Captcha  string `form:"captcha" json:"captcha"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginRes struct {
	ErrorRes
	Data struct {
		Token  string `form:"token" json:"token"`
		Expire int64  `form:"expire" json:"expire"`
	} `json:"data"`
}

// UpdateUserInfo 使用 UpdateUserInfo 结构体也是为了避免前端传入了,不该传入的参数,
// 对可以随便修改无关紧要的字段 使用这个就行;
// 敏感数据或隐私数据,需要通过相应验证后才能修改,一般会独立写接口单独更新
type UpdateUserInfo struct {
	UID             string `form:"uid" json:"uid" xml:"uid"  gorm:"unique;not null;comment:用户唯一标识符"`
	Nickname        string `form:"nickname" json:"nickname" xml:"nickname" gorm:"comment:用户昵称"`
	Avatar          string `json:"avatar"  gorm:"comment:头像"`
	Introduction    string `json:"introduction"  gorm:"comment:介绍"`
	PersonalWebsite string `json:"personalWebsite"  gorm:"comment:个人网站"`
}

// ModifyPassword 修改密码使用的结构体
type ModifyPassword struct {
	UID         string `json:"-"`
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ModifyEmail 修改邮箱使用的结构体
type ModifyEmail struct {
	UID     string `json:"-"`
	Captcha string `json:"captcha"`
	Email   string `json:"email" binding:"required,email"`
}

// ModifyMobile 修改手机号使用的结构体
type ModifyMobile struct {
	UID     string `json:"-"`
	Captcha string `json:"captcha" binding:"required"`
	Mobile  string `json:"mobile" binding:"required"`
}

type SysUserReq struct {
	PageReq
	Username   string `json:"username" xml:"username" form:"username"`
	Email      string `json:"email" xml:"email" form:"email" `
	Mobile     string `form:"mobile" json:"mobile"  form:"mobile"`
	Organize   string `form:"organize" json:"organize"`
	Department string `form:"department" json:"department"`
	Position   string `form:"position" json:"position"`
	Status     string `form:"status" json:"status"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

// Value 实现gorm value, scan接口,对roles解析支持
func (i *Roles) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}

func (i *Roles) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), i)
}

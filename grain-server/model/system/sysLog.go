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

type SysLog struct {
	// ID
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// 用户UID
	UID string `json:"uid" xml:"uid" bson:"uid"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	// 删除时间
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" bson:"deleted_at" gorm:"index"`
	// 名称
	Name string `form:"name" json:"name" xml:"name" gorm:"comment:名称"`
	// 角色
	Role string `form:"role" json:"role" xml:"role" gorm:"comment:角色"`
	// 用户名称
	Username string `form:"username" json:"username" xml:"username" gorm:"column:username;comment:用户名"`
	// 昵称
	Nickname string `form:"nickname" json:"nickname" xml:"nickname" gorm:"comment:用户名称"`
	// 请求方法
	Method string `form:"method" json:"method" xml:"method" gorm:"comment:请求方法"`
	// 请求路径
	Path string `form:"path" json:"path" xml:"path" gorm:"comment:请求路径"`
	// 请求数据
	ReqData any `form:"reqData" json:"reqData" xml:"reqData" gorm:"comment:请求数据"`
	// 响应code
	ResCode int `form:"resCode" json:"resCode" xml:"resCode" gorm:"comment:返回Code"`
	// 响应数据
	ResData any `form:"resData" json:"resData" xml:"resData" gorm:"comment:返回数据"`
	// 客户端请求IP
	ClientIP string `form:"clientIP" json:"clientIP" xml:"clientIP" gorm:"comment:客户端请求IP"`
	// 请求时间
	RequestAt time.Time `form:"requestAt" json:"requestAt" xml:"requestAt" gorm:"comment:请求时间"`
	// 响应时间
	ResponseAt time.Time `form:"responseAt" json:"responseAt" xml:"responseAt" gorm:"comment:响应时间"`
	// 延时
	Latency int64 `form:"latency" json:"latency" xml:"latency" gorm:"comment:延迟"`
	// 错误信息
	ErrorMessage string `form:"errorMessage" json:"errorMessage" xml:"errorMessage" gorm:"comment:错误Msg"`
	// 状态码
	StatusCode int `form:"statusCode" json:"statusCode" xml:"statusCode" gorm:"comment:状态码"`
	// 数据大小
	BodySize int `form:"bodySize" json:"bodySize" xml:"bodySize" gorm:"comment:bodySize"`
	// 日志类型 目前主要区分登录日志
	LogType string `form:"logType" json:"logType"`
}

// BeforeSave 钩子函数：在保存文档之前执行
func (m *SysLog) BeforeSave() {
	currentTime := time.Now()
	m.CreatedAt = currentTime
}

// BeforeUpdate 钩子函数：在更新文档之前执行
func (m *SysLog) BeforeUpdate() {
	currentTime := time.Now()
	m.UpdatedAt = currentTime
}

// SysLogReq 一般用来查询操作日志的数据
type SysLogReq struct {
	PageReq
	Name      string `form:"name" json:"name" xml:"name" gorm:"comment:名称"`
	Role      string `form:"role" json:"role" xml:"role" gorm:"comment:角色"`
	Username  string `form:"username" json:"username" xml:"username" gorm:"column:username;comment:用户名"`
	QueryTime string ` form:"queryTime" json:"queryTime"`
}

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

type SysApi struct {
	Model
	Path        string `json:"path" gorm:"uniqueIndex:api_path_method_idx;index:idx_path;comment:路径"`
	Method      string `json:"method" gorm:"uniqueIndex:api_path_method_idx;index:idx_method;comment:方法"`
	Description string `json:"description"  gorm:"comment:描述"`
	ApiGroup    string `json:"group" gorm:"comment:分组"`
}

type SysApiReq struct {
	PageReq
	Path        string `json:"path"  form:"path"`
	Method      string `json:"method"  form:"method"`
	Description string `json:"description"   form:"description"`
	Group       string `json:"group"  form:"group"`
}

func (SysApi) TableName() string {
	return "sys_apis"
}

// ApiGroup 用户返回给前端使用的结构体
// 一般只在设置casbin权限时候才使用该结构体组装数据提供给前端使用
type ApiGroup struct {
	// 对应 SysApi ID
	ID uint `json:"id"`
	// 对应 SysApi ApiGroup
	ApiGroup string `json:"group"`
	// 对应 SysApi  Description
	Description string `json:"description"`
	// 对应 SysApi Path
	Path string `json:"path,omitempty"`
	// 对应 SysApi Method
	Method string `json:"method"`
	// xx分组下的数据,
	//父ID一样的全部放在一个分组里
	Children []ApiGroup `json:"children"`
}

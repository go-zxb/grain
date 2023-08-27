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

// Upload 文件附件结构体
type Upload struct {
	Model
	// 所属用户
	UID string `json:"uid" xml:"uid" gorm:"comment:用户唯一标识符"`
	// 文件的名称
	FileName string `form:"fileName" json:"fileName" xml:"fileName" gorm:"comment:"`
	// 文件的链接
	FileUrl string `form:"fileUrl" json:"fileUrl" xml:"fileUrl" gorm:"comment:"`
	// 上传的文件用途
	FilePurpose string `form:"filePurpose" json:"filePurpose" xml:"filePurpose" gorm:"comment:文件用途"`
	// 文件的类型
	FileType string `form:"fileType" json:"fileType"`
}

func (Upload) TableName() string {
	return "uploads"
}

// UploadReq 一般用于查询数据
type UploadReq struct {
	PageReq
	FileName  string `form:"fileName" json:"fileName" xml:"fileName" gorm:"comment:"`
	QueryTime string ` form:"queryTime" json:"queryTime"`
}

// UploadRes 一般用于返回响应数据
// 为什么会有 username,nickname 字段?
// 这是为了 在某些情况下 我想知道这条数据的所属者是谁,
// 那么就可以把所属者的信息放在这里就ok了
type UploadRes struct {
	Model
	Nickname string `json:"nickname;omitempty" xml:"nickname" gorm:"-"`
	Username string `json:"username;omitempty" xml:"username" gorm:"-"`
	FileName string `json:"fileName" xml:"fileName" gorm:"comment:"`
	FileUrl  string `json:"fileUrl" xml:"fileUrl" gorm:"comment:"`
}

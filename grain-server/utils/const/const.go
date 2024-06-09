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

package consts

const (
	FormatOrTypeError = 4000
	ReqFail           = 5000
	//系统用户
	CreateUserFail             = 1001
	InvalidParameter           = 1002
	IncorrectAccountORPassword = 1003
	GetUserInfoFail            = 1004
	GetUserByIDInfoFail        = 1005
	GetUserListInfoFail        = 1006
	UpdateUserInfoFail         = 1007
	ModifyPasswordFail         = 1008
	ModifyMobileFail           = 1009
	ModifyEmailFail            = 1010
	EditSysUserFail            = 1011
	SetDefaultRoleFail         = 1012
	DeleteSysUserByIdFail      = 1013
	DeleteSysUserByIdsFail     = 1014
	UploadAvatarFail           = 1015

	//验证码
	SendMobileCaptchaFail     = 1101
	SendEmailCaptchaFail      = 1102
	SendUserEmailCaptchaFail  = 1103
	SendUserMobileCaptchaFail = 1104

	// casbin
	GetAuthApiListFail = 1200
	UpdateCasbinFail   = 1201

	// 系统用户角色
	CreateRoleFail     = 1300
	GetRoleListFail    = 1301
	UpdateRoleFail     = 1302
	DeleteRoleByIdFail = 1303
	DeleteRoleListFail = 1304

	NotRoleList = 1340
)

var (
	TokenBlack = "tokenBlack:"
	UserInfo   = "userInfo:"
)

var Language = 0
var Maps = make(map[int]map[int]string)

func init() {
	Maps[0] = map[int]string{
		IncorrectAccountORPassword: "帐户或密码不正确",
		GetUserInfoFail:            "获取用户信息失败",
		GetUserByIDInfoFail:        "获取用户信息失败",
		GetUserListInfoFail:        "获取用户列表信息失败",
		UpdateUserInfoFail:         "更新用户信息失败",
		ModifyPasswordFail:         "修改密码失败",
		ModifyMobileFail:           "修改手机失败",
		ModifyEmailFail:            "修改电子邮件失败",
		EditSysUserFail:            "编辑系统用户失败",
		SetDefaultRoleFail:         "设置默认角色失败",
		DeleteSysUserByIdFail:      "删除系统用户失败",
		UploadAvatarFail:           "上传头像失败",

		//验证码
		SendMobileCaptchaFail:     "发送手机验证码失败",
		SendEmailCaptchaFail:      "发送电子邮件验证码失败",
		SendUserEmailCaptchaFail:  "发送用户电子邮件验证码失败",
		SendUserMobileCaptchaFail: "发送用户手机验证码失败",

		// casbin
		GetAuthApiListFail: "获取已分配权限的Api列表失败",
		UpdateCasbinFail:   "更新权限失败",

		// 系统用户角色
		CreateRoleFail:  "创建用户角色失败",
		GetRoleListFail: "获取用户角色分页数据失败",
		NotRoleList:     "暂无角色数据",
	}

	Maps[1] = map[int]string{
		IncorrectAccountORPassword: "The account or password is incorrect",
		GetUserInfoFail:            "Failed to get user information",
		GetUserByIDInfoFail:        "Failed to get user information",
		GetUserListInfoFail:        "Failed to get user list information",
		UpdateUserInfoFail:         "Failed to update user information",
		ModifyPasswordFail:         "Failed to change password",
		ModifyMobileFail:           "Failed to modify phone",
		ModifyEmailFail:            "Failed to modify email",
		EditSysUserFail:            "Failed to edit system user",
		SetDefaultRoleFail:         "Setting the default role failed",
		DeleteSysUserByIdFail:      "Failed to delete system user",
		UploadAvatarFail:           "Failed to upload avatar",

		//验证码
		SendMobileCaptchaFail:     "Failed to send phone verification code",
		SendEmailCaptchaFail:      "Failed to send email verification code",
		SendUserEmailCaptchaFail:  "Failed to send user email verification code",
		SendUserMobileCaptchaFail: "Failed to send user's mobile phone verification code",

		// casbin
		GetAuthApiListFail: "Failed to get the list of APIs with assigned permissions",
		UpdateCasbinFail:   "Failed to update permissions",
	}
}

func ErrMsg(err int) string {
	return Maps[Language][err]
}

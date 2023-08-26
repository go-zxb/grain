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

package handler

import (
	"github.com/gin-gonic/gin"
	xjson "github.com/go-grain/go-utils/json"
	"github.com/go-grain/go-utils/response"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/utils"
	"github.com/go-grain/grain/utils/const"
	"github.com/go-grain/grain/utils/upload"
	"strconv"
	"time"
)

type SysUserHandle struct {
	res response.Response
	sv  *service.SysUserService
}

func NewSysUserHandle(sv *service.SysUserService) *SysUserHandle {
	return &SysUserHandle{
		sv: sv,
	}
}

// Login 登录
// @Summary 登录
// @Description 用户登录接口，使用用户名和密码进行登录
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.LoginReq true "用户信息"
// @Success 200  {object} model.LoginRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/login [post]
func (r *SysUserHandle) Login(ctx *gin.Context) {
	reply := r.res.New()
	user := model.LoginReq{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("请求数据有误").Fail(ctx)
		return
	}
	token, err := r.sv.Login(&user, ctx)
	if err != nil {
		reply.WithCode(consts.IncorrectAccountORPassword).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("欢迎回来").WithData(gin.H{"token": token}).Success(ctx)
}

// LogOut 退出登录
// @Security ApiKeyAuth
// @Summary 退出登录
// @Description 退出登录接口
// @Tags 系统用户
// @Accept json
// @Produce json
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/logout [post]
func (r *SysUserHandle) LogOut(ctx *gin.Context) {
	reply := r.res.New()
	err := r.sv.LogOut(ctx)
	if err != nil {
		reply.WithCode(500).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithCode(200).WithMessage("退出登录成功").Success(ctx)
}

// GetLoginUserInfo 获取个人信息
// @Security ApiKeyAuth
// @Summary 获取个人信息
// @Description 获取个人信息接口
// @Tags 系统用户
// @Accept json
// @Produce json
// @Success 200  {object} model.SysUser  "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/info [get]
func (r *SysUserHandle) GetLoginUserInfo(ctx *gin.Context) {
	reply := r.res.New()
	userInfo, err := r.sv.GetLoginUserInfo(ctx)
	if err != nil {
		reply.WithCode(consts.GetUserInfoFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("ok").WithData(userInfo).Success(ctx)
}

// CreateSysUser 创建一个系统用户
// @Security ApiKeyAuth
// @Summary 创建一个系统用户
// @Description 管理员创建一个系统用户
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.CreateSysUser true "用户信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/create [post]
func (r *SysUserHandle) CreateSysUser(ctx *gin.Context) {
	reply := r.res.New()
	c := model.CreateSysUser{}
	err := ctx.ShouldBindJSON(&c)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	s := model.SysUser{}
	err = utils.StructToStruct(&c, &s)
	if err != nil {
		reply.WithCode(consts.CreateUserFail).WithMessage(err.Error()).Fail(ctx)
		return
	}

	err = r.sv.CreateSysUser(&s, ctx)
	if err != nil {
		reply.WithCode(consts.CreateUserFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("创建用户成功").Success(ctx)
}

// GetSysUserById 根据 userID 获取用户信息
// @Security ApiKeyAuth
// @Summary 根据 userID 获取系统用户信息
// @Description 根据 userID 获取系统用户信息
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param id query  int true "用户ID "
// @Success 200 {object} model.SysUser "成功"
// @Failure 400,404 {object} model.ErrorRes "失败"
// @Router /sysUser [get]
func (r *SysUserHandle) GetSysUserById(ctx *gin.Context) {
	reply := r.res.New()
	sysUserId, _ := strconv.Atoi(ctx.Query("id"))
	sysUserInfo, err := r.sv.GetSysUserById(uint(sysUserId), ctx)
	if err != nil {
		reply.WithCode(consts.GetUserByIDInfoFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").WithData(sysUserInfo).Success(ctx)
}

// GetSysUserList 批量获取用户列表
// @Security ApiKeyAuth
// @Summary 批量获取用户列表
// @Description 批量获取用户列表
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.SysUserReq true "分页数据"
// @Success 200  {object} model.SysUser "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/list [get]
func (r *SysUserHandle) GetSysUserList(ctx *gin.Context) {
	reply := r.res.New()
	req := model.SysUserReq{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetSysUserList(&req, ctx)
	if err != nil {
		reply.WithCode(consts.GetUserListInfoFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").WithData(list).Success(ctx)
}

// UpdateSysUser 更新个人信息
// @Security ApiKeyAuth
// @Summary 更新个人信息
// @Description 用户更新个人信息
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.UpdateUserInfo true "用户信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/update [post]
func (r *SysUserHandle) UpdateSysUser(ctx *gin.Context) {
	reply := r.res.New()
	sysUser := model.UpdateUserInfo{}
	err := ctx.ShouldBindJSON(&sysUser)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("解析SysUser参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateSysUser(&sysUser, ctx)
	if err != nil {
		reply.WithCode(consts.UpdateUserInfoFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("更新个人信息成功").Success(ctx)
}

// ModifyPassword 修改密码
// @Security ApiKeyAuth
// @Summary 修改密码
// @Description 用户修改自己的密码
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param sysUser body model.ModifyPassword true "修改密码信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/modifyPassword [put]
func (r *SysUserHandle) ModifyPassword(ctx *gin.Context) {
	reply := r.res.New()
	sysUser := model.ModifyPassword{}
	err := ctx.ShouldBindJSON(&sysUser)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.ModifyPassword(&sysUser, ctx)
	if err != nil {
		reply.WithCode(consts.ModifyPasswordFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("修改密码成功").Success(ctx)
}

// ModifyMobile 修改手机号
// @Security ApiKeyAuth
// @Summary 修改手机号
// @Description 修改手机号
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.ModifyMobile true "用户信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/modifyMobile [put]
func (r *SysUserHandle) ModifyMobile(ctx *gin.Context) {
	reply := r.res.New()
	sysUser := model.ModifyMobile{}
	err := ctx.ShouldBindJSON(&sysUser)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.ModifyMobile(&sysUser, ctx)
	if err != nil {
		reply.WithCode(consts.ModifyMobileFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("修改手机号成功").Success(ctx)
}

// ConfirmModifyEmail 确认修改邮箱
// @Security ApiKeyAuth
// @Summary 确认修改邮箱
// @Description 确认修改邮箱
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param key query string true "key"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/confirmModifyEmail [get]
func (r *SysUserHandle) ConfirmModifyEmail(ctx *gin.Context) {
	key := ctx.Query("key")
	err := r.sv.ConfirmModifyEmail(key, ctx)
	if err != nil {
		_, _ = ctx.Writer.WriteString("错误...")
		return
	}
	_, _ = ctx.Writer.WriteString("修改邮箱成功")
}

// ModifyEmail 提交修改邮箱信息
// @Security ApiKeyAuth
// @Summary 提交修改邮箱信息
// @Description 提交修改邮箱信息
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.ModifyEmail true "邮箱信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/modifyEmail [put]
func (r *SysUserHandle) ModifyEmail(ctx *gin.Context) {
	reply := r.res.New()
	sysUser := model.ModifyEmail{}
	err := ctx.ShouldBindJSON(&sysUser)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.ModifyEmail(&sysUser, ctx)
	if err != nil {
		reply.WithCode(consts.ModifyEmailFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("修改信息已提交,已向目标邮箱发送了确认链接,请进一步完成确认修改").Success(ctx)
}

// EditSysUser 编辑用户信息
// @Security ApiKeyAuth
// @Summary 编辑用户信息
// @Description 编辑用户信息
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.SysUser true "用户信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/editSysUser [put]
func (r *SysUserHandle) EditSysUser(ctx *gin.Context) {
	reply := r.res.New()
	sysUser := model.SysUser{}
	err := ctx.ShouldBindJSON(&sysUser)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.EditUserInfo(&sysUser, ctx)
	if err != nil {
		reply.WithCode(consts.EditSysUserFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("编辑用户信息成功").Success(ctx)
}

// SetDefaultRole 设置默认角色
// @Security ApiKeyAuth
// @Summary 设置默认角色
// @Description 设置默认角色
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param data body model.DefaultRole true "角色信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/setDefaultRole [put]
func (r *SysUserHandle) SetDefaultRole(ctx *gin.Context) {
	reply := r.res.New()
	user := model.SysUser{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.SetDefaultRole(&user, ctx)
	if err != nil {
		reply.WithCode(consts.SetDefaultRoleFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("更新默认角色成功").Success(ctx)
}

// DeleteSysUserById 删除用户
// @Security ApiKeyAuth
// @Summary 删除用户
// @Description 删除用户
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser [put]
func (r *SysUserHandle) DeleteSysUserById(ctx *gin.Context) {
	reply := r.res.New()
	sysUserId, _ := strconv.Atoi(ctx.Query("id"))
	err := r.sv.DeleteSysUserById(uint(sysUserId), ctx)
	if err != nil {
		reply.WithCode(consts.DeleteSysUserByIdFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除用户成功").Success(ctx)
}

// DeleteSysUserByIdList 批量删除用户
// @Security ApiKeyAuth
// @Summary 批量删除用户
// @Description 批量删除用户
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param list path []int true "用户ID列表"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/list [put]
func (r *SysUserHandle) DeleteSysUserByIdList(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		SysUserIds []uint `json:"sysUserIds"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteSysUserByIdList(api.SysUserIds, ctx)
	if err != nil {
		reply.WithCode(consts.DeleteSysUserByIdsFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("批量删除用户成功").Success(ctx)
}

// UploadAvatar 上传头像
// @Security ApiKeyAuth
// @Summary 上传头像
// @Description 上传头像
// @Tags 系统用户
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Param file formData file true "头像"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/modifyAvatar [post]
func (r *SysUserHandle) UploadAvatar(ctx *gin.Context) {
	reply := r.res.New()

	file, err := upload.UploadFile(ctx, "avatar")
	if err != nil {
		reply.WithCode(500).WithMessage(err.Error()).Fail(ctx)
		return
	}

	file.FilePurpose = "系统用户头像"
	err = r.sv.UploadAvatar(file, ctx)
	if err != nil {
		reply.WithCode(consts.UploadAvatarFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("上传文件成功").WithData(xjson.G{"fileUrl": file.FileUrl}).Success(ctx)
}

func (r *SysUserHandle) Certification(ctx *gin.Context) {
	type Record struct {
		CertificationType    string `json:"certificationType,omitempty"`
		CertificationContent string `json:"certificationContent,omitempty"`
		Status               string `json:"status,omitempty"`
		Time                 string `json:"time,omitempty"`
	}
	type name struct {
		AccountType          string    `json:"accountType,omitempty"`
		Status               int       `json:"status,omitempty"`
		Time                 time.Time `json:"time,omitempty"`
		LegalPerson          string    `json:"legalPerson,omitempty"`
		CertificateType      string    `json:"certificateType,omitempty"`
		AuthenticationNumber string    `json:"authenticationNumber,omitempty"`
		Record               []Record  `json:"record"`
	}
	res := name{
		AccountType:          "开发者",
		Status:               1,
		Time:                 time.Now(),
		LegalPerson:          "Grain",
		CertificateType:      "身份证",
		AuthenticationNumber: "5202211997239882",
		Record: []Record{
			{
				CertificationType:    "ddd",
				CertificationContent: "dd",
				Status:               "dd",
				Time:                 "dfd",
			},
		},
	}
	reply := r.res.New()
	reply.WithCode(200).WithMessage("ok").WithData(gin.H{"enterpriseInfo": res}).Success(ctx)
}

// SwitchRole 切换角色
// @Security ApiKeyAuth
// @Summary 切换角色
// @Description 切换角色
// @Tags 系统用户
// @Accept json
// @Produce json
// @Param role query int true "角色"
// @Success 200  {object} model.LoginRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysUser/switchRole [post]
func (r *SysUserHandle) SwitchRole(ctx *gin.Context) {

}

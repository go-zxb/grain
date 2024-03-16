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
	utils "github.com/go-grain/go-utils"
	"github.com/go-grain/go-utils/response"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/utils/const"
)

type RoleHandle struct {
	res response.Response
	sv  *service.RoleService
}

func NewRoleHandle(sv *service.RoleService) *RoleHandle {
	return &RoleHandle{
		sv: sv,
	}
}

func (r *RoleHandle) InitRole() error {
	return r.sv.InitRole()
}

// CreateRole 创建一个系统角色
// @Security ApiKeyAuth
// @Summary 创建一个系统角色
// @Description 管理员创建一个系统角色
// @Tags 系统角色
// @Accept json
// @Produce json
// @Param data body model.CreateSysRole true "用户角色信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysRole [post]
func (r *RoleHandle) CreateRole(ctx *gin.Context) {
	res := r.res.New()
	role := model.CreateSysRole{}
	err := ctx.ShouldBindJSON(&role)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateRole(&role, ctx)
	if err != nil {
		res.WithCode(consts.CreateRoleFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("创建角色成功").Success(ctx)
}

// GetRoleList 批量获取系统角色列表
// @Security ApiKeyAuth
// @Summary 批量获取系统角色列表
// @Description 批量获取系统角色列表
// @Tags 系统角色
// @Accept json
// @Produce json
// @Param data body model.SysUserReq true "分页数据"
// @Success 200  {object} model.SysUser "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysRole/list [post]
func (r *RoleHandle) GetRoleList(ctx *gin.Context) {
	res := r.res.New()
	req := model.SysRoleQueryPage{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetRoleList(&req, ctx)
	if err != nil {
		res.WithCode(consts.GetRoleListFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(list).Success(ctx)
}

// UpdateRole 更新用户角色信息
// @Security ApiKeyAuth
// @Summary 更新用户角色信息
// @Description 更新用户角色信息
// @Tags 系统角色
// @Accept json
// @Produce json
// @Param data body model.UpdateUserInfo true "用户角色信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysRole/update [post]
func (r *RoleHandle) UpdateRole(ctx *gin.Context) {
	res := r.res.New()
	role := model.SysRole{}
	err := ctx.ShouldBindJSON(&role)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateRole(&role, ctx)
	if err != nil {
		res.WithCode(consts.UpdateRoleFail).WithMessage("更新角色失败").Fail(ctx)
		return
	}
	res.WithMessage("更新角色成功").Success(ctx)
}

// DeleteRoleById 根据角色ID删除系统角色
// @Security ApiKeyAuth
// @Summary 删除系统角色
// @Description 删除系统角色
// @Tags 系统角色
// @Accept json
// @Produce json
// @Param ID query int true "角色ID"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysRole [delete]
func (r *RoleHandle) DeleteRoleById(ctx *gin.Context) {
	reply := r.res.New()
	role := utils.String2Int(ctx.Query("id"))
	if role == 0 {
		reply.WithCode(consts.ReqFail).WithMessage("ID不能为空").Fail(ctx)
		return
	}
	err := r.sv.DeleteRoleById(uint(role), ctx)
	if err != nil {
		reply.WithCode(consts.DeleteRoleByIdFail).WithMessage("删除角色失败").Fail(ctx)
		return
	}
	reply.WithMessage("删除角色成功").Success(ctx)
}

// DeleteRoleByIds 批量删除系统角色
// @Security ApiKeyAuth
// @Summary 批量删除系统角色
// @Description 批量删除系统角色
// @Tags 系统角色
// @Accept json
// @Produce json
// @Param data body []int true "系统角色ID列表"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysRole/list [delete]
func (r *RoleHandle) DeleteRoleByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		Roles []uint `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage("IDArray不能为空").Fail(ctx)
		return
	}
	err = r.sv.DeleteRoleByIds(api.Roles, ctx)
	if err != nil {
		reply.WithCode(consts.DeleteRoleListFail).WithMessage("批量删除角色失败").Fail(ctx)
		return
	}
	reply.WithMessage("批量删除角色成功").Success(ctx)
}

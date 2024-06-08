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
	"fmt"
	"github.com/gin-gonic/gin"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/pkg/response"
	consts "github.com/go-grain/grain/utils/const"
	"strconv"
)

type MenuHandle struct {
	res response.Response
	sv  *service.MenuService
}

func NewMenuHandle(sv *service.MenuService) *MenuHandle {
	return &MenuHandle{
		sv: sv,
	}
}

func (r *MenuHandle) InitMenu() error {
	return r.sv.InitMenu()
}

// CreateMenu
// @Security ApiKeyAuth
// @Summary 创建菜单
// @Description 创建菜单
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Param data body  model.SysMenu true "菜单信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /sysMenu [post]
func (r *MenuHandle) CreateMenu(ctx *gin.Context) {
	res := r.res.New()
	menu := model.SysMenu{}
	err := ctx.ShouldBindJSON(&menu)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateMenu(&menu, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("创建菜单成功").Success(ctx)
}

// GetUserMenu
// @Security ApiKeyAuth
// @Summary 用户获取动态菜单
// @Description 用户获取动态菜单
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Success 200 {object} model.SysMenu "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu/userMenu [get]
func (r *MenuHandle) GetUserMenu(ctx *gin.Context) {
	res := r.res.New()
	role := ctx.GetString("role")
	menuInfo, err := r.sv.GetUserMenu(role, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(menuInfo).Success(ctx)
}

// GetMenuAndPermission
// @Security ApiKeyAuth
// @Summary 设置菜单权限
// @Description 设置菜单权限
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Success 200 {object} model.SysMenu "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu/userMenu [get]
func (r *MenuHandle) GetMenuAndPermission(ctx *gin.Context) {
	res := r.res.New()
	role := ctx.Query("role")
	menuInfo, selectKeys, err := r.sv.GetMenuAndPermission(role, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(menuInfo).WithData2(selectKeys).Success(ctx)
}

// SetMenuAndPermission
// @Security ApiKeyAuth
// @Summary 设置菜单权限
// @Description 设置菜单权限
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Success 200 {object} model.SysMenu "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu/userMenu [get]
func (r *MenuHandle) SetMenuAndPermission(ctx *gin.Context) {
	reply := r.res.New()
	type Menu struct {
		Role string `form:"role" json:"role"`
		Keys []uint `form:"keys" json:"keys"`
	}
	Keys := Menu{}
	if err := ctx.ShouldBindJSON(&Keys); err != nil {
		reply.WithCode(500).WithMessage(err.Error()).Fail(ctx)
		return
	}
	fmt.Println(Keys)
	err := r.sv.SetMenuAndPermission(Keys.Keys, Keys.Role)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").Success(ctx)
}

// GetMenuList
// @Security ApiKeyAuth
// @Summary 获取菜单分页数据
// @Description 获取菜单分页数据
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Param data body model.SysMenuReq true "分页列表请求参数"
// @Success 200 {object} model.SysMenu "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu/list [get]
func (r *MenuHandle) GetMenuList(ctx *gin.Context) {
	res := r.res.New()
	req := model.SysMenuReq{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetMenuList(&req, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(list).Success(ctx)
}

// UpdateMenu
// @Security ApiKeyAuth
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Param data body  model.SysMenu true "更新菜单信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu [put]
func (r *MenuHandle) UpdateMenu(ctx *gin.Context) {
	res := r.res.New()
	menu := model.SysMenu{}
	err := ctx.ShouldBindJSON(&menu)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateMenu(&menu, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("更新菜单成功").Success(ctx)
}

// DeleteMenuById
// @Security ApiKeyAuth
// @Summary 删除一条菜单数据
// @Description 根据菜单ID删除菜单
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Param id query  int true "根据菜单ID删除菜单 "
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu [delete]
func (r *MenuHandle) DeleteMenuById(ctx *gin.Context) {
	res := r.res.New()
	menuId, _ := strconv.Atoi(ctx.Query("id"))
	err := r.sv.DeleteMenuById(uint(menuId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("删除菜单成功").Success(ctx)
}

// DeleteMenuByIds
// @Security ApiKeyAuth
// @Summary 批量删除菜单
// @Description "根据菜单ID批量删除菜单"
// @Tags 动态菜单
// @Accept json
// @Produce json
// @Param data body []uint true "根据菜单ID批量删除菜单"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysMenu/deleteMenuByIds [delete]
func (r *MenuHandle) DeleteMenuByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		MenuIds []uint `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteMenuByIds(api.MenuIds, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除菜单成功").Success(ctx)
}

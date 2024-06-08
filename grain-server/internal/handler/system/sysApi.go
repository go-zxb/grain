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
	service "github.com/go-grain/grain/internal/service/system"
	model "github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/pkg/convert"
	"github.com/go-grain/grain/pkg/response"
	"github.com/go-grain/grain/utils/const"
)

type ApiHandle struct {
	res response.Response
	sv  *service.ApiService
}

func NewApiHandle(sv *service.ApiService) *ApiHandle {
	return &ApiHandle{
		sv: sv,
	}
}

func (r *ApiHandle) InitApi() error {
	return r.sv.InitApi()
}

// CreateApi 创建一个API接口
// @Security ApiKeyAuth
// @Summary 创建一个API接口
// @Description 创建一个API接口
// @Tags API接口
// @Accept json
// @Produce json
// @Param data body model.SysApi true "API接口信息"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi [post]
func (r *ApiHandle) CreateApi(ctx *gin.Context) {
	res := r.res.New()
	api := model.SysApi{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateApi(&api, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("创建Api成功").Success(ctx)
}

// GetApiList 批量获取API接口列表
// @Security ApiKeyAuth
// @Summary 批量获取API接口列表
// @Description 批量获取API接口列表
// @Tags API接口
// @Accept json
// @Produce json
// @Param data body model.SysApiReq true "分页数据"
// @Success 200  {object} []model.SysApi "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi/list [get]
func (r *ApiHandle) GetApiList(ctx *gin.Context) {
	res := r.res.New()
	req := model.SysApiReq{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetApiList(&req)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").
		WithData(list).
		WithTotal(req.Total).
		WithPage(req.Page).
		WithPageSize(req.PageSize).
		Success(ctx)
}

// GetApiAndPermissions 批量获取所有API接口和已授权的api列表
// @Security ApiKeyAuth
// @Summary 批量获取所有API接口和已授权的api列表
// @Description 批量获取所有API接口和已授权的api列表
// @Tags API接口
// @Accept json
// @Produce json
// @Success 200  {object} model.SysUser "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi/apiAndPermissions [get]
func (r *ApiHandle) GetApiAndPermissions(ctx *gin.Context) {
	res := r.res.New()
	role := ctx.Query("role")
	list, err := r.sv.GetApiAndPermissions(role)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").
		WithData(list).
		Success(ctx)
}

// GetApiGroup 批量获取API接口分组
// @Security ApiKeyAuth
// @Summary 批量获取API接口分组
// @Description 批量获取API接口分组
// @Tags API接口
// @Accept json
// @Produce json
// @Success 200  {object} model.SysUser "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi/apiGroups [get]
func (r *ApiHandle) GetApiGroup(ctx *gin.Context) {
	res := r.res.New()
	list, err := r.sv.GetApiGroup()
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").
		WithData(list).
		Success(ctx)
}

// UpdateApi 更新APi接口
// @Security ApiKeyAuth
// @Summary 更新APi接口
// @Description 更新APi接口
// @Tags API接口
// @Accept json
// @Produce json
// @Param sysUser body model.SysApi true "APi接口信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi [put]
func (r *ApiHandle) UpdateApi(ctx *gin.Context) {
	res := r.res.New()
	user := model.SysApi{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage("解析Api参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateApi(&user, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("更新Api成功").Success(ctx)
}

// DeleteApiById 根据API接口ID删除API接口
// @Security ApiKeyAuth
// @Summary 根据API接口ID删除API接口
// @Description 根据API接口ID删除API接口
// @Tags API接口
// @Accept json
// @Produce json
// @Param ID query int true "API接口ID"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi [delete]
func (r *ApiHandle) DeleteApiById(ctx *gin.Context) {
	reply := r.res.New()
	id := convert.String2Int(ctx.Query("id"))
	if id == 0 {
		reply.WithCode(consts.InvalidParameter).WithMessage("参数不能为空").Fail(ctx)
		return
	}
	err := r.sv.DeleteApiById(uint(id), ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除Api成功").Success(ctx)
}

// DeleteApiByIds 根据API接口ID删除API接口
// @Security ApiKeyAuth
// @Summary 根据API接口ID删除API接口
// @Description 根据API接口ID删除API接口
// @Tags API接口
// @Accept json
// @Produce json
// @Param data body  []int true "API接口Ids"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /sysApi/deleteApiByIds [delete]
func (r *ApiHandle) DeleteApiByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		Ids []uint `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteApiByIds(api.Ids, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除Api成功").Success(ctx)
}

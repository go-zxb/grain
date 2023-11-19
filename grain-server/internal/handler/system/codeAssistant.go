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
	consts "github.com/go-grain/grain/utils/const"
	"strconv"
)

type CodeAssistantHandle struct {
	res response.Response
	sv  *service.CodeAssistantService
}

func NewCodeAssistantHandle(sv *service.CodeAssistantService) *CodeAssistantHandle {
	return &CodeAssistantHandle{
		sv: sv,
	}
}

// CreateProject
// @Security ApiKeyAuth
// @Summary 创建项目
// @Description 创建项目
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param data body  model.CreateProject true "项目信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /project [post]
func (r *CodeAssistantHandle) CreateProject(ctx *gin.Context) {
	res := r.res.New()
	codeFactory := model.Project{}
	err := ctx.ShouldBindJSON(&codeFactory)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateProject(&codeFactory, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("操作成功").Success(ctx)
}

// UpdateProject
// @Security ApiKeyAuth
// @Summary 更新项目
// @Description 更新项目
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param data body  model.CreateProject true "项目信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /project [put]
func (r *CodeAssistantHandle) UpdateProject(ctx *gin.Context) {
	res := r.res.New()
	codeFactory := model.Project{}
	err := ctx.ShouldBindJSON(&codeFactory)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.UpdateProject(&codeFactory, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("操作成功").Success(ctx)
}

// CreateModel
// @Security ApiKeyAuth
// @Summary 创建模型
// @Description 创建模型
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param data body  model.CreateModels true "模型信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /models [post]
func (r *CodeAssistantHandle) CreateModel(ctx *gin.Context) {
	res := r.res.New()
	codeFactory := model.Models{}
	err := ctx.ShouldBindJSON(&codeFactory)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateModel(&codeFactory, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("操作成功").Success(ctx)
}

// UpdateModel
// @Security ApiKeyAuth
// @Summary 更新模型
// @Description 更新模型
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param data body  model.CreateModels true "模型信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /models [put]
func (r *CodeAssistantHandle) UpdateModel(ctx *gin.Context) {
	res := r.res.New()
	codeFactory := model.Models{}
	err := ctx.ShouldBindJSON(&codeFactory)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.UpdateModel(&codeFactory, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("操作成功").Success(ctx)
}

// CreateField
// @Security ApiKeyAuth
// @Summary 创建字段
// @Description 创建字段
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param data body  model.Fields true "字段信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /field [post]
func (r *CodeAssistantHandle) CreateField(ctx *gin.Context) {
	res := r.res.New()
	codeFactory := model.Fields{}
	err := ctx.ShouldBindJSON(&codeFactory)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateField(&codeFactory, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("操作成功").Success(ctx)
}

// UpdateField
// @Security ApiKeyAuth
// @Summary 更新字段
// @Description 更新字段
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param data body  model.Fields true "字段信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /field [put]
func (r *CodeAssistantHandle) UpdateField(ctx *gin.Context) {
	res := r.res.New()
	codeFactory := model.Fields{}
	err := ctx.ShouldBindJSON(&codeFactory)
	if err != nil {
		res.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.UpdateField(&codeFactory, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("操作成功").Success(ctx)
}

// ViewCode
// @Security ApiKeyAuth
// @Summary 预览代码
// @Description 预览代码
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param mid query  int true "模型ID "
// @Success 200 {object} model.ViewCode "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /viewCode [get]
func (r *CodeAssistantHandle) ViewCode(ctx *gin.Context) {
	res := r.res.New()
	mid := ctx.Query("mid")
	data, err := r.sv.ViewCode(uint(utils.String2Int(mid)), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}

	res.WithMessage("ok").WithData(data).Success(ctx)
}

// GenerateCode
// @Security ApiKeyAuth
// @Summary 生成代码
// @Description 生成代码
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param mid,fb query  int true "模型ID "
// @Success 200 {object} model.ViewCode "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /generateCode [post]
func (r *CodeAssistantHandle) GenerateCode(ctx *gin.Context) {
	res := r.res.New()
	mid := utils.String2Int(ctx.Query("mid"))
	forceBuild := utils.String2Int(ctx.Query("fb"))
	err := r.sv.GenerateCode(uint(mid), uint(forceBuild), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("生成代码成功").Success(ctx)
}

// GetProjectList
// @Security ApiKeyAuth
// @Summary 获取项目分页数据
// @Description 获取项目分页数据
// @Tags 代码助手
// @Accept json
// @Produce json
// @Success 200 {object} []model.Project "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /project/list [get]
func (r *CodeAssistantHandle) GetProjectList(ctx *gin.Context) {
	reply := r.res.New()
	list, err := r.sv.GetProjectList(ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").
		WithData(list).
		WithPage(1).
		WithPageSize(1).
		Success(ctx)
}

// GetModelList
// @Security ApiKeyAuth
// @Summary 获取模型数据
// @Description 获取模型数据
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param parentId query  int true "模型列表父ID"
// @Success 200 {object} model.Models "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /models/list [get]
func (r *CodeAssistantHandle) GetModelList(ctx *gin.Context) {
	res := r.res.New()
	parentId, _ := strconv.Atoi(ctx.Query("parentId"))
	list, err := r.sv.GetModelList(uint(parentId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(list).WithPage(1).Success(ctx)
}

// GetFieldList
// @Security ApiKeyAuth
// @Summary 获取字段分页数据
// @Description 获取字段分页数据
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param parentId query  int true "字段列表父ID "
// @Success 200 {object} model.Field "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /field/list [get]
func (r *CodeAssistantHandle) GetFieldList(ctx *gin.Context) {
	res := r.res.New()
	parentId, _ := strconv.Atoi(ctx.Query("parentId"))
	list, err := r.sv.GetFieldList(uint(parentId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(list).Success(ctx)
}

// DeleteProjectById
// @Security ApiKeyAuth
// @Summary 删除项目
// @Description 根据项目ID删除项目
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param pid query  int true "根据项目ID删除项目 "
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /project [delete]
func (r *CodeAssistantHandle) DeleteProjectById(ctx *gin.Context) {
	res := r.res.New()
	menuId, _ := strconv.Atoi(ctx.Query("pid"))
	err := r.sv.DeleteProjectByIds(uint(menuId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("删除项目成功").Success(ctx)
}

// DeleteModelById
// @Security ApiKeyAuth
// @Summary 删除模型
// @Description 根据模型ID删除模型
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param mid query  int true "根据模型ID删除模型 "
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /models [delete]
func (r *CodeAssistantHandle) DeleteModelById(ctx *gin.Context) {
	res := r.res.New()
	menuId, _ := strconv.Atoi(ctx.Query("mid"))
	err := r.sv.DeleteModelById(uint(menuId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("删除项目成功").Success(ctx)
}

// DeleteFieldById
// @Security ApiKeyAuth
// @Summary 删除字段
// @Description 根据字段ID删除字段
// @Tags 代码助手
// @Accept json
// @Produce json
// @Param fid query  int true "根据字段ID删除字段 "
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /field [delete]
func (r *CodeAssistantHandle) DeleteFieldById(ctx *gin.Context) {
	res := r.res.New()
	menuId, _ := strconv.Atoi(ctx.Query("fid"))
	err := r.sv.DeleteFieldById(uint(menuId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("删除项目成功").Success(ctx)
}

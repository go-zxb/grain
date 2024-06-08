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
	"github.com/go-grain/grain/pkg/response"
	consts "github.com/go-grain/grain/utils/const"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SysLogHandle struct {
	res response.Response
	sv  *service.SysLogService
}

func NewSysLogHandle(sv *service.SysLogService) *SysLogHandle {
	return &SysLogHandle{
		sv: sv,
	}
}

// GetSysLogList
// @Security ApiKeyAuth
// @Summary 获取系统日志分页数据
// @Description 获取系统日志分页数据
// @Tags 系统日志
// @Accept json
// @Produce json
// @Param data body model.SysLogReq true "分页列表请求参数"
// @Success 200 {object} model.SysLog "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysLog/list [get]
func (r *SysLogHandle) GetSysLogList(ctx *gin.Context) {
	res := r.res.New()
	req := model.SysLogReq{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetSysLogList(&req, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithTotal(req.Total).WithData(list).Success(ctx)
}

// DeleteSysLogById
// @Security ApiKeyAuth
// @Summary 删除系统日志
// @Description 根据系统日志ID删除系统日志
// @Tags 系统日志
// @Accept json
// @Produce json
// @Param id query  int true "根据系统日志ID删除系统日志 "
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysLog [delete]
func (r *SysLogHandle) DeleteSysLogById(ctx *gin.Context) {
	reply := r.res.New()
	id := ctx.Query("id")
	if id == "" {
		reply.WithCode(consts.InvalidParameter).WithMessage("ID不能为空").Fail(ctx)
		return
	}
	err := r.sv.DeleteSysLogById(id, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除操作日志成功").Success(ctx)
}

// DeleteSysLogByIds
// @Security ApiKeyAuth
// @Summary 删除系统日志
// @Description "根据系统日志ID批量删除系统日志"
// @Tags 系统日志
// @Accept json
// @Produce json
// @Param data body []uint true "根据系统日志ID批量删除系统日志"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /sysLog/list [delete]
func (r *SysLogHandle) DeleteSysLogByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		SysLogIds []primitive.ObjectID `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteSysLogByIds(api.SysLogIds, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage("删除操作日志失败").Fail(ctx)
		return
	}
	reply.WithMessage("删除操作日志成功").Success(ctx)
}

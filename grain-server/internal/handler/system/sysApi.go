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
	model "github.com/go-grain/grain/model/system"
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

func (r *ApiHandle) CreateApi(ctx *gin.Context) {
	res := r.res.New()
	api := model.SysApi{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateApi(&api, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("创建Api成功").Success(ctx)
}

func (r *ApiHandle) GetApiList(ctx *gin.Context) {
	res := r.res.New()
	req := model.SysApiReq{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage("参数解析失败").Fail(ctx)
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

func (r *ApiHandle) UpdateApi(ctx *gin.Context) {
	res := r.res.New()
	user := model.SysApi{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage("解析Api参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateApi(&user, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("更新Api成功").Success(ctx)
}

func (r *ApiHandle) DeleteApiById(ctx *gin.Context) {
	reply := r.res.New()
	id := utils.String2Int(ctx.Query("id"))
	if id == 0 {
		reply.WithCode(consts.ReqFail).WithMessage("参数不能为空").Fail(ctx)
		return
	}
	err := r.sv.DeleteApiById(uint(id), ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除Api成功").Success(ctx)
}

func (r *ApiHandle) DeleteApiByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		Ids []uint `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteApiByIds(api.Ids, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("删除Api成功").Success(ctx)
}

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
	"github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/pkg/response"
	"github.com/go-grain/grain/utils/const"
)

type CasbinHandle struct {
	res response.Response
	sv  *service.CasbinService
}

func NewCasbinHandle(sv *service.CasbinService) *CasbinHandle {
	return &CasbinHandle{sv: sv}
}

func (r *CasbinHandle) InitCasbinHandle() error {
	return r.sv.InitCasbinRoleRule()
}

// Update 更新角色权限
// @Security ApiKeyAuth
// @Summary 更新角色权限
// @Description 更新角色权限
// @Tags Casbin权限
// @Accept json
// @Produce json
// @Param sysUser body model.CasbinReq true "分配角色权限"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /casbin [put]
func (r *CasbinHandle) Update(ctx *gin.Context) {
	reply := r.res.New()
	var c model.CasbinReq
	err := ctx.ShouldBindJSON(&c)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}

	err = r.sv.Update(&c, ctx)
	if err != nil {
		reply.WithCode(consts.UpdateCasbinFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("更新权限成功").Success(ctx)
}

// AuthApiList 获取xx角色可访问的接口列表
// @Security ApiKeyAuth
// @Summary 获取xx角色可访问的接口列表
// @Description 获取xx角色可访问的接口列表
// @Tags Casbin权限
// @Accept json
// @Produce json
// @Param role path string true "根据Role获取xx角色可访问的接口列表"
// @Success 200 {object} model.CasbinRule "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /casbin/authApiList [get]
func (r *CasbinHandle) AuthApiList(ctx *gin.Context) {
	reply := r.res.New()
	role := ctx.Query("role")
	if role == "" {
		reply.WithCode(consts.ReqFail).WithMessage("请求参数有误").Fail(ctx)
		return
	}
	list, err := r.sv.AuthApiList(role)
	if err != nil {
		reply.WithCode(consts.GetAuthApiListFail).WithMessage(err.Error()).Fail(ctx)
		return
	}

	reply.WithMessage("ok").WithData(list).Success(ctx)
}

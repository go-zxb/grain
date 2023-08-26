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
	_ "github.com/go-grain/grain/model/system"
	model "github.com/go-grain/grain/model/system"
	consts "github.com/go-grain/grain/utils/const"
	"github.com/go-grain/grain/utils/upload"
	"strconv"
)

type UploadHandle struct {
	res response.Response
	sv  *service.UploadService
}

func NewUploadHandle(sv *service.UploadService) *UploadHandle {
	return &UploadHandle{
		sv: sv,
	}
}

// UploadFile 上传文件
// @Security ApiKeyAuth
// @Summary 上传文件
// @Description 上传文件
// @Tags 上传文件
// @Accept json
// @Produce json
// @Accept multipart/form-data
// @Param file formData file true "文件"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /upload [post]
func (r *UploadHandle) UploadFile(ctx *gin.Context) {
	reply := r.res.New()

	file, err := upload.UploadFile(ctx, "systemFile")
	if err != nil {
		reply.WithCode(500).WithMessage(err.Error()).Fail(ctx)
		return
	}

	file.FilePurpose = "普通文件"
	err = r.sv.CreateUpload(file, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("上传文件成功").WithData(xjson.G{"fileUrl": file.FileUrl}).Success(ctx)
}

// GetUploadList
// @Security ApiKeyAuth
// @Summary 获取文件列表
// @Description 获取文件列表
// @Tags 上传文件
// @Accept json
// @Produce json
// @Param data body model.UploadReq true "分页数据"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /upload [get]
func (r *UploadHandle) GetUploadList(ctx *gin.Context) {
	res := r.res.New()
	req := model.UploadReq{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetUploadList(&req, ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(list).Success(ctx)
}

// DeleteUploadById
// @Security ApiKeyAuth
// @Summary 删除文件
// @Description 删除文件
// @Tags 上传文件
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /upload [delete]
func (r *UploadHandle) DeleteUploadById(ctx *gin.Context) {
	res := r.res.New()
	uploadId, _ := strconv.Atoi(ctx.Query("id"))
	err := r.sv.DeleteUploadById(uint(uploadId), ctx)
	if err != nil {
		res.WithCode(consts.ReqFail).WithMessage("删除文件失败").Fail(ctx)
		return
	}
	res.WithMessage("删除文件成功").Success(ctx)
}

// DeleteUploadByIds 批量删除用户
// @Security ApiKeyAuth
// @Summary 批量删除文件
// @Description 批量删除文件
// @Tags 上传文件
// @Accept json
// @Produce json
// @Param list path []int true "文件ID列表"
// @Success 200  {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /upload/list [delete]
func (r *UploadHandle) DeleteUploadByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		UploadIds []uint `json:"uploadIds"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteUploadByIds(api.UploadIds, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage("删除文件失败").Fail(ctx)
		return
	}
	reply.WithMessage("删除文件成功").Success(ctx)
}

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
	"github.com/go-grain/grain/utils/const"
)

type CaptchaHandle struct {
	res response.Response
	sv  *service.CaptchaService
}

func NewCaptchaHandle(sv *service.CaptchaService) *CaptchaHandle {
	return &CaptchaHandle{
		sv: sv,
	}
}

// SendMobileCaptcha 向指定的手机号发送验证码
// @Security ApiKeyAuth
// @Summary 向指定的手机号发送验证码
// @Description 向指定的手机号发送验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param sysUser body model.Mobile true "手机号"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /captcha/sendMobileCaptcha [post]
func (r *CaptchaHandle) SendMobileCaptcha(ctx *gin.Context) {
	reply := r.res.New()
	mobile := model.Mobile{}
	err := ctx.ShouldBindJSON(&mobile)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.SendMobileCaptcha(&mobile, ctx)
	if err != nil {
		reply.WithCode(consts.SendMobileCaptchaFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("验证码发送成功").Success(ctx)
}

// SendEmailCaptcha 向指定的邮箱地址发送验证码
// @Security ApiKeyAuth
// @Summary 向指定的邮箱地址发送验证码
// @Description 向指定的邮箱地址发送验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param sysUser body model.Email true "邮箱地址"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /captcha/sendEmailCaptcha [post]
func (r *CaptchaHandle) SendEmailCaptcha(ctx *gin.Context) {
	reply := r.res.New()
	email := model.Email{}
	err := ctx.ShouldBindJSON(&email)
	if err != nil {
		reply.WithCode(consts.InvalidParameter).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.SendEmailCaptcha(&email, ctx)
	if err != nil {
		reply.WithCode(consts.SendEmailCaptchaFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("验证码发送成功").Success(ctx)
}

// SendUserEmailCaptcha 发送用户邮箱验证码
// @Security ApiKeyAuth
// @Summary 发送用户邮箱验证码
// @Description 发送用户邮箱验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /captcha/sendUserEmailCaptcha [post]
func (r *CaptchaHandle) SendUserEmailCaptcha(ctx *gin.Context) {
	reply := r.res.New()
	err := r.sv.SendUserEmailCaptcha(ctx)
	if err != nil {
		reply.WithCode(consts.SendUserEmailCaptchaFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("验证码发送成功").Success(ctx)
}

// SendUserMobileCaptcha 发送用户手机验证码
// @Security ApiKeyAuth
// @Summary 发送用户手机验证码
// @Description 发送用户手机验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Param sysUser body model.Mobile true "手机号"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 500  {object} model.ErrorRes "失败"
// @Router /captcha/sendUserMobileCaptcha [post]
func (r *CaptchaHandle) SendUserMobileCaptcha(ctx *gin.Context) {
	reply := r.res.New()
	err := r.sv.SendUserMobileCaptcha(ctx)
	if err != nil {
		reply.WithCode(consts.SendUserMobileCaptchaFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("验证码发送成功").Success(ctx)
}

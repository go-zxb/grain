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

package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	utils "github.com/go-grain/go-utils"
	gEmail "github.com/go-grain/go-utils/email"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/jordan-wright/email"
	"strings"
)

type CaptchaService struct {
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Logger
}

func NewCaptcha(rdb redis.IRedis, conf *config.Config, log *log.Logger) *CaptchaService {
	return &CaptchaService{rdb: rdb, conf: conf, log: log}
}

func (s *CaptchaService) SendMobileCaptcha(mobile *model.Mobile, ctx *gin.Context) error {
	mobile.Mobile = strings.TrimSpace(mobile.Mobile)

	i, err := s.rdb.GetInt(fmt.Sprintf("captchaCount:%s", ctx.ClientIP()))
	if err != nil {
		return errors.New("获取验证码失败 服务器内部错误")
	}
	if i >= 5 {
		return errors.New("频繁请求")
	}

	captcha := utils.RandomString(config.GetConfig().System.CaptchaLength)

	err = s.rdb.SetInt(fmt.Sprintf("captcha:%s:%d", ctx.ClientIP(), captcha), captcha, 300)
	if err != nil {
		return errors.New("获取验证码失败")
	}
	xlog.Info(captcha)
	//调用第三方接口发送 手机验证码
	//以后在实现
	return nil
}

func (s *CaptchaService) SendUserEmailCaptcha(ctx *gin.Context) error {
	sysUser, err := query.SysUser.Where(query.SysUser.UID.Eq(ctx.GetString("uid"))).First()
	if err != nil {
		return err
	}

	i, err := s.rdb.GetInt(fmt.Sprintf("captchaCount:%s", ctx.GetString("uid")))
	if err != nil {
		return errors.New("获取验证码失败 服务器内部错误")
	}
	if i >= 5 {
		return errors.New("频繁请求")
	}

	captcha := utils.RandomString(config.GetConfig().System.CaptchaLength)

	err = s.rdb.SetInt(fmt.Sprintf("captcha:%s:%d", ctx.GetString("uid"), captcha), captcha, 300)
	if err != nil {
		return errors.New("获取验证码失败")
	}

	err = s.Send(sysUser.Email, captcha, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *CaptchaService) SendUserMobileCaptcha(ctx *gin.Context) error {
	mobile, err := query.SysUser.Where(query.SysUser.UID.Eq(ctx.GetString("uid"))).First()
	if err != nil {
		return err
	}

	i, err := s.rdb.GetInt(fmt.Sprintf("captchaCount:%s", ctx.GetString("uid")))
	if err != nil {
		return errors.New("获取验证码失败 服务器内部错误")
	}
	if i >= 5 {
		return errors.New("频繁请求")
	}

	captcha := utils.RandomString(config.GetConfig().System.CaptchaLength)

	err = s.rdb.SetInt(fmt.Sprintf("captcha:%s:%d", ctx.GetString("uid"), captcha), captcha, 300)
	if err != nil {
		return errors.New("获取验证码失败")
	}

	xlog.Info(captcha, mobile.Mobile)
	//调用第三方接口发送 手机验证码
	//以后在实现
	return nil
}

func (s *CaptchaService) SendEmailCaptcha(req *model.Email, ctx *gin.Context) error {
	req.Email = strings.TrimSpace(req.Email)

	i, err := s.rdb.GetInt(fmt.Sprintf("captchaCount:%s", ctx.ClientIP()))
	if err != nil {
		return errors.New("获取验证码失败 服务器内部错误")
	}
	if i >= 5 {
		return errors.New("频繁请求")
	}

	captcha := utils.RandomString(config.GetConfig().System.CaptchaLength)

	err = s.rdb.SetInt(fmt.Sprintf("captcha:%s:%d", ctx.ClientIP(), captcha), captcha, 300)
	if err != nil {
		return errors.New("获取验证码失败")
	}

	err = s.Send(req.Email, captcha, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *CaptchaService) Send(xemail string, captcha int64, ctx *gin.Context) error {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = config.GetConfig().System.SiteName + "<" + s.conf.SysEmail.EmailUsername + ">"
	// 设置接收方的邮箱
	e.To = []string{xemail}
	//设置主题
	e.Subject = "邮箱验证码:5分钟内有效"
	//设置文件发送的内容
	e.HTML = []byte(fmt.Sprintf("<html><body><h4>你的验证码是: <br>%d</h4></body></html>", captcha))

	emailServer, err := gEmail.NewMailServer(s.conf.SysEmail.EmailUsername, s.conf.SysEmail.EmailPassword, s.conf.SysEmail.EmailHost, s.conf.SysEmail.EmailHost, utils.Int2String(s.conf.SysEmail.EmailPort))
	if err != nil {
		xlog.Error("初始化Email服务失败", emailServer)
		return errors.New("获取验证码失败 请检查邮箱服务配置")
	}

	err = emailServer.SendEmail(e)
	if err != nil {
		xlog.Error(err)
		return errors.New("获取验证码失败")
	}

	c := ctx.GetString("uid")
	if c == "" {
		c = ctx.ClientIP()
	}

	_, err = s.rdb.IncrInt(fmt.Sprintf("captchaCount:%s", c), 1)
	if err != nil {
		return err
	}
	s.rdb.SetEx(fmt.Sprintf("captchaCount:%s", ctx.GetString("uid")), 1800)

	// debug 模式 打印在控制台 方便查看
	if s.conf.Gin.Model == "debug" {
		xlog.Info(captcha)
	}
	return nil
}

func (s *CaptchaService) CustomEmail(req *model.Email, subject, HTML string) error {
	req.Email = strings.TrimSpace(req.Email)
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = config.GetConfig().System.SiteName + "<" + s.conf.SysEmail.EmailUsername + ">"
	// 设置接收方的邮箱
	e.To = []string{req.Email}
	//设置主题
	e.Subject = subject
	//设置文件发送的内容
	e.HTML = []byte(HTML)

	emailServer, err := gEmail.NewMailServer(s.conf.SysEmail.EmailUsername, s.conf.SysEmail.EmailPassword, s.conf.SysEmail.EmailHost, s.conf.SysEmail.EmailHost, utils.Int2String(s.conf.SysEmail.EmailPort))
	if err != nil {
		xlog.Error("初始化Email服务失败", emailServer)
		return errors.New("请检查邮箱服务配置")
	}

	err = emailServer.SendEmail(e)
	if err != nil {
		return errors.New("发送邮箱信息失败")
	}

	return nil
}

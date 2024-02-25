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
	"github.com/go-grain/go-utils/encrypt"
	xjson "github.com/go-grain/go-utils/json"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/utils/const"
	"net/url"
	"strings"
	"time"
)

type ISysUserRepo interface {
	Login(user *model.LoginReq) (*model.SysUser, error)
	CreateSysUser(user *model.SysUser) error
	GetSysUserById(id uint) (u *model.SysUser, err error)
	GetSysUserByUId(uid string) (u *model.SysUser, err error)
	GetSysUserList(req *model.SysUserReq) ([]*model.SysUser, error)
	UpdateSysUser(user *model.UpdateUserInfo) error
	EditSysUser(user *model.SysUser) error
	SetDefaultRole(user *model.SysUser) error
	DeleteSysUserById(userId uint) error
	DeleteSysUserByIds(userIds []uint) error
	UploadAvatar(avatar *model.Upload, uid string) error
}

type SysUserService struct {
	repo    ISysUserRepo
	rdb     redis.IRedis
	conf    *config.Config
	log     *log.Logger
	captcha *CaptchaService
}

func NewSysUserService(repo ISysUserRepo, rdb redis.IRedis, conf *config.Config, logger *log.Logger) *SysUserService {
	return &SysUserService{
		repo:    repo,
		rdb:     rdb,
		conf:    conf,
		log:     logger,
		captcha: NewCaptcha(rdb, conf, logger),
	}
}

func InitSysUser(conf *config.Config) error {
	defaultAdminRole := conf.System.DefaultAdminRole
	defaultRole := conf.System.DefaultRole
	sysUser := []*model.SysUser{
		{UID: utils.UID(), Nickname: "张漳", Username: "admin", Password: encrypt.EncryptPassword("public"), Roles: &model.Roles{defaultAdminRole, defaultRole}, Role: defaultAdminRole, Status: "yes"},
		{UID: utils.UID(), Nickname: "张漳", Username: "grain", Password: encrypt.EncryptPassword("public"), Roles: &model.Roles{defaultRole}, Role: defaultRole, Status: "yes"},
	}

	count, err := query.Q.SysUser.Count()
	if err != nil {
		return err
	}
	// 有数据就默认已被初始化过,直接返回nil
	if count > 0 {
		return nil
	}
	err = query.Q.SysUser.Create(sysUser...)
	if err != nil {
		return err
	}
	return nil
}

func (s *SysUserService) Login(login *model.LoginReq, ctx *gin.Context) (string, error) {
	user, err := s.repo.Login(login)
	if err != nil {
		return "", err
	}

	ctx.Set("LogType", "login")

	if !encrypt.ComparePasswords(user.Password, login.Password) {
		s.log.Sava(s.log.OperationLog(400, "用户登录", login, ctx))
		return "", errors.New("账号或密码不正确")
	}

	if user.Status == "no" {
		s.log.Sava(s.log.OperationLog(400, "用户登录", login, ctx))
		return "", errors.New("账号已被冻结,无法正常登录")
	}

	jwt := utils.Jwt{}
	token, err := jwt.GenerateToken(user.UID, user.Role, s.conf.JWT.SecretKey, s.conf.JWT.ExpirationSeconds)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "用户登录", login, ctx))
		return "", err
	}
	s.log.Sava(s.log.OperationLog(200, "用户登录", login, ctx))
	return token, err
}

func (s *SysUserService) GetLoginUserInfo(ctx *gin.Context) (*model.SysUser, error) {
	info, err := s.repo.GetSysUserByUId(ctx.GetString("uid"))
	if err != nil {
		return nil, err
	}
	for _, s2 := range *info.Roles {
		find, err := query.SysRole.Where(query.SysRole.Role.Eq(s2)).First()
		if err != nil {
			continue
		}
		info.RoleStr = append(info.RoleStr, model.RoleStr{
			Value: s2,
			Label: find.RoleName,
		})
	}
	info.Password = ""
	info.Role = ctx.GetString("role")
	if info.Avatar == "" {
		info.Avatar = ""
	} else {
		info.Avatar = s.conf.Server.FileDomain + "/" + info.Avatar
	}
	return info, err
}

func (s *SysUserService) LogOut(ctx *gin.Context) error {
	s.rdb.Set(
		fmt.Sprintf("%s%s",
			consts.TokenBlack,
			ctx.GetString("token")),
		100, time.Duration(ctx.GetInt64("expTokenAt")))
	return nil
}

func (s *SysUserService) CreateSysUser(sysUser *model.SysUser, ctx *gin.Context) error {
	sysUser.UID = utils.UID()
	sysUser.ID = 0
	sysUser.Password = encrypt.EncryptPassword(sysUser.Password)
	err := s.repo.CreateSysUser(sysUser)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "创建系统用户", sysUser, ctx))
		if strings.Contains(err.Error(), " for key") {
			return errors.New("提交的参数重复")
		}
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "创建系统用户", sysUser, ctx))
	return nil
}

func (s *SysUserService) GetSysUserById(sysUserId uint, ctx *gin.Context) (*model.SysUser, error) {
	return s.repo.GetSysUserById(sysUserId)
}

func (s *SysUserService) GetSysUserList(req *model.SysUserReq, ctx *gin.Context) ([]*model.SysUser, error) {
	list, err := s.repo.GetSysUserList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	return list, err
}

func (s *SysUserService) UpdateSysUser(sysUser *model.UpdateUserInfo, ctx *gin.Context) error {
	sysUser.UID = ctx.GetString("uid")
	err := s.repo.UpdateSysUser(sysUser)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新系统用户信息", sysUser, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新系统用户信息", sysUser, ctx))
	return nil
}

func (s *SysUserService) ModifyPassword(sysUser *model.ModifyPassword, ctx *gin.Context) error {
	sysUser.UID = ctx.GetString("uid")
	user, err := s.repo.GetSysUserByUId(sysUser.UID)
	if err != nil {
		return err
	}
	if !encrypt.ComparePasswords(user.Password, sysUser.OldPassword) {
		return errors.New("旧密码不正确")
	}

	newUserInfo := model.SysUser{
		Model:    model.Model{ID: user.ID},
		Password: encrypt.EncryptPassword(sysUser.NewPassword),
	}
	err = s.repo.EditSysUser(&newUserInfo)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "修改密码", sysUser, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "修改密码", sysUser, ctx))
	return nil
}

// ConfirmModifyEmail 确认修改邮箱
func (s *SysUserService) ConfirmModifyEmail(key string, ctx *gin.Context) error {
	newUserInfo := model.SysUser{}
	//aes 解密key
	encrypt, err := encrypt.AesDecrypt(key, []byte("b06d734d53dc73c7"), []byte("0000000000000000"))
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "确认修改邮箱", key, ctx, err.Error()))
		return err
	}
	key = fmt.Sprintf("%s:%s", "confirmEmail", encrypt)
	err = s.rdb.GetObject(key, &newUserInfo)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "确认修改邮箱", newUserInfo, ctx, err.Error()))
		return err
	}

	err = s.repo.EditSysUser(&newUserInfo)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "确认修改邮箱", newUserInfo, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "确认修改邮箱", newUserInfo, ctx))
	s.rdb.Del(key)
	return nil
}

// ModifyEmail 提交修改邮箱任务,系统会向目标邮箱发送确认修改链接,当用户点击链接成功访问后系统才更新修改的邮箱
func (s *SysUserService) ModifyEmail(email *model.ModifyEmail, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	userInfo, err := s.repo.GetSysUserByUId(uid)
	if err != nil {
		return err
	}

	rdbKey := fmt.Sprintf("%s:%s:%s", "captcha", uid, email.Captcha)
	if userInfo.Email != "" {
		if email.Captcha == "" {
			return errors.New("验证码不能为空")
		}
		captcha := s.rdb.Get(rdbKey)
		if captcha != email.Captcha {
			return errors.New("验证码不正确")
		}
	}

	newUserInfo := model.SysUser{
		Model: model.Model{ID: userInfo.ID},
		Email: email.Email,
	}

	err = s.rdb.SetObject(fmt.Sprintf("%s:%s", "confirmEmail", uid), &newUserInfo, 86400)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "提交修改邮箱待确认", email, ctx))
		return err
	}

	s.log.Sava(s.log.OperationLog(200, "提交修改邮箱待确认", email, ctx))
	s.rdb.Del(rdbKey)

	aesEncrypt, err := encrypt.AesEncrypt([]byte(uid), []byte("b06d734d53dc73c7"), []byte("0000000000000000"))
	if err != nil {
		return err
	}
	aesEncrypt = url.QueryEscape(aesEncrypt)

	html := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>邮箱确认链接</title>
  </head>
  <body>
    <p>请点击以下链接以确认修改您的邮箱地址：</p>
    <p>
      <a href="` + s.conf.Server.FileDomain + "/confirmModifyEmail?key=" + aesEncrypt + `">` + s.conf.Server.FileDomain + "/confirmModifyEmail?key=" + aesEncrypt + `</a>
    </p>
  </body>
</html>`

	err = s.captcha.CustomEmail(&model.Email{Email: email.Email}, "确认修改邮箱", html)
	if err != nil {
		return err
	}

	return nil
}

func (s *SysUserService) ModifyMobile(mobile *model.ModifyMobile, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	ip := ctx.ClientIP()
	rdbKey := fmt.Sprintf("%s:%s:%s", "captcha", ip, mobile.Captcha)
	captcha := s.rdb.Get(rdbKey)
	if captcha != mobile.Captcha {
		return errors.New("验证码不正确")
	}

	userInfo, err := s.repo.GetSysUserByUId(uid)
	if err != nil {
		return err
	}
	newUserInfo := model.SysUser{
		Model:  model.Model{ID: userInfo.ID},
		Mobile: mobile.Mobile,
	}
	err = s.repo.EditSysUser(&newUserInfo)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "修改手机号", mobile, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "修改手机号", mobile, ctx))
	s.rdb.Del(rdbKey)
	return nil
}

func (s *SysUserService) EditUserInfo(sysUser *model.SysUser, ctx *gin.Context) error {
	have := false
	role := s.conf.System.DefaultRole
	for i, s2 := range *sysUser.Roles {
		if sysUser.Role == s2 {
			have = true
			break
		}
		if i == 0 {
			role = s2
		}
	}

	if !have {
		sysUser.Role = role
	}

	if sysUser.Password != "" {
		sysUser.Password = encrypt.EncryptPassword(sysUser.Password)
	}

	err := s.repo.EditSysUser(sysUser)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新系统用户信息", sysUser, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新系统用户信息", sysUser, ctx))
	return nil
}

func (s *SysUserService) SetDefaultRole(user *model.SysUser, ctx *gin.Context) error {
	err := s.repo.SetDefaultRole(user)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "设置默认角色", user, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "设置默认角色", user, ctx))
	return nil
}

func (s *SysUserService) DeleteSysUserById(id uint, ctx *gin.Context) error {
	err := s.repo.DeleteSysUserById(id)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除用户", xjson.G{"id": id}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除用户", xjson.G{"id": id}, ctx))
	return nil
}

func (s *SysUserService) DeleteSysUserByIds(ids []uint, ctx *gin.Context) error {
	err := s.repo.DeleteSysUserByIds(ids)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除用户", xjson.G{"id": ids}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除用户", xjson.G{"id": ids}, ctx))
	return nil
}

func (s *SysUserService) UploadAvatar(avatar *model.Upload, ctx *gin.Context) error {
	err := s.repo.UploadAvatar(avatar, ctx.GetString("uid"))
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新系统用户头像", avatar, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新系统用户头像", avatar, ctx))
	return nil
}

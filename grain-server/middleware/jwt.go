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

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	utils "github.com/go-grain/go-utils"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/go-utils/response"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	model "github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/utils/const"
	"net/http"
	"time"
)

func JwtAuth(rdb redis.IRedis) gin.HandlerFunc {
	conf := config.GetConfig()
	return func(ctx *gin.Context) {
		reply := response.Response{}
		jwt := utils.Jwt{}
		tokenString := ctx.GetHeader("G-Token")
		tokenClaims, err := jwt.ParseToken(tokenString, conf.JWT.SecretKey)
		if err != nil {
			reply.WithCode(http.StatusUnauthorized).WithMessage(err.Error()).Fail(ctx)
			ctx.Abort()
			return
		}

		black, _ := rdb.GetInt(fmt.Sprintf("%s%s", consts.TokenBlack, utils.MD5(tokenString)))
		switch black {
		case 120:
			reply.WithCode(http.StatusUnauthorized).WithMessage("账号进入黑名单列表,无法在继续为您服务").Fail(ctx)
			ctx.Abort()
			return
		case 110:
			reply.WithCode(http.StatusUnauthorized).WithMessage("登录异常,请重新登录").Fail(ctx)
			ctx.Abort()
			return
		case 100:
			reply.WithCode(http.StatusUnauthorized).WithMessage("无效请求 账号已退出登录").Fail(ctx)
			ctx.Abort()
			return
		}
		//获取用户信息
		sysUser := &model.SysUser{}
		if err = rdb.GetObject(tokenClaims.Uid, sysUser); err != nil {
			sysUser, err = query.Q.SysUser.Where(query.SysUser.UID.Eq(tokenClaims.Uid)).First()
			_ = rdb.SetObject(tokenClaims.Uid, sysUser, 180)
		}

		//把用户相关信息都塞到ctx去,方便下游使用
		if err == nil {
			ctx.Set("username", sysUser.Username)
			ctx.Set("nickname", sysUser.Nickname)
			ctx.Set("email", sysUser.Email)
			ctx.Set("mobil", sysUser.Mobile)
		}
		expired := int64(tokenClaims.ExpiresAt.Time.Sub(time.Now()).Seconds())
		ctx.Set("expTokenAt", expired)
		ctx.Set("uid", tokenClaims.Uid)
		ctx.Set("role", tokenClaims.Role)
		ctx.Set("token", utils.MD5(tokenString))
		ctx.Next()
	}
}

func SwitchRole(rdb redis.IRedis) gin.HandlerFunc {
	conf := config.GetConfig()
	return func(ctx *gin.Context) {
		reply := response.Response{}
		role := ctx.Query("role")
		if role == "" {
			reply.WithCode(500).WithMessage("角色ID不能为空").Fail(ctx)
			ctx.Abort()
			return
		}
		if role == ctx.GetString("role") {
			reply.WithCode(500).WithMessage("当前已是该角色,无须切换").Fail(ctx)
			ctx.Abort()
			return
		}
		sysUser := &model.SysUser{}
		if err := rdb.GetObject(ctx.GetString("uid"), sysUser); err != nil {
			sysUser, err = query.Q.SysUser.Where(query.SysUser.UID.Eq(ctx.GetString("uid"))).First()
			_ = rdb.SetObject(ctx.GetString("uid"), sysUser, 180)
		}
		for _, s := range *sysUser.Roles {
			if s == role {
				jwt := utils.Jwt{}
				token, _ := jwt.GenerateToken(ctx.GetString("uid"), role, conf.JWT.SecretKey, conf.JWT.ExpirationSeconds)
				reply.WithMessage("切换角色成功").WithData(gin.H{"token": token}).Success(ctx)
				ctx.Abort()
				break
			}
		}
		ctx.Abort()
		return
	}
}

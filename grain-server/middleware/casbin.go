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
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/go-utils/response"
	"net/http"
)

func Casbin(enforcer *casbin.CachedEnforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reply := response.Response{}
		// 权限验证
		enforce, err := enforcer.Enforce(ctx.GetString("role"), ctx.Request.URL.Path, ctx.Request.Method)
		if err != nil {
			reply.WithCode(http.StatusInternalServerError).WithMessage(err.Error()).Fail(ctx)
			ctx.Abort()
			return
		}
		if !enforce {
			reply.WithCode(http.StatusForbidden).WithMessage("无权限").Fail(ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

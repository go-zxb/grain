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

package core

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/go-utils/response"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/data"
	"github.com/go-grain/grain/internal/repo/system/query"
	sysRouter "github.com/go-grain/grain/internal/router/system"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/middleware"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

type Grain struct {
	db       *gorm.DB
	sysLog   log.Logger
	engine   *gin.Engine
	conf     *config.Config
	rdb      redis.IRedis
	enforcer *casbin.CachedEnforcer
}

func (r *Grain) InitConf() (err error) {
	r.conf, err = config.InitConfig()
	if err != nil {
		return
	}

	os.Mkdir(".tmp/", 0o664)
	file, err := os.OpenFile(".tmp/grain.log", os.O_CREATE|os.O_RDWR, 0o664)
	if err != nil {
		return err
	}
	r.sysLog = log.With(log.NewStdLogger(file),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)

	r.db, err = data.InitDB(*r.conf)
	if err != nil {
		return
	}

	r.rdb, err = data.InitRedis()
	if err != nil {
		return
	}

	r.enforcer = service.NewCasbin(r.db)

	return
}

func (r *Grain) InitRouter() {
	r.engine = gin.Default()
	gin.SetMode(r.conf.Gin.Model)
	r.engine.Use(middleware.Cors())

	routerGroup := r.engine.Group("api/v1")
	r.engine.NoRoute(func(ctx *gin.Context) {
		reply := response.Response{}
		reply.WithCode(404).WithMessage("请求路径不正确").Fail(ctx)
	})

	sysRouter.InitRouterSwag(routerGroup)
	sysRouter.NewCaptchaRouter(routerGroup, r.rdb, r.conf, r.sysLog).InitRouters()
	sysRouter.NewSysLogRouter(routerGroup, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters()
	sysRouter.NewApiRouter(routerGroup, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters().InitApi()
	sysRouter.NewOrganizeRouter(routerGroup, r.db, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters()
	sysRouter.NewMenuRouter(routerGroup, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters().InitMenu()
	sysRouter.NewRoleRouter(routerGroup, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters().InitRole()
	sysRouter.NewUploadRouter(routerGroup, r.engine, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters()
	sysRouter.NewCasbinRouter(routerGroup, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters().InitCasbin()
	sysRouter.NewCodeAssistantRouter(routerGroup, r.db, r.rdb, r.conf, r.sysLog, r.enforcer).InitRouters()
	sysRouter.NewSysUserRouter(r.engine, routerGroup, r.rdb, r.conf, r.enforcer, r.sysLog).InitRouters().InitUser()
}

func (r *Grain) RunGin() {
	if err := r.engine.Run(r.conf.Gin.Host); err != nil {
		panic(err)
	}
}

func (r *Grain) InitGenQuery() {
	query.SetDefault(r.db)
}

func (r *Grain) LoadPolicy() {
	if err := r.enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}

func Run() {
	grain := Grain{}

	if err := grain.InitConf(); err != nil {
		panic(err)
	}

	grain.InitGenQuery()

	grain.InitRouter()

	grain.LoadPolicy()

	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("swag文档地址:http://127.0.0.1:8080/api/v1/swagger/index.html")
	}()

	grain.RunGin()
}

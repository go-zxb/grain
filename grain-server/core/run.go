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

type IInit interface {
	init(grain *Grain) error
}

type Grain struct {
	db       *gorm.DB
	sysLog   log.Logger
	engine   *gin.Engine
	conf     *config.Config
	rdb      redis.IRedis
	enforcer *casbin.CachedEnforcer
}

type InitConf struct{}

func (InitConf) init(grain *Grain) (err error) {
	grain.conf, err = config.InitConfig()
	if err != nil {
		return
	}

	os.Mkdir(".tmp/", 0o664)
	file, err := os.OpenFile(".tmp/grain.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o664)
	if err != nil {
		return err
	}

	grain.sysLog = log.With(log.NewStdLogger(file),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)

	grain.db, err = data.InitDB(*grain.conf)
	if err != nil {
		return
	}

	grain.rdb, err = data.InitRedis()
	if err != nil {
		return
	}

	grain.enforcer = service.NewCasbin(grain.db)

	return
}

type InitRouter struct{}

func (InitRouter) init(grain *Grain) (err error) {
	grain.engine = gin.Default()
	gin.SetMode(grain.conf.Gin.Model)
	grain.engine.Use(middleware.Cors())

	routerGroup := grain.engine.Group("api/v1")
	grain.engine.NoRoute(func(ctx *gin.Context) {
		reply := response.Response{}
		reply.WithCode(404).WithMessage("请求路径不正确").Fail(ctx)
	})

	sysRouter.InitRouterSwag(routerGroup)
	sysRouter.NewCaptchaRouter(routerGroup, grain.rdb, grain.conf, grain.sysLog).InitRouters()
	sysRouter.NewSysLogRouter(routerGroup, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters()
	sysRouter.NewApiRouter(routerGroup, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters().InitApi()
	sysRouter.NewOrganizeRouter(routerGroup, grain.db, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters()
	sysRouter.NewMenuRouter(routerGroup, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters().InitMenu()
	sysRouter.NewRoleRouter(routerGroup, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters().InitRole()
	sysRouter.NewUploadRouter(routerGroup, grain.engine, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters()
	sysRouter.NewCasbinRouter(routerGroup, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters().InitCasbin()
	sysRouter.NewCodeAssistantRouter(routerGroup, grain.db, grain.rdb, grain.conf, grain.sysLog, grain.enforcer).InitRouters()
	sysRouter.NewSysUserRouter(grain.engine, routerGroup, grain.rdb, grain.conf, grain.enforcer, grain.sysLog).InitRouters().InitUser()
	return nil
}

type RunGin struct{}

func (RunGin) init(grain *Grain) (err error) {
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("swag文档地址:http://127.0.0.1:8080/api/v1/swagger/index.html")
	}()
	if err := grain.engine.Run(grain.conf.Gin.Host); err != nil {
		return err
	}
	return nil
}

type InitGenQuery struct{}

func (InitGenQuery) init(grain *Grain) (err error) {
	query.SetDefault(grain.db)
	return nil
}

type LoadPolicy struct{}

func (LoadPolicy) init(grain *Grain) (err error) {
	if err := grain.enforcer.LoadPolicy(); err != nil {
		return err
	}
	return nil
}

func (Grain) Do(grain *Grain, init []IInit) {
	for _, iInit := range init {
		err := iInit.init(grain)
		if err != nil {
			panic(err)
		}
	}
}

func Run() {

	grain := Grain{}
	init := []IInit{
		&InitConf{},
		&InitGenQuery{},
		&LoadPolicy{},
		&InitRouter{},
		&RunGin{},
	}

	grain.Do(&grain, init)

}

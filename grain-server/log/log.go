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

package log

import (
	"context"
	"github.com/gin-gonic/gin"
	xjson "github.com/go-grain/go-utils/json"
	model "github.com/go-grain/grain/model/system"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

// Logger 目前日志保存在MongoDB
type Logger struct {
	Mongo *mongo.Collection
}

func NewLog(mongo *mongo.Collection) (*Logger, error) {
	return &Logger{Mongo: mongo}, nil
}

func (l *Logger) SetMongo(mongo *mongo.Collection) {
	l.Mongo = mongo
}

func (l *Logger) Sava(data any) {
	// 获取对象的反射值
	v := reflect.ValueOf(data)

	// 获取方法的反射值
	method := v.MethodByName("BeforeSave")

	// 检查方法是否存在
	if method.IsValid() {
		// 调用方法
		method.Call(nil)
	}
	_, err := l.Mongo.InsertOne(context.TODO(), data)
	if err != nil {
		return
	}
}

func (l *Logger) OperationLog(code int, name string, data any, ctx *gin.Context, resData ...any) *model.SysLog {
	return &model.SysLog{
		Name:       name,
		UID:        ctx.GetString("uid"),
		Role:       ctx.GetString("role"),
		LogType:    ctx.GetString("LogType"),
		Username:   ctx.GetString("username"),
		Nickname:   ctx.GetString("nickname"),
		Method:     ctx.Request.Method,
		Path:       ctx.Request.URL.Path,
		ReqData:    l.Marshal(data),
		ResCode:    code,
		ResData:    resData,
		ClientIP:   ctx.ClientIP(),
		RequestAt:  time.Now(),
		ResponseAt: time.Now(),
		StatusCode: ctx.Writer.Status(),
	}
}

func (l *Logger) Marshal(data any) string {
	marshal := xjson.Marshal(data)
	if marshal == nil {
		return "json Marshal fail"
	}
	return string(marshal)
}

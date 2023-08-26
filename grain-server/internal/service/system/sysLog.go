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
	"github.com/gin-gonic/gin"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/log"
	model "github.com/go-grain/grain/model/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISysLogRepo interface {
	CreateSysLog(operationLog *model.SysLog) error
	GetSysLogList(req *model.SysLogReq) ([]*model.SysLog, error)
	DeleteSysLogById(id primitive.ObjectID, uid string) error
	DeleteSysLogByIds(ids []primitive.ObjectID, uid string) error
}

type SysLogService struct {
	repo ISysLogRepo
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Logger
}

func NewSysLogService(repo ISysLogRepo, rdb redis.IRedis, conf *config.Config, logger *log.Logger) *SysLogService {
	return &SysLogService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  logger,
	}
}

func (s *SysLogService) GetSysLogList(req *model.SysLogReq, ctx *gin.Context) ([]*model.SysLog, error) {
	list, err := s.repo.GetSysLogList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	return list, err
}

func (s *SysLogService) DeleteSysLogById(operationLogId string, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	objectID, err := primitive.ObjectIDFromHex(operationLogId)
	if err != nil {
		return err
	}
	err = s.repo.DeleteSysLogById(objectID, uid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除日志", operationLogId, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除日志", operationLogId, ctx))
	return nil
}

func (s *SysLogService) DeleteSysLogByIds(operationLogIds []primitive.ObjectID, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	err := s.repo.DeleteSysLogByIds(operationLogIds, uid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "批量删除日志", operationLogIds, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "批量删除日志", operationLogIds, ctx))
	return nil
}

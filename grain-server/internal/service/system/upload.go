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
)

type IUploadRepo interface {
	CreateUpload(upload *model.Upload) error
	GetUploadList(req *model.UploadReq) ([]*model.Upload, error)
	DeleteUploadById(uploadId uint, uid string) error
	DeleteUploadByIds(uploadIds []uint, uid string) error
}

type UploadService struct {
	repo IUploadRepo
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Logger
}

func NewUploadService(repo IUploadRepo, rdb redis.IRedis, conf *config.Config, logger *log.Logger) *UploadService {
	return &UploadService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  logger,
	}
}

func (s *UploadService) CreateUpload(upload *model.Upload, ctx *gin.Context) error {
	upload.UID = ctx.GetString("uid")
	err := s.repo.CreateUpload(upload)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "上传文件", upload, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "上传文件", upload, ctx))
	return nil
}

func (s *UploadService) GetUploadList(req *model.UploadReq, ctx *gin.Context) ([]*model.Upload, error) {
	list, err := s.repo.GetUploadList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	for _, i2 := range list {
		i2.FileUrl = s.conf.Server.FileDomain + "/" + i2.FileUrl
	}
	return list, err
}

func (s *UploadService) DeleteUploadById(uploadId uint, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	err := s.repo.DeleteUploadById(uploadId, uid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除上传文件", uploadId, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除上传文件", uploadId, ctx))
	return nil
}

func (s *UploadService) DeleteUploadByIds(uploadIds []uint, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	err := s.repo.DeleteUploadByIds(uploadIds, uid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除上传文件", uploadIds, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除上传文件", uploadIds, ctx))
	return nil
}

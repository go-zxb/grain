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
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/log"
	model "github.com/go-grain/grain/model/system"
	redisx "github.com/go-grain/grain/pkg/redis"
)

type IUploadRepo interface {
	CreateUpload(upload *model.Upload) error
	GetUploadList(req *model.UploadReq) ([]*model.Upload, error)
	DeleteUploadById(uploadId uint, uid string) error
	DeleteUploadByIds(uploadIds []uint, uid string) error
}

type UploadService struct {
	repo IUploadRepo
	rdb  redisx.IRedis
	conf *config.Config
	log  *log.Helper
}

func NewUploadService(repo IUploadRepo, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *UploadService {
	return &UploadService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  log.NewHelper(logger),
	}
}

func (s *UploadService) CreateUpload(upload *model.Upload, ctx *gin.Context) error {
	upload.UID = ctx.GetString("uid")
	if err := s.repo.CreateUpload(upload); err != nil {
		s.log.Errorw("errMsg", "上传文件", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "上传文件")
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
	if err := s.repo.DeleteUploadById(uploadId, uid); err != nil {
		s.log.Errorw("errMsg", "删除上传文件", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "删除上传文件")
	return nil
}

func (s *UploadService) DeleteUploadByIds(uploadIds []uint, ctx *gin.Context) error {
	uid := ctx.GetString("uid")
	if err := s.repo.DeleteUploadByIds(uploadIds, uid); err != nil {
		s.log.Errorw("errMsg", "删除上传文件", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "删除上传文件")
	return nil
}

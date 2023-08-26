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

package upload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	utils "github.com/go-grain/go-utils"
	model "github.com/go-grain/grain/model/system"
	"io"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(ctx *gin.Context, classify string) (*model.Upload, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	now := time.Now()
	y := now.Year()
	m := int(now.Month())
	d := now.Day()
	path := fmt.Sprintf("uploads/%s/%d/%d-%d/", classify, y, m, d)
	filename := ""

	e := utils.Ext(file.Filename)
	if e == "" {
		filename = path + utils.RandomCharset(5)
	}
	filename = path + utils.RandomCharset(5) + "." + e

	exist := utils.PathIsNotExist(path)
	if exist {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	// 检测重名
	if !utils.FileIsNotExist(filename) {
		for { //慢慢找一个和本地文件名不重复的随机字符串做文件名
			filename = path + utils.RandomCharset(5) + "." + e
			if !utils.FileIsNotExist(filename) {
				continue
			} else {
				break
			}
		}
	}

	dst, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(file.Filename)
	fileType := ""
	switch ext {
	case ".jpg", ".jpeg", ".png":
		fileType = "图片 "
	case ".zip", ".gz", "7z":
		fileType = "压缩包"
	case ".mp4", ".mov", "flv", "mp3":
		fileType = "音视频"
	case ".go":
		fileType = "Go文件"
	default:
		fileType = "NA"
	}

	upload := &model.Upload{
		FileName: file.Filename,
		FileUrl:  filename,
		FileType: fileType,
	}
	return upload, nil
}

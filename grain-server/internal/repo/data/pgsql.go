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

package data

import (
	"github.com/go-grain/grain/config"
	"github.com/go-pay/gopay/pkg/xlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitPgsql(conf config.DataBase) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,               // 慢 SQL 阈值
			LogLevel:                  config.GetConfig().DataBase.LogLevel, // 日志级别
			IgnoreRecordNotFoundError: true,                                 // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,                                 // 禁用彩色打印
		},
	)

	pgsqlConfig := postgres.Config{
		DriverName: "postgres",
		DSN:        conf.Pgsql.Source,
	}
	gormDB, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	sqlDB, _ := gormDB.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxIdleTime(time.Second * 5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	xlog.Info("初始化Pgsql成功")

	db = &DB{DB: gormDB}
	err = db.autoMigrate()
	if err != nil {
		xlog.Info("Pgsql AutoMigrate error", err.Error())
		return nil, err
	}
	return gormDB, err
}

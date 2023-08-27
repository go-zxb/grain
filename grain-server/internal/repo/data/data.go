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
	"errors"
	"github.com/go-grain/grain/config"
	sysModel "github.com/go-grain/grain/model/system"
	"gorm.io/gorm"
)

const (
	// dbMySQL Gorm Drivers mysql || postgres || sqlite || sqlserver
	dbMySQL    string = "mysql"
	dbPostgres string = "postgres"
	dbTidb     string = "tidb"
)

var db *DB

type DB struct {
	DB *gorm.DB
}

func InitDB(conf config.Config) (*gorm.DB, error) {
	switch conf.DataBase.Driver {
	case dbMySQL:
		mysql, err := InitMysql(conf)
		if err != nil {
			return nil, err
		}
		return mysql, err
	case dbPostgres:
		pgsql, err := InitPgsql(conf.DataBase)
		if err != nil {
			return nil, err
		}
		return pgsql, err
	case dbTidb:
		tidb, err := InitTiDB(conf.DataBase)
		if err != nil {
			return nil, err
		}
		return tidb, err
	default:
		return nil, errors.New("数据库配置有问题")
	}
}

func NewDB() *DB {
	return db
}

func (db *DB) autoMigrate() error {
	err := db.DB.AutoMigrate(
		sysModel.SysRole{},
		sysModel.SysUser{},
		sysModel.SysApi{},
		sysModel.SysMenu{},
		sysModel.Upload{},
		sysModel.Project{},
		sysModel.Models{},
		sysModel.Fields{},
	)
	if err != nil {
		return err
	}
	return nil
}

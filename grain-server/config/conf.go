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

package config

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	xjson "github.com/go-grain/go-utils/json"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

//go:embed config.yaml
var file string

//go:embed email.key
var emailData string

var xViper *viper.Viper
var config *Config

type System struct {
	SiteName         string `mapstructure:"site_name" json:"site_name" yaml:"site_name"`
	CaptchaLength    int    `mapstructure:"captcha_length" json:"captcha_length" yaml:"captcha_length"`
	DefaultRole      string `mapstructure:"default_role" json:"default_role" yaml:"default_role"`
	DefaultAdminRole string `mapstructure:"default_admin_role" json:"default_admin_role" yaml:"default_admin_role"`
}

type SysEmail struct {
	EmailHost     string `mapstructure:"email_host" json:"email_host" yaml:"email_host"`
	EmailPort     int    `mapstructure:"email_port" json:"email_port" yaml:"email_port"`
	EmailUsername string `mapstructure:"email_username" json:"email_username" yaml:"email_username"`
	EmailPassword string `mapstructure:"email_password" json:"email_password" yaml:"email_password"`
}

type Log struct {
	Level     zapcore.Level `mapstructure:"level" json:"level" yaml:"level"`
	LogPath   string        `mapstructure:"log_path" json:"log_path" yaml:"log_path"`
	SplitSize int           `mapstructure:"split_size" json:"split_size" yaml:"split_size"`
}

type Server struct {
	FileDomain string `mapstructure:"file_domain" json:"file_domain" yaml:"file_domain"`
}

type Gin struct {
	Host  string `mapstructure:"host" json:"host" yaml:"host"`
	Model string `mapstructure:"model" json:"model" yaml:"model"`
}

type DataBase struct {
	Driver   string          `mapstructure:"driver" json:"driver" yaml:"driver"`
	LogLevel logger.LogLevel `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	MySql    struct {
		Source string `mapstructure:"source" json:"source" yaml:"source"`
	} `mapstructure:"mysql" json:"mysql" yaml:"mysql"`

	TiDB struct {
		Source string `mapstructure:"source" json:"source" yaml:"source"`
	} `mapstructure:"ti_db" json:"ti_db" yaml:"ti_db"`

	Pgsql struct {
		Source string `mapstructure:"source" json:"source" yaml:"source"`
	} `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`

	Sqlite struct {
		Source string `mapstructure:"source" json:"source" yaml:"source"`
	} `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`

	Redis struct {
		UserName     string        `mapstructure:"user_name" json:"user_name" yaml:"user_name"`
		Password     string        `mapstructure:"password" json:"password" yaml:"password"`
		Addr         string        `mapstructure:"addr" json:"addr" yaml:"addr"`
		DB           int           `mapstructure:"db" json:"db" yaml:"db"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
	} `yaml:"redis"`

	Mongo struct {
		URL string `mapstructure:"url" json:"url" yaml:"url"`
	} `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
}

type JWT struct {
	SecretKey         string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	ExpirationSeconds int64  `mapstructure:"expiration_seconds" json:"expiration_seconds" yaml:"expiration_seconds"`
	Issuer            string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}

type Config struct {
	Gin      Gin      `mapstructure:"gin" json:"gin" yaml:"gin"`
	System   System   `mapstructure:"system" json:"system" yaml:"system"`
	SysEmail SysEmail `mapstructure:"email" json:"email" yaml:"email"`
	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
	Server   Server   `mapstructure:"server" json:"server" yaml:"server"`
	DataBase DataBase `mapstructure:"database" json:"database" yaml:"database"`
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

func GetConfig() *Config {
	return config
}

func SetKey(key string, value interface{}) error {
	xViper.Set(key, value)
	err := xViper.WriteConfig()
	if err != nil {
		fmt.Printf("Failed to write config file: %s", err)
		return err
	}
	return nil
}

func InitConfig() (*Config, error) {
	// 如果没有配置文件 则从导出一份默认配置文件到本地
	_, err := os.Stat("config/config.yaml")
	if os.IsNotExist(err) {
		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			return nil, errors.New("配置文件初始化失败")
		}
		err = os.WriteFile("config/config.yaml", []byte(file), os.ModePerm)
		if err != nil {
			return nil, errors.New("配置文件初始化失败")
		}
	}

	path := flag.String("config", "config/config.yaml", "添加配置文件路径")
	showHelp := flag.Bool("h", false, "显示帮助")

	flag.Parse()

	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var conf *viper.Viper
	conf = viper.New()
	conf.SetConfigFile(*path)
	err = conf.ReadInConfig()
	if err != nil {
		panic(any(err.Error()))
		return nil, err
	}
	v := Config{}
	err = conf.Unmarshal(&v)
	if err != nil {
		panic(any(err.Error()))
		return nil, err
	}
	conf.WatchConfig()

	conf.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已更改: ", e.Name)
		if err = conf.Unmarshal(&v); err != nil {
			fmt.Println(err)
		}
	})
	if err = conf.Unmarshal(&v); err != nil {
		fmt.Println(err)
	}

	//为了push到github不暴露邮箱配置信息,放在别的地方解析过来,作为菜鸟的我,只能使用这种简单粗暴的方式实现了
	if v.Gin.Model == "debug" {
		_ = xjson.Unmarshal([]byte(emailData), &v.SysEmail)
	}

	xViper = conf
	config = &v
	return config, nil
}

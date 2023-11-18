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
	"fmt"
	"github.com/gin-gonic/gin"
	xjson "github.com/go-grain/go-utils/json"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/internal/repo/system/query"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/model/system"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type IApiRepo interface {
	CreateApi(api *model.SysApi) error
	GetApiList(req *model.SysApiReq) ([]*model.SysApi, error)
	GetAllApi() ([]*model.SysApi, error)
	UpdateApi(api *model.SysApi) error
	DeleteApiByIds(ids []uint) error
	DeleteApiById(id uint) error
	AuthApiList(role string) (list []*model.CasbinRule, err error)
}

type ApiService struct {
	repo IApiRepo
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Logger
}

func NewApiService(repo IApiRepo, rdb redis.IRedis, conf *config.Config, logger *log.Logger) *ApiService {
	return &ApiService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  logger,
	}
}

func InitApi() error {
	apis := []*model.SysApi{

		// 系统角色
		{Path: "/api/v1/sysRole", Description: "编辑角色", ApiGroup: "系统角色", Method: "PUT"},
		{Path: "/api/v1/sysRole", Description: "创建角色", ApiGroup: "系统角色", Method: "POST"},
		{Path: "/api/v1/sysRole/list", Description: "获取角色列表", ApiGroup: "系统角色", Method: "GET"},
		{Path: "/api/v1/sysRole/deleteSysRoleByIds", Description: "批量删除角色", ApiGroup: "系统角色", Method: "DELETE"},

		// casbin
		{Path: "/api/v1/casbin", Description: "更新角色权限", ApiGroup: "系统权限", Method: "PUT"},
		{Path: "/api/v1/casbin/list", Description: "获取已授权的Api列表", ApiGroup: "系统权限", Method: "GET"},

		// 系统Api
		{Path: "/api/v1/sysApi", Description: "创建Api", ApiGroup: "系统Api", Method: "POST"},
		{Path: "/api/v1/sysApi", Description: "编辑Api", ApiGroup: "系统Api", Method: "PUT"},
		{Path: "/api/v1/sysApi", Description: "删除Api", ApiGroup: "系统Api", Method: "DELETE"},
		{Path: "/api/v1/sysApi/list", Description: "获取Api列表", ApiGroup: "系统Api", Method: "GET"},
		{Path: "/api/v1/sysApi/deleteSysApiByIds", Description: "批量删除Api", ApiGroup: "系统Api", Method: "DELETE"},
		{Path: "/api/v1/sysApi/apiGroups", Description: "获取Api分组列表", ApiGroup: "系统Api", Method: "GET"},
		{Path: "/api/v1/sysApi/apiAndPermissions", Description: "获取已授权的Api列表", ApiGroup: "系统Api", Method: "GET"},

		// 系统用户
		{Path: "/api/v1/sysUser", Description: "删除系统用户", ApiGroup: "系统用户", Method: "DELETE"},
		{Path: "/api/v1/sysUser/info", Description: "获取用户信息", ApiGroup: "系统用户", Method: "GET"},
		{Path: "/api/v1/sysUser/editUserInfo", Description: "编辑系统用户", ApiGroup: "系统用户", Method: "PUT"},
		{Path: "/api/v1/sysUser/update", Description: "更新个人信息", ApiGroup: "系统用户", Method: "PUT"},
		{Path: "/api/v1/sysUser/create", Description: "创建系统用户", ApiGroup: "系统用户", Method: "POST"},
		{Path: "/api/v1/sysUser/list", Description: "获取系统用户列表", ApiGroup: "系统用户", Method: "GET"},
		{Path: "/api/v1/sysUser/deleteSysUserByIds", Description: "批量删除系统用户", ApiGroup: "系统用户", Method: "DELETE"},
		{Path: "/api/v1/sysUser/avatar", Description: "更新系统用户头像", ApiGroup: "系统用户", Method: "POST"},
		{Path: "/api/v1/sysUser/setDefaultRole", Description: "设置默认角色", ApiGroup: "系统用户", Method: "PUT"},

		//系统菜单
		{Path: "/api/v1/sysMenu", Description: "编辑菜单", ApiGroup: "系统菜单", Method: "PUT"},
		{Path: "/api/v1/sysMenu", Description: "创建菜单", ApiGroup: "系统菜单", Method: "POST"},
		{Path: "/api/v1/sysMenu", Description: "删除菜单", ApiGroup: "系统菜单", Method: "DELETE"},
		{Path: "/api/v1/sysMenu/list", Description: "获取菜单列表", ApiGroup: "系统菜单", Method: "GET"},
		{Path: "/api/v1/sysMenu/deleteSysMenuByIds", Description: "批量删除菜单", ApiGroup: "系统菜单", Method: "DELETE"},
		{Path: "/api/v1/sysMenu/userMenu", Description: "获取动态菜单", ApiGroup: "系统用户", Method: "GET"},

		// 代码助手
		{Path: "/api/v1/codeAssistant/projects", Description: "添加或更新项目信息", ApiGroup: "代码助手", Method: "POST"},
		{Path: "/api/v1/codeAssistant/models", Description: "添加或更新模型信息", ApiGroup: "代码助手", Method: "POST"},
		{Path: "/api/v1/codeAssistant/fields", Description: "添加或更新字段信息", ApiGroup: "代码助手", Method: "POST"},
		{Path: "/api/v1/codeAssistant/projects", Description: "删除项目", ApiGroup: "代码助手", Method: "DELETE"},
		{Path: "/api/v1/codeAssistant/models", Description: "删除模型", ApiGroup: "代码助手", Method: "DELETE"},
		{Path: "/api/v1/codeAssistant/fields", Description: "删除字段", ApiGroup: "代码助手", Method: "DELETE"},
		{Path: "/api/v1/codeAssistant/projects/list", Description: "获取项目列表", ApiGroup: "代码助手", Method: "GET"},
		{Path: "/api/v1/codeAssistant/models/list", Description: "获取模型列表", ApiGroup: "代码助手", Method: "GET"},
		{Path: "/api/v1/codeAssistant/fields/list", Description: "获取字段列表", ApiGroup: "代码助手", Method: "GET"},
		{Path: "/api/v1/codeAssistant/viewCode", Description: "预览代码", ApiGroup: "代码助手", Method: "GET"},

		// 系统日志
		{Path: "/api/v1/sysLog/list", Description: "获取系统操作日志列表", ApiGroup: "系统菜单", Method: "GET"},
		{Path: "/api/v1/sysLog", Description: "删除系统操作日志", ApiGroup: "系统菜单", Method: "DELETE"},
	}
	q := query.Q.SysApi

	count, err := q.Count()
	if err != nil {
		return err
	}

	// 数据相等则退出 否则就更新
	if count != 0 {
		return nil
	}

	for _, api := range apis {
		_, err = q.Where(q.Path.Eq(api.Path)).Where(q.Method.Eq(api.Method)).First()
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		err = q.Create(api)
		if err != nil {
			continue
		}
	}

	return nil
}

// CreateApi 创建api
func (s *ApiService) CreateApi(api *model.SysApi, ctx *gin.Context) error {
	//获取分组信息 为了偷懒就这样处理了分组问题
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(api.Description, -1)
	if len(matches) >= 1 {
		api.ApiGroup = matches[0][1]
		api.Description = strings.ReplaceAll(api.Description, fmt.Sprintf("[%s]", api.ApiGroup), "")
	}
	err := s.repo.CreateApi(api)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "创建Api", api, ctx))
		if strings.Contains(err.Error(), "duplicated key not allowed") {
			return errors.New("提交的参数重复")
		}
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "创建Api", api, ctx))
	return nil
}

func (s *ApiService) GetApiList(req *model.SysApiReq) ([]*model.SysApi, error) {
	list, err := s.repo.GetApiList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	return list, nil
}

func (s *ApiService) GetApiGroup() (any, error) {
	list, err := s.repo.GetAllApi()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}

	type G struct {
		Group string `json:"group"`
	}
	apis := make(map[string]string)
	for _, i2 := range list {
		apis[i2.ApiGroup] = i2.ApiGroup
	}

	var res []G
	for _, s3 := range apis {
		res = append(res, G{Group: s3})
	}

	return res, err
}

func (s *ApiService) GetApiAndPermissions(role string) (any, error) {
	list, err := s.repo.GetAllApi()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	apis := make(map[string]*model.ApiGroup)
	for _, api := range list {
		//判断相同的分组是否已存在,存在就直接添加到该组里面
		_, ok := apis[api.ApiGroup]
		if ok {
			Children := model.ApiGroup{
				ApiGroup:    api.ApiGroup,
				ID:          api.ID,
				Description: api.Description,
				Method:      api.Method,
				Path:        api.Path,
			}
			apis[api.ApiGroup].Children = append(apis[api.ApiGroup].Children, Children)
		} else {
			Children := []model.ApiGroup{
				{
					ApiGroup:    api.ApiGroup,
					ID:          api.ID,
					Description: api.Description,
					Method:      api.Method,
					Path:        api.Path,
				},
			}
			apis[api.ApiGroup] = &model.ApiGroup{
				ApiGroup:    api.ApiGroup,
				ID:          api.ID + 10000,
				Description: api.ApiGroup,
				Children:    Children,
				Path:        api.Path,
			}
		}

	}
	//所有api 放到一个切片
	var apiSlice = make([]any, 0)
	for _, g := range apis {
		apiSlice = append(apiSlice, g)
	}
	//获取已授权的api
	authApi, err := s.repo.AuthApiList(role)
	if err != nil {
		return nil, err
	}

	var authID []uint
	for _, i2 := range list {
		for _, rule := range authApi {
			// 查找已授权的Api ID
			if i2.Path == rule.V1 && i2.Method == rule.V2 {
				authID = append(authID, i2.ID)
			}
		}
	}

	res := gin.H{"apiList": apiSlice, "authApi": authID}
	return res, nil
}

func (s *ApiService) UpdateApi(api *model.SysApi, ctx *gin.Context) error {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(api.Description, -1)
	if len(matches) >= 1 {
		api.ApiGroup = matches[0][1]
		api.Description = strings.ReplaceAll(api.Description, fmt.Sprintf("[%s]", api.ApiGroup), "")
	}
	err := s.repo.UpdateApi(api)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新Api", api, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新Api", api, ctx))
	return nil
}

func (s *ApiService) DeleteApiById(id uint, ctx *gin.Context) error {
	err := s.repo.DeleteApiById(id)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除Api", xjson.G{"id": id}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除Api", xjson.G{"id": id}, ctx))
	return nil

}

func (s *ApiService) DeleteApiByIds(ids []uint, ctx *gin.Context) error {
	err := s.repo.DeleteApiByIds(ids)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除Api", xjson.G{"id": ids}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除Api", xjson.G{"id": ids}, ctx))
	return nil

}

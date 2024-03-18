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
	"github.com/Xuanwo/gg"
	"github.com/gin-gonic/gin"
	utils "github.com/go-grain/go-utils"
	"github.com/go-grain/go-utils/goutil"
	xjson "github.com/go-grain/go-utils/json"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/log"
	model "github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/stencil"
	"github.com/go-pay/gopay/pkg/xlog"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"text/template"
)

const (
	MongoDB = "mongo"
	mysql   = "mysql"
)

type ICodeAssistantRepo interface {
	CreateProject(p *model.Project) error
	UpdateProject(p *model.Project) error
	DeleteProjectById(pid uint) error
	CreateModel(m *model.Models) error
	UpdateModel(m *model.Models) error
	DeleteModelById(mid uint) error
	CreateField(f *model.Fields) error
	UpdateField(f *model.Fields) error
	DeleteFieldById(fid uint) error
	GetProjectList() ([]*model.Project, error)
	GetProject(pid uint) (*model.Project, error)
	GetModel(mid uint) (*model.Models, error)
	GetModels(parentId uint) ([]*model.Models, error)
	GetFields(parentId uint) ([]*model.Fields, error)
}

type CodeAssistantService struct {
	repo ICodeAssistantRepo
	rdb  redis.IRedis
	conf *config.Config
	log  *log.Logger
}

func NewCodeAssistantService(repo ICodeAssistantRepo, rdb redis.IRedis, conf *config.Config, logger *log.Logger) *CodeAssistantService {
	return &CodeAssistantService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  logger,
	}
}

func (s *CodeAssistantService) CreateProject(p *model.Project, ctx *gin.Context) error {
	p.Model = model.Model{}
	err := s.repo.CreateProject(p)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "创建项目", p, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "创建项目", p, ctx))
	return nil
}

func (s *CodeAssistantService) UpdateProject(p *model.Project, ctx *gin.Context) error {
	err := s.repo.UpdateProject(p)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新项目", p, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新项目", p, ctx))
	return nil
}

func (s *CodeAssistantService) CreateModel(m *model.Models, ctx *gin.Context) error {
	err := s.repo.CreateModel(m)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "创建模块", m, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "创建模块", m, ctx))
	return nil
}

func (s *CodeAssistantService) UpdateModel(m *model.Models, ctx *gin.Context) error {
	err := s.repo.UpdateModel(m)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新模块", m, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新模块", m, ctx))
	return nil
}

func (s *CodeAssistantService) CreateField(f *model.Fields, ctx *gin.Context) error {
	err := s.repo.CreateField(f)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "创建字段", f, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "创建字段", f, ctx))
	return nil
}

func (s *CodeAssistantService) UpdateField(f *model.Fields, ctx *gin.Context) error {
	err := s.repo.UpdateField(f)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "更新字段", f, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "更新字段", f, ctx))
	return nil
}

func (s *CodeAssistantService) GetProjectList(ctx *gin.Context) ([]*model.Project, error) {
	list, err := s.repo.GetProjectList()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}

	for _, project := range list {
		projectModels, err := s.repo.GetModels(project.ID)
		if err != nil {
			return nil, err
		}
		project.Models = projectModels
		for _, models := range project.Models {
			fields, err := s.repo.GetFields(models.ID)
			if err != nil {
				return nil, err
			}
			models.Fields = fields
		}
	}
	return list, err
}

func (s *CodeAssistantService) GetModelList(parentId uint, ctx *gin.Context) ([]*model.Models, error) {
	list, err := s.repo.GetModels(parentId)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	return list, err
}

func (s *CodeAssistantService) GetFieldList(parentId uint, ctx *gin.Context) ([]*model.Fields, error) {
	list, err := s.repo.GetFields(parentId)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}

	return list, err
}

func (s *CodeAssistantService) DeleteProjectByIds(pid uint, ctx *gin.Context) error {
	err := s.repo.DeleteProjectById(pid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除项目", xjson.G{"id": pid}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除项目", xjson.G{"id": pid}, ctx))
	return nil
}

func (s *CodeAssistantService) DeleteModelById(mid uint, ctx *gin.Context) error {
	err := s.repo.DeleteModelById(mid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除模型", xjson.G{"id": mid}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除模型", xjson.G{"id": mid}, ctx))
	return nil
}

func (s *CodeAssistantService) DeleteFieldById(fid uint, ctx *gin.Context) error {
	err := s.repo.DeleteFieldById(fid)
	if err != nil {
		s.log.Sava(s.log.OperationLog(400, "删除字段", xjson.G{"id": fid}, ctx))
		return err
	}
	s.log.Sava(s.log.OperationLog(200, "删除字段", xjson.G{"id": fid}, ctx))
	return nil
}

func (s *CodeAssistantService) ViewCode(mId uint, ctx *gin.Context) (*model.ViewCode, error) {
	modelData, err := s.repo.GetModel(mId)
	if err != nil {
		return nil, err
	}
	if modelData.ID == 0 {
		return nil, errors.New("数据异常")
	}
	//获取项目信息
	p, err := s.repo.GetProject(modelData.ParentId)
	if err != nil {
		return nil, err
	}

	if p.ID == 0 {
		return nil, errors.New("获取项目信息失败,请使用简单粗暴的方式查看数据库数据是否正常")
	}
	//获取字段列表
	fields, err := s.repo.GetFields(modelData.ID)
	if err != nil {
		return nil, err
	}

	modelData.Fields = fields
	modelData.ProjectName = p.ProjectName
	modelData.WebProjectPath = ".tmp"
	modelData.ProjectPath = ".tmp"

	if s.conf.Gin.Model == "debug" {
		modelData.ProjectName += "/.tmp"
	}

	modelData.FirstLetter = utils.ToLower(modelData.StructName[0:1])
	if modelData.QueryTime == "yes" {
		//处理是否支持前端按时间范围搜索数据
		modelData.IsQueryCriteria = true
	}
	maps, err := tmplateData(*modelData)
	if err != nil {
		return nil, err
	}
	viewCode := model.ViewCode{}
	modelData.GoFieldTo()
	for i, m1 := range maps {
		m1.Type = i
		err = s.generateCode(*modelData, *m1)
		if err != nil {
			return nil, err
		}

		switch i {
		case "model":
			viewCode.Model = utils.ReadFile(m1.Filename)
		case "router":
			viewCode.Router = utils.ReadFile(m1.Filename)
		case "handler":
			viewCode.Handle = utils.ReadFile(m1.Filename)
		case "service":
			viewCode.Service = utils.ReadFile(m1.Filename)
		case "repo":
			viewCode.Repo = utils.ReadFile(m1.Filename)
		}
	}
	fields, _ = s.repo.GetFields(modelData.ID)
	modelData.Fields = fields
	s.WebGenerate(*modelData)
	viewCode.APi = utils.ReadFile(fmt.Sprintf("%s/src/api/business/%s.ts", modelData.WebProjectPath, modelData.Name))
	viewCode.Vue = utils.ReadFile(fmt.Sprintf("%s/src/views/business/%s/index.vue", modelData.WebProjectPath, modelData.Name))
	viewCode.ZhCN = utils.ReadFile(fmt.Sprintf("%s/src/views/business/%s/locale/zh-CN.ts", modelData.WebProjectPath, modelData.Name))

	fields, _ = s.repo.GetFields(modelData.ID)
	modelData.Fields = fields
	s.FlutterGenerate(*modelData)
	viewCode.FlutterModel = utils.ReadFile(fmt.Sprintf("%s/flutter/models/%s/%s.dart", modelData.WebProjectPath, modelData.Name, modelData.Name))
	viewCode.FlutterAPi = utils.ReadFile(fmt.Sprintf("%s/flutter/api/%s/%s.dart", modelData.WebProjectPath, modelData.Name, modelData.Name))

	err = GenBuild(modelData)
	if err != nil {
		return nil, err
	}
	err = InjectGenBuild(modelData)
	if err != nil {
		return nil, err
	}

	//_ = os.RemoveAll(".tmp")
	return &viewCode, nil
}

func (s *CodeAssistantService) GenerateCode(mid uint, forceBuild uint, ctx *gin.Context) error {
	modelData, err := s.repo.GetModel(mid)
	if err != nil {
		return err
	}
	if modelData.ID == 0 {
		return errors.New("数据异常")
	}
	//获取项目信息
	p, err := s.repo.GetProject(modelData.ParentId)
	if err != nil {
		return err
	}

	if p.ID == 0 {
		return errors.New("获取项目信息失败,请使用简单粗暴的方式查看数据库数据是否正常")
	}
	//获取字段列表
	fields, err := s.repo.GetFields(modelData.ID)
	if err != nil {
		return err
	}

	modelData.Fields = fields
	modelData.ProjectName = p.ProjectName
	if p.ProjectPath == "" {
		//处理空路径,空路径会直接在当前运行路径 tmp下生成,如果是使用本项目作为基础框架可填写 ./即可
		modelData.ProjectPath = ".tmp"
	} else {
		modelData.ProjectPath = p.ProjectPath
	}

	if p.WebProjectPath == "" {
		//处理空路径,空路径会直接在当前运行路径 tmp下生成,如果是使用本项目作为基础框架可填写 ../xxxweb即可
		modelData.WebProjectPath = ".tmp"
	} else {
		modelData.WebProjectPath = p.WebProjectPath
	}

	if s.conf.Gin.Model == "debug" {
		modelData.ProjectName += "/.tmp"
	}

	modelData.FirstLetter = utils.ToLower(modelData.StructName[0:1])
	if modelData.QueryTime == "yes" {
		//处理是否支持前端按时间范围搜索数据
		modelData.IsQueryCriteria = true
	}
	maps, err := tmplateData(*modelData)
	if err != nil {
		return err
	}
	modelData.GoFieldTo()
	for i, m1 := range maps {
		m1.Type = i
		err = s.generateCode(*modelData, *m1)
		if err != nil {
			return err
		}
	}
	if modelData.DatabaseName == mysql {
		//生成Gorm gen 相关代码
		if err = GenBuild(modelData); err != nil {
			return err
		}

		//生成Gorm gen 相关代码
		if err := InjectGenBuild(modelData); err != nil {
			return err
		}
	}

	s.WebGenerate(*modelData)
	return nil
}

// GenBuild .
func GenBuild(model *model.Models) error {
	var code = `g := gen.NewGenerator(gen.Config{
OutPath:"` + fmt.Sprintf("%s/%s/%s/query", model.ProjectPath, "repo", model.Name) + `",
Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
})// generate from struct in project
g.ApplyBasic(
model.` + model.StructName + `{},
)
g.Execute()`

	f := gg.NewGroup()
	f.AddLineComment("自动生成的代码 请勿编辑")
	f.AddPackage(model.Name)
	f.NewImport().
		AddPath(fmt.Sprintf("%s/model/%s", model.ProjectName, model.Name)).
		AddPath("gorm.io/gen")
	f.NewFunction(fmt.Sprintf("%sGenBuildQuery", utils.ToTitle(model.StructName))).AddBody(
		gg.String(code),
	)

	filename := fmt.Sprintf("%s/cmd/gen/%s/%s.go", model.ProjectPath, model.Name, model.Name)
	filepath := fmt.Sprintf("%s/cmd/gen/%s", model.ProjectPath, model.Name)
	if !utils.FileIsNotExist(filename) {
		_ = os.Remove(filename)
	}

	if utils.PathIsNotExist(filepath) {
		err := os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	open, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		xlog.Info(err)
		return err
	}
	defer open.Close()

	_, err = open.WriteString(f.String())
	if err != nil {
		return err
	}
	err = goutil.FmtCode(filename)
	if err != nil {
		return err
	}
	return nil
}

// InjectGenBuild 向cmd/gen/build.go文件 注入xxxGenBuild()方便生成query代码
func InjectGenBuild(model *model.Models) error {
	filename := model.ProjectPath + "/cmd/gen/build.go"
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var initAppFunc *ast.FuncDecl
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == "main" {
			initAppFunc = funcDecl
			break
		}
	}

	_ = os.Remove(filename)
	open, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if err != nil {
		xlog.Info(err)
		return err
	}
	old, _ := io.ReadAll(open)
	defer open.Close()

	f := gg.NewGroup()
	f.AddLineComment("// Package " + model.Name + " 自动生成的代码 请勿编辑")
	f.AddPackage("main")
	f.NewImport().AddPath(model.ProjectName + "/cmd/gen/" + model.Name)
	f.AddLineComment("自动生成代码后,运行前请先运行该文件生成orm相关依赖文件")
	body := fmt.Sprintf("%s.%sGenBuildQuery()", model.Name, model.StructName)
	function := f.NewFunction("main").AddBody(gg.String(body))

	if initAppFunc != nil {
		for _, stmt := range initAppFunc.Body.List {
			switch e := stmt.(type) {
			case *ast.ExprStmt:
				if valx, ok := e.X.(*ast.CallExpr); ok {
					if valxf, ok := valx.Fun.(*ast.SelectorExpr); ok {
						if v, ok := valxf.X.(*ast.Ident); ok {
							function.AddBody(fmt.Sprintf("%s.%s()", v.Name, valxf.Sel.Name))
						}
					}
					if v, ok := valx.Fun.(*ast.Ident); ok {
						function.AddBody(v.Name + "()")
					}
				}
			}
		}
	}

	if _, err = open.WriteString(f.String()); err != nil {
		_, _ = open.Write(old)
		return err
	}

	if err = goutil.FmtCode(filename); err != nil {
		return err
	}
	return nil
}

func (s *CodeAssistantService) generateCode(m model.Models, p model.CodePath) (err error) {
	dir := p.FilePath
	//检查目录是否存在
	exists := utils.PathIsNotExist(dir)

	if exists {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			xlog.Error(err)
			return err
		}
	}

	//读取模板文件
	var b []byte
	b, err = p.FS.ReadFile(p.TemplatePath)
	if err != nil {
		xlog.Error(err)
		return err
	}
	rt := string(b)

	rt = strings.ReplaceAll(rt, "{{a-model}}", m.Name)
	rt = strings.ReplaceAll(rt, "{{c-First}}", m.FirstLetter)
	rt = strings.ReplaceAll(rt, "{{s-model}}", utils.ToTitle(m.Name))

	if p.Type == "query" {
		for _, field := range m.Fields {
			switch field.Type {
			case "[]byte":
				field.Type = "byte"
			case "time.Time":
				field.Type = "time"
			}
			field.Type = utils.ToTitle(field.Type)
		}
	} else {
		m.GoFieldTo()
	}

	temp, err := template.New(utils.ToLower(p.Type)).Parse(rt)
	if err != nil {
		xlog.Error(err)
		return err
	}

	filename := p.Filename
	//检查文件是否存在
	exists = utils.FileIsNotExist(filename)

	if s.conf.Gin.Model == "debug" {
		_ = os.Remove(filename)
	} else {
		if !exists {
			return errors.New("文件已存在: " + fmt.Sprintf("%s/%s.go", p.Type, m.Name))
		}
	}

	open, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		xlog.Info(err)
		return err
	}
	defer open.Close()

	//渲染模板文件
	err = temp.Execute(open, m)
	if err != nil {
		xlog.Error(err)
		return err
	}
	//格式化代码
	_ = utils.FormatGoCode(filename)
	return nil
}

func tmplateData(m model.Models) (map[string]*model.CodePath, error) {
	//把所有相关信息都存放在一个map里面
	maps := make(map[string]*model.CodePath)
	name := utils.ToLower(m.StructName)

	if m.DatabaseName == mysql {
		maps["model"] = &model.CodePath{
			FS:           &stencil.MysqlTemplateFS,
			TemplatePath: "mysql/model.grain",
			FilePath:     fmt.Sprintf("%s/internal/%s/%s", m.ProjectPath, "model", name),
			Filename:     fmt.Sprintf("%s/internal/%s/%s/%s.go", m.ProjectPath, "model", name, name),
		}
		maps["router"] = &model.CodePath{
			FS:           &stencil.MysqlTemplateFS,
			TemplatePath: "mysql/router.grain",
			FilePath:     fmt.Sprintf("%s/internal/%s/%s", m.ProjectPath, "router", name),
			Filename:     fmt.Sprintf("%s/internal/%s/%s/%s.go", m.ProjectPath, "router", name, name),
		}
		maps["handler"] = &model.CodePath{
			FS:           &stencil.MysqlTemplateFS,
			TemplatePath: "mysql/handler.grain",
			FilePath:     fmt.Sprintf("%s/internal/%s/%s", m.ProjectPath, "handler", name),
			Filename:     fmt.Sprintf("%s/internal/%s/%s/%s.go", m.ProjectPath, "handler", name, name),
		}
		maps["service"] = &model.CodePath{
			FS:           &stencil.MysqlTemplateFS,
			TemplatePath: "mysql/service.grain",
			FilePath:     fmt.Sprintf("%s/internal/%s/%s", m.ProjectPath, "service", name),
			Filename:     fmt.Sprintf("%s/internal/%s/%s/%s.go", m.ProjectPath, "service", name, name),
		}
		maps["repo"] = &model.CodePath{
			FS:           &stencil.MysqlTemplateFS,
			TemplatePath: "mysql/repo.grain",
			FilePath:     fmt.Sprintf("%s/internal/%s/%s", m.ProjectPath, "repo", name),
			Filename:     fmt.Sprintf("%s/internal/%s/%s/%s.go", m.ProjectPath, "repo", name, name),
		}

		maps["admin_handler"] = &model.CodePath{
			FS:           &stencil.MysqlAdminTemplateFS,
			TemplatePath: "mysql/admin/handler.grain",
			FilePath:     fmt.Sprintf("%s/internal/admin/handler/%s", m.ProjectPath, name),
			Filename:     fmt.Sprintf("%s/internal/admin/handler/%s/%s.go", m.ProjectPath, name, name),
		}
		maps["admin_service"] = &model.CodePath{
			FS:           &stencil.MysqlAdminTemplateFS,
			TemplatePath: "mysql/admin/service.grain",
			FilePath:     fmt.Sprintf("%s/internal/admin/service/%s", m.ProjectPath, name),
			Filename:     fmt.Sprintf("%s/internal/admin/service/%s/%s.go", m.ProjectPath, name, name),
		}
		maps["admin_repo"] = &model.CodePath{
			FS:           &stencil.MysqlAdminTemplateFS,
			TemplatePath: "mysql/admin/repo.grain",
			FilePath:     fmt.Sprintf("%s/internal/admin/repo/%s", m.ProjectPath, name),
			Filename:     fmt.Sprintf("%s/internal/admin/repo/%s/%s.go", m.ProjectPath, name, name),
		}
	} else if m.DatabaseName == "MongoDB" {
		return nil, errors.New("暂不支持MongoDB")
	} else {
		return nil, errors.New("数据库使用 MongoDB 还是 MySQL?")
	}
	return maps, nil
}

func (s *CodeAssistantService) WebGenerate(m model.Models) []string {
	mm := make(map[string]*model.CodePath)

	path := m.WebProjectPath
	if path[len(path)-1:] == "/" {
		m.WebProjectPath = path[:len(path)-1]
	}

	name := utils.ToLower(m.StructName)
	mm["web_index.vue"] = &model.CodePath{
		FS:           &stencil.WebTemplateFS,
		TemplatePath: "web/vue.grain",
		FilePath:     fmt.Sprintf("%s/src/views/business/%s", m.WebProjectPath, name),
		Filename:     fmt.Sprintf("%s/src/views/business/%s/index.vue", m.WebProjectPath, name),
	}
	mm["web_api.ts"] = &model.CodePath{
		FS:           &stencil.WebTemplateFS,
		TemplatePath: "web/api.grain",
		FilePath:     fmt.Sprintf("%s/src/api/business", m.WebProjectPath),
		Filename:     fmt.Sprintf("%s/src/api/business/%s.ts", m.WebProjectPath, name),
	}
	mm["web_zh-ch.ts"] = &model.CodePath{
		FS:           &stencil.WebTemplateFS,
		TemplatePath: "web/zh-CN.grain",
		FilePath:     fmt.Sprintf("%s/src/views/business/%s/locale/", m.WebProjectPath, name),
		Filename:     fmt.Sprintf("%s/src/views/business/%s/locale/zh-CN.ts", m.WebProjectPath, name),
	}
	mm["menu.ts"] = &model.CodePath{
		FS:           &stencil.WebTemplateFS,
		TemplatePath: "web/menu.grain",
		FilePath:     fmt.Sprintf("%s/src/router/routes/modules/business", m.WebProjectPath),
		Filename:     fmt.Sprintf("%s/src/router/routes/modules/business/%s.ts", m.WebProjectPath, name),
	}

	for _, field := range m.Fields {
		switch field.Type {
		case "int":
			field.Type = "number"
		case "int8":
			field.Type = "number"
		case "int16":
			field.Type = "number"
		case "int32":
			field.Type = "number"
		case "int64":
			field.Type = "number"
		case "uint":
			field.Type = "number"
		case "uint8":
			field.Type = "number"
		case "uint16":
			field.Type = "number"
		case "uint32":
			field.Type = "number"
		case "uint64":
			field.Type = "number"
		case "bool":
			field.Type = "boolean"
		case "time.Time":
			field.Type = "string"
		}
	}

	var errStr []string
	m.WebFieldToLower()
	for i, _ := range mm {
		//生成代码
		err := s.WebGeneratedCode(m, mm[i])
		if err != nil {
			errStr = append(errStr, err.Error()+"\n")
		}
	}
	return errStr
}

func (s *CodeAssistantService) WebGeneratedCode(m model.Models, web *model.CodePath) (err error) {
	//检查目录是否存在
	exists := utils.PathIsNotExist(web.FilePath)
	if err != nil {
		xlog.Error(err)
		return err
	}

	if exists {
		err = os.MkdirAll(web.FilePath, os.ModePerm)
		if err != nil {
			xlog.Error(err)
			return err
		}
	}

	//读取模板文件
	b, err := web.FS.ReadFile(web.TemplatePath)
	if err != nil {
		xlog.Error(err)
		return err
	}
	rt := string(b)

	//一些不知道怎么处理的 直接使用简单粗暴的方式替换处理
	rt = strings.ReplaceAll(rt, "{{ModelNameA}}", m.StructName)
	rt = strings.ReplaceAll(rt, "{{ModelNameB}}", m.Name)

	temp, err := template.New(utils.ToLower(m.Type)).Parse(rt)
	if err != nil {
		xlog.Info(err, "New模版失败", web.TemplatePath)
		return err
	}

	//检查文件是否存在
	exists = utils.FileIsNotExist(web.Filename)

	if s.conf.Gin.Model == "debug" {
		_ = os.Remove(web.Filename)
	} else {
		if !exists {
			xlog.Info("文件已存在", fmt.Sprintf("%s.go", m.Type))
			return errors.New("文件已存在: " + fmt.Sprintf("%s/%s.go", m.Type, m.Name))
		}
	}

	open, err := os.OpenFile(web.Filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		xlog.Info(err)
		return err
	}
	defer open.Close()
	//渲染模板文件
	err = temp.Execute(open, m)
	if err != nil {
		xlog.Error(err)
		return err
	}
	return nil
}

func (s *CodeAssistantService) FlutterGenerate(m model.Models) []string {
	mm := make(map[string]*model.CodePath)

	path := m.WebProjectPath
	if path[len(path)-1:] == "/" {
		m.WebProjectPath = path[:len(path)-1]
	}

	name := utils.ToLower(m.StructName)
	mm["model.dart"] = &model.CodePath{
		FS:           &stencil.FlutterTemplateFS,
		TemplatePath: "flutter/model.grain",
		FilePath:     fmt.Sprintf("%s/flutter/models/%s", m.ProjectPath, name),
		Filename:     fmt.Sprintf("%s/flutter/models/%s/%s.dart", m.ProjectPath, name, m.Name),
	}

	mm["api.dart"] = &model.CodePath{
		FS:           &stencil.FlutterTemplateFS,
		TemplatePath: "flutter/api.grain",
		FilePath:     fmt.Sprintf("%s/flutter/api/%s", m.ProjectPath, name),
		Filename:     fmt.Sprintf("%s/flutter/api/%s/%s.dart", m.ProjectPath, name, m.Name),
	}

	var errStr []string
	//m.WebFieldToLower()
	for i, _ := range mm {
		//生成代码
		err := s.FlutterGeneratedCode(m, mm[i])
		if err != nil {
			errStr = append(errStr, err.Error()+"\n")
		}
	}
	return errStr
}

func (s *CodeAssistantService) FlutterGeneratedCode(m model.Models, flutter *model.CodePath) (err error) {

	//检查目录是否存在
	exists := utils.PathIsNotExist(flutter.FilePath)
	if err != nil {
		xlog.Error(err)
		return err
	}

	if exists {
		err = os.MkdirAll(flutter.FilePath, os.ModePerm)
		if err != nil {
			xlog.Error(err)
			return err
		}
	}

	//读取模板文件
	b, err := flutter.FS.ReadFile(flutter.TemplatePath)
	if err != nil {
		xlog.Error(err)
		return err
	}
	rt := string(b)

	temp, err := template.New(utils.ToLower(m.Type)).Parse(rt)
	if err != nil {
		xlog.Info(err, "New模版失败", flutter.TemplatePath)
		return err
	}

	//检查文件是否存在
	exists = utils.FileIsNotExist(flutter.Filename)

	if s.conf.Gin.Model == "debug" {
		_ = os.Remove(flutter.Filename)
	} else {
		if !exists {
			return errors.New("文件已存在: " + fmt.Sprintf("%s.dart", m.Name))
		}
	}

	open, err := os.OpenFile(flutter.Filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		xlog.Info(err)
		return err
	}
	defer open.Close()

	for _, field := range m.Fields {
		if field.Type == "" {
			field.Type = "String?"
		} else {
			switch field.Type {
			case "string":
				field.Type = "String?"
			case "int", "int64", "float32", "float64":
				field.Type = "num?"
			case "bool":
				field.Type = "bool?"
			default:
				field.Type += "?"
			}
		}

	}

	//渲染模板文件
	err = temp.Execute(open, m)
	if err != nil {
		xlog.Error(err)
		return err
	}
	return nil
}

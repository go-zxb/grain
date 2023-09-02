package model

import (
	"embed"
	utils "github.com/go-grain/go-utils"
)

type ViewCode struct {
	Model   string `json:"model"`
	Router  string `json:"router"`
	Handle  string `json:"handle"`
	Service string `json:"service"`
	Repo    string `json:"repo"`
	Vue     string `json:"vue"`
	APi     string `json:"api"`
	ZhCN    string `json:"zhCn"`
}

type Project struct {
	Model
	ProjectName    string    `json:"projectName" gorm:"comment:项目名称"`
	ProjectPath    string    `json:"projectPath"  gorm:"comment:项目路径"`
	WebProjectPath string    `json:"webProjectPath" gorm:"comment:Web项目路径"`
	Description    string    `json:"description" gorm:"description:描述"`
	Models         []*Models `json:"model" gorm:"-"`
	IsInit         string    `json:"is_init" gorm:"default:no;comment:是否已初始化"`
}

func (Project) TableName() string {
	return "project"
}

type CreateProject struct {
	ProjectName    string `json:"projectName" gorm:"comment:项目名称"`
	ProjectPath    string `json:"projectPath"  gorm:"comment:项目路径"`
	WebProjectPath string `json:"webProjectPath" gorm:"comment:Web项目路径"`
	Description    string `json:"description" gorm:"description:描述"`
}

type ProjectReq struct {
	PageReq
}

type UpdateProject struct {
	UpdateType string   `json:"updateType"`
	Project    *Project `json:"project"`
	Models     *Models  `json:"model"`
	Fields     *Fields  `json:"fields"`
}

type Models struct {
	Model
	ParentId          uint      `json:"parentId" gorm:"comment:父ID"`
	ProjectName       string    `json:"projectName" gorm:"-"`
	ProjectPath       string    `json:"projectPath"  gorm:"-"`
	WebProjectPath    string    `json:"webProjectPath" gorm:"-"`
	StructName        string    `json:"structName"  gorm:"comment:结构体名称"`
	ToLowerStructName string    `json:"toLowerStructName"  gorm:"-"`
	Description       string    `json:"description" gorm:"description:描述"`
	QueryTime         string    `json:"queryTime"  gorm:"comment:时间范围查询"`
	IsInit            string    `json:"isInit" gorm:"default:no;comment:是否已初始化"`
	DatabaseName      string    `json:"databaseName"  gorm:"comment:用什么数据库 MySQL MongoDB?"`
	Name              string    `json:"name" gorm:"-"`
	Type              string    `gorm:"-"`
	FirstLetter       string    `gorm:"-"`
	IsQueryCriteria   bool      `gorm:"-"`
	Fields            []*Fields `json:"fields"  gorm:"-"`
}

func (Models) TableName() string {
	return "models"
}

type CreateModels struct {
	StructName   string `json:"structName"  gorm:"comment:结构体名称"`
	DatabaseName string `json:"databaseName"  gorm:"comment:用什么数据库 MySQL MongoDB?"`
	Description  string `json:"description" gorm:"description:描述"`
	QueryTime    string `json:"queryTime"  gorm:"comment:时间范围查询"`
}

type Fields struct {
	Model
	ParentId      uint   `json:"parentId"  gorm:"comment:父ID"`
	Name          string `json:"name" gorm:"comment:字段名"`
	NameLower     string `json:"name_lower" gorm:"-"`
	Type          string `json:"type" gorm:"comment:字段类型"`
	JsonTag       string `json:"jsonTag" gorm:"comment:Json标签tag"`
	Description   string `json:"description" gorm:"description:描述"`
	QueryCriteria string `json:"queryCriteria"  gorm:"comment:查询条件"`
	MysqlField    string `json:"mysqlField" gorm:"column:sql_field;comment:Mysql字段名"`
	RequiredValue string `json:"requiredValue" gorm:"comment:校验规则"`
	Required      string `json:"required" gorm:"comment:是否必传参数"`
}

func (Fields) TableName() string {
	return "fields"
}

type CodePath struct {
	TemplatePath string
	FilePath     string
	Filename     string
	Type         string
	FS           *embed.FS
}

type WebStruct struct {
	FileName string
	FilePath string
	Suffix   string
	Template string
	Type     int
}

func (m *Models) GoFieldTo() {
	//处理字段
	m.Name = utils.ToLower(m.StructName)
	m.ToLowerStructName = utils.ToLower(m.StructName)
	m.StructName = utils.ToTitle(m.StructName)
	for _, f := range m.Fields {

		if f.Type == "" {
			f.Type = "string"
		}

		if f.JsonTag == "" {
			f.JsonTag = f.Name
		}

		f.Type = utils.ToLower(f.Type)
		f.Name = utils.ToTitle(f.Name)
		f.NameLower = utils.ToLower(f.Name)
		f.JsonTag = utils.ToLower(f.JsonTag)

	}
}

func (m *Models) WebFieldToLower() {
	for _, f := range m.Fields {

		if f.QueryCriteria != "" {
			m.IsQueryCriteria = true
		}

		if f.Type == "" {
			f.Type = "string"
		}

		f.NameLower = utils.ToLower(f.Name)
		f.Name = utils.ToLower(f.Name)
		f.Type = utils.ToLower(f.Type)

	}
}

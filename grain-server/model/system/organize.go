package model

type Organize struct {
	Model
	ParentId uint   `form:"parentId" json:"parentId" gorm:"comment:父ID"`
	Name     string `form:"name" json:"name" binding:"required" gorm:"comment:组织或部门名称"`
	Leader   string `form:"leader" json:"leader" gorm:"comment:部门领导"`
	OeType   int    `form:"oeType" json:"oeType" gorm:"comment:1是组织,2是部门,3是员工"`
	Position string `form:"position" json:"position"  gorm:"comment:职位名称"`
}

func (Organize) TableName() string {
	return "organize"
}

type CreateOrganize struct {
	ParentId uint   `form:"parentId" json:"parentId" gorm:"comment:父ID"`
	Name     string `form:"name" json:"name" binding:"required" gorm:"comment:组织或部门名称"`
	OeType   int    `json:"oeType"`
}

type UpdateOrganize struct {
	ParentId uint   `form:"parentId" json:"parentId"`
	Name     string `form:"name" json:"name"`
	OeType   int    `json:"oeType"`
}

type OrganizeQuery struct {
	PageReq
	ID    uint   `form:"id" json:"id"`
	QType string `form:"qType" json:"qType"`
}

type OrganizeRes struct {
	Model
	ParentId uint   `json:"parentId" gorm:"comment:父ID"`
	Name     string `json:"name" gorm:"comment:组织或部门名称"`
	OeType   int    `json:"oeType"`
}

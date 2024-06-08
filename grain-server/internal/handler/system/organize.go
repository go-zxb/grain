package handler

import (
	"github.com/gin-gonic/gin"
	service "github.com/go-grain/grain/internal/service/system"
	model "github.com/go-grain/grain/model/system"
	"github.com/go-grain/grain/pkg/response"
	"github.com/go-grain/grain/utils/const"
	"strconv"
)

type OrganizeHandle struct {
	res response.Response
	sv  *service.OrganizeService
}

func NewOrganizeHandle(sv *service.OrganizeService) *OrganizeHandle {
	return &OrganizeHandle{
		sv: sv,
	}
}

// CreateOrganize
// @Security ApiKeyAuth
// @Summary 创建组织管理
// @Description 创建组织管理
// @Tags 组织管理
// @Accept json
// @Produce json
// @Param data body  model.CreateOrganize true "组织管理信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Router /organize [post]
func (r *OrganizeHandle) CreateOrganize(ctx *gin.Context) {
	reply := r.res.New()
	organize := model.Organize{}
	err := ctx.ShouldBindJSON(&organize)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateOrganize(&organize, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("创建组织管理成功").Success(ctx)
}

// UpdateOrganize
// @Security ApiKeyAuth
// @Summary 更新组织管理
// @Description 更新组织管理信息
// @Tags 组织管理
// @Accept json
// @Produce json
// @Param data body  model.UpdateOrganize true "更新组织管理信息"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /organize [put]
func (r *OrganizeHandle) UpdateOrganize(ctx *gin.Context) {
	reply := r.res.New()
	organize := model.Organize{}
	err := ctx.ShouldBindJSON(&organize)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage("解析组织管理参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateOrganize(&organize, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("更新组织管理成功").Success(ctx)
}

// GetOrganizeById
// @Security ApiKeyAuth
// @Summary 根据组织管理ID获取信息
// @Description 根据组织管理ID获取信息
// @Tags 组织管理
// @Accept json
// @Produce json
// @Param id query  int true "组织管理ID "
// @Success 200 {object} model.Organize "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /organize [get]
func (r *OrganizeHandle) GetOrganizeById(ctx *gin.Context) {
	reply := r.res.New()
	organizeId, _ := strconv.Atoi(ctx.Query("id"))
	organizeInfo, err := r.sv.GetOrganizeById(uint(organizeId), ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").WithData(organizeInfo).Success(ctx)
}

// GetOrganizeList
// @Security ApiKeyAuth
// @Summary 获取组织管理分页数据
// @Description 获取组织管理分页数据
// @Tags 组织管理
// @Accept json
// @Produce json
// @Param data body model.OrganizeQuery true "分页列表请求参数"
// @Success 200 {object} model.Organize "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /organize/list [get]
func (r *OrganizeHandle) GetOrganizeList(ctx *gin.Context) {
	reply := r.res.New()
	query := model.OrganizeQuery{}
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetOrganizeList(&query, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").WithData(list).
		WithTotal(query.Total).
		WithPage(query.Page).
		WithPageSize(query.PageSize).
		Success(ctx)
}

func (r *OrganizeHandle) GetOrganizeListGroup(ctx *gin.Context) {
	reply := r.res.New()
	query := model.OrganizeQuery{}
	err := ctx.ShouldBindQuery(&query)
	list, err := r.sv.GetOrganizeListGroup(&query, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("成功").WithData(list).Success(ctx)
}

// DeleteOrganizeById
// @Security ApiKeyAuth
// @Summary 删除组织管理
// @Description 根据组织管理ID删除组织管理
// @Tags 组织管理
// @Accept json
// @Produce json
// @Param id query  int true "根据组织管理ID删除组织管理 "
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /organize [delete]
func (r *OrganizeHandle) DeleteOrganizeById(ctx *gin.Context) {
	reply := r.res.New()
	organizeId, _ := strconv.Atoi(ctx.Query("id"))
	err := r.sv.DeleteOrganizeById(uint(organizeId), ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage("删除组织管理失败").Fail(ctx)
		return
	}
	reply.WithMessage("删除组织管理成功").Success(ctx)
}

// DeleteOrganizeByIds
// @Security ApiKeyAuth
// @Summary 删除组织管理
// @Description "根据组织管理ID批量删除组织管理"
// @Tags 组织管理
// @Accept json
// @Produce json
// @Param data body []uint true "根据组织管理ID批量删除组织管理"
// @Success 200 {object} model.ErrorRes "成功"
// @Failure 400 {object} model.ErrorRes "格式错误"
// @Failure 401 {object} model.ErrorRes "未经授权"
// @Failure 404 {object} model.ErrorRes "资源不存在"
// @Router /organize/organizeByIds [delete]
func (r *OrganizeHandle) DeleteOrganizeByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		OrganizeIds []uint `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.DeleteOrganizeByIds(api.OrganizeIds, ctx)
	if err != nil {
		reply.WithCode(consts.ReqFail).WithMessage("删除组织管理失败").Fail(ctx)
		return
	}
	reply.WithMessage("删除组织管理成功").Success(ctx)
}

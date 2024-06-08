package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/config"
	"github.com/go-grain/grain/log"
	model "github.com/go-grain/grain/model/system"
	redisx "github.com/go-grain/grain/pkg/redis"
	"github.com/go-pay/gopay/pkg/xlog"
)

type IOrganizeRepo interface {
	CreateOrganize(organize *model.Organize) error
	GetOrganizeById(id uint) (u *model.Organize, err error)
	GetOrganizeList(req *model.OrganizeQuery) ([]*model.Organize, error)
	//GetOrganizeListGroup(req *model.OrganizeQuery) ([]*model.Organize, error)
	UpdateOrganize(organize *model.Organize) error
	DeleteOrganizeById(organizeId uint) error
	DeleteOrganizeByIds(organizeIds []uint) error
}

type OrganizeService struct {
	repo IOrganizeRepo
	rdb  redisx.IRedis
	conf *config.Config
	log  *log.Helper
}

func NewOrganizeService(repo IOrganizeRepo, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *OrganizeService {
	return &OrganizeService{
		repo: repo,
		rdb:  rdb,
		conf: conf,
		log:  log.NewHelper(logger),
	}
}

func (s *OrganizeService) CreateOrganize(organize *model.Organize, ctx *gin.Context) error {
	if err := s.repo.CreateOrganize(organize); err != nil {
		s.log.Errorw("errMsg", "创建项目", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "创建组织管理")
	return nil
}

func (s *OrganizeService) GetOrganizeById(organizeId uint, ctx *gin.Context) (*model.Organize, error) {
	return s.repo.GetOrganizeById(organizeId)
}

func (s *OrganizeService) GetOrganizeList(req *model.OrganizeQuery, ctx *gin.Context) ([]*model.Organize, error) {
	list, err := s.repo.GetOrganizeList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	return list, err
}

type dt struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	Item []dt   `json:"item"`
}

type oe struct {
	OeType   int    `json:"oeType"`
	Name     string `json:"name"`
	ID       uint   `json:"id"`
	Children []dt   `json:"children"`
}

func (s *OrganizeService) GetOrganizeListGroup(req *model.OrganizeQuery, ctx *gin.Context) (interface{}, error) {
	list, err := s.repo.GetOrganizeList(req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("暂无更多数据")
	}
	var o []*oe

	//遍历组织
	xlog.Info(len(list))
	for _, organize := range list {
		//查询xx组织下的部门
		req.QType = "2"
		req.ID = organize.ID
		department, err := s.repo.GetOrganizeList(req)
		if err != nil {
			continue
		}

		//xx组织
		oen := &oe{
			OeType: 1,
			ID:     organize.ID,
			Name:   organize.Name,
		}

		//创建一个部门
		var dtn []dt
		//遍历部门数据 获取部门下面的职位
		for _, m := range department {

			d := dt{
				Name: m.Name,
				ID:   m.ID,
			}

			//获取部门下面的职位
			req.QType = "3"
			req.ID = m.ID
			positions, err := s.repo.GetOrganizeList(req)
			if err != nil {
				continue
			}

			// 把职位数据遍历添加到dtn部门下面
			for _, ps := range positions {
				d.Item = append(d.Item, dt{
					Name: ps.Name,
					ID:   ps.ID,
				})
			}
			dtn = append(dtn, d)
		}

		oen.Children = dtn
		o = append(o, oen)
	}

	return o, err
}

func (s *OrganizeService) UpdateOrganize(organize *model.Organize, ctx *gin.Context) error {
	if err := s.repo.UpdateOrganize(organize); err != nil {
		s.log.Errorw("errMsg", "更新组织管理", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "更新组织管理")
	return nil
}

func (s *OrganizeService) DeleteOrganizeById(organizeId uint, ctx *gin.Context) error {
	if err := s.repo.DeleteOrganizeById(organizeId); err != nil {
		s.log.Errorw("errMsg", "删除组织管理", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "删除组织管理")
	return nil
}

func (s *OrganizeService) DeleteOrganizeByIds(organizeIds []uint, ctx *gin.Context) error {
	if err := s.repo.DeleteOrganizeByIds(organizeIds); err != nil {
		s.log.Errorw("errMsg", "批量删除组织管理", "err", err.Error())
		return err
	}
	s.log.Infow("errMsg", "批量删除组织管理")
	return nil
}

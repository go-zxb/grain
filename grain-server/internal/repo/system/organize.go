package repo

import (
	"fmt"
	utils "github.com/go-grain/go-utils"
	"github.com/go-grain/go-utils/redis"
	"github.com/go-grain/grain/internal/repo/data"
	"github.com/go-grain/grain/internal/repo/system/query"
	service "github.com/go-grain/grain/internal/service/system"
	model "github.com/go-grain/grain/model/system"
	"gorm.io/gorm"
)

type OrganizeRepo struct {
	db    *data.DB
	rdb   redis.IRedis
	query *query.Query
}

func NewOrganizeRepo(db *gorm.DB, rdb redis.IRedis) service.IOrganizeRepo {
	//为了偷懒自动生成代码后直接在这里AutoMigrate,你可以放到data里面去统一管理
	_ = db.AutoMigrate(model.Organize{})
	//SetDefault 你也可以放到core>start文件里面去统一初始化
	query.SetDefault(db)
	return &OrganizeRepo{
		db:    &data.DB{DB: db},
		rdb:   rdb,
		query: query.Q,
	}
}

func (r *OrganizeRepo) CreateOrganize(organize *model.Organize) error {
	return r.query.Organize.Create(organize)
}

func (r *OrganizeRepo) UpdateOrganize(organize *model.Organize) error {
	if _, err := r.query.Organize.Where(r.query.Organize.ID.Eq(organize.ID)).Updates(organize); err != nil {
		return err
	}
	Neworganize, err := r.query.Organize.Where(r.query.Organize.ID.Eq(organize.ID)).First()
	if err != nil {
		return err
	}
	_ = r.rdb.SetObject(fmt.Sprintf("%s:%d", organize.TableName(), organize.ID), Neworganize, 180)
	return nil
}

func (r *OrganizeRepo) GetOrganizeById(id uint) (organize *model.Organize, err error) {
	err = r.rdb.GetObject("", organize)
	if err != nil || organize.ID == 0 {
		organize, err = r.query.Organize.Where(r.query.Organize.ID.Eq(id)).First()
		if err != nil {
			return nil, err
		}
		r.rdb.SetObject(utils.ToString(id), organize, 180)
	}
	return
}

func (r *OrganizeRepo) GetOrganizeList(req *model.OrganizeQuery) (list []*model.Organize, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 || req.PageSize >= 100 {
		req.PageSize = 20
	}

	o := r.query.Organize
	q := o.Where()

	if req.QType == "1" {
		q = q.Where(o.OeType.Eq(1))
	}

	if req.QType == "2" {
		q = q.Where(o.OeType.Eq(2))
	}

	if req.QType == "3" {
		q = q.Where(o.OeType.Eq(3))
	}

	if req.ID != 0 {
		q = q.Where(o.ParentId.Eq(req.ID))
	}

	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	req.Total = count
	q = q.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize)
	list, err = q.Order(o.CreatedAt).Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *OrganizeRepo) DeleteOrganizeById(id uint) error {
	if _, err := r.query.Organize.Where(r.query.Organize.ID.Eq(id)).Delete(); err != nil {
		return err
	}
	return nil
}

func (r *OrganizeRepo) DeleteOrganizeByIds(ids []uint) error {
	if _, err := r.query.Organize.Where(r.query.Organize.ID.In(ids...)).Delete(); err != nil {
		return err
	}
	return nil
}

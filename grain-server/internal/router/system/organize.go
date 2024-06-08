package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/grain/config"
	handler "github.com/go-grain/grain/internal/handler/system"
	repo "github.com/go-grain/grain/internal/repo/system"
	service "github.com/go-grain/grain/internal/service/system"
	"github.com/go-grain/grain/log"
	"github.com/go-grain/grain/middleware"
	redisx "github.com/go-grain/grain/pkg/redis"
	"gorm.io/gorm"
)

type OrganizeRouter struct {
	public  gin.IRoutes
	private gin.IRoutes
	api     *handler.OrganizeHandle
}

func NewOrganizeRouter(routerGroup *gin.RouterGroup, db *gorm.DB, rdb redisx.IRedis, conf *config.Config, logger log.Logger, enforcer *casbin.CachedEnforcer) *OrganizeRouter {
	data := repo.NewOrganizeRepo(db, rdb)
	sv := service.NewOrganizeService(data, rdb, conf, logger)
	return &OrganizeRouter{
		public:  routerGroup.Group("organize"),
		api:     handler.NewOrganizeHandle(sv),
		private: routerGroup.Group("organize").Use(middleware.JwtAuth(rdb), middleware.Casbin(enforcer)),
	}
}

func (r *OrganizeRouter) InitRouters() {
	r.private.PUT("", r.api.UpdateOrganize)
	r.private.POST("", r.api.CreateOrganize)
	r.private.GET("list", r.api.GetOrganizeList)
	r.private.GET("listGroup", r.api.GetOrganizeListGroup)
	r.private.DELETE("organizeById", r.api.DeleteOrganizeById)
	r.private.DELETE("organizeByIds", r.api.DeleteOrganizeByIds)
}

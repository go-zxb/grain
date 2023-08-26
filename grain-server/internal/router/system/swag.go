package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouterSwag(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

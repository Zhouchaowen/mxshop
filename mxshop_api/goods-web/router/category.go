package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods-web/api/category"
	"mxshop_api/goods-web/api/goods"
	"mxshop_api/goods-web/middlewares"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categories")
	{
		CategoryRouter.GET("", category.List)
		CategoryRouter.DELETE("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),category.Delete)
		CategoryRouter.GET("/:id",category.Detail)
		CategoryRouter.PUT("",middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.New)
		CategoryRouter.PATCH("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.Update)
	}
}

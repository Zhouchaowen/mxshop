package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods-web/api/goods"
	"mxshop_api/user-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("",middlewares.JWTAuth(),middlewares.IsAdminAuth(), goods.New) // 管理员
		GoodsRouter.GET("/:id",goods.Detail)
		GoodsRouter.DELETE("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.Delete)
		GoodsRouter.GET("/:id/stocks",goods.Stocks)
		GoodsRouter.PUT("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.Update)
		GoodsRouter.PATCH("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.UpdateStatus)
	}
}

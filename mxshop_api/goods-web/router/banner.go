package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_api/goods-web/api/banners"
	"mxshop_api/goods-web/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("banner")
	{
		CategoryRouter.GET("", banners.List)
		CategoryRouter.DELETE("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),banners.Delete)
		CategoryRouter.PUT("",middlewares.JWTAuth(),middlewares.IsAdminAuth(),banners.New)
		CategoryRouter.PATCH("/:id",middlewares.JWTAuth(),middlewares.IsAdminAuth(),banners.Update)
	}
}


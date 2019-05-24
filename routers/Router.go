package routers

import (
	"github.com/gin-gonic/gin"
	"shorturl/controllers"
)

func Route(Router *gin.Engine) {

	Router.GET("/:code", controllers.Path)

	api := Router.Group("v1")
	{
		api.POST("/create", controllers.Create)
		api.POST("/query", controllers.Query)
	}
}

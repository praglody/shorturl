package routers

import (
	"github.com/gin-gonic/gin"
	. "shorturl/controllers"
	"shorturl/middlewares"
)

func Route(Router *gin.Engine) {

	Router.GET("/:code", Index.Path)

	api := Router.Group("v1").Use(middlewares.Auth())
	{
		api.POST("/create", Index.Create)
		api.POST("/multicreate", Index.MultiCreate)
		api.POST("/query", Index.Query)
	}
}

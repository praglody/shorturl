package routers

import (
	"github.com/gin-gonic/gin"
	. "shorturl/controllers"
	"shorturl/middlewares"
)

func Route(Router *gin.Engine) {
	Router.Use(middlewares.Request()).
		StaticFile("/", "./public").
		GET("/:code", Short.Path)

	api := Router.Group("api/v1").Use(middlewares.Request())
	{
		api.POST("/create", Short.Create)
		api.POST("/multicreate", Short.MultiCreate)
		api.POST("/query", Short.Query)
	}
}

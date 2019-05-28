package routers

import (
	"github.com/gin-gonic/gin"
	. "shorturl/controllers"
)

func Route(Router *gin.Engine) {

	Router.GET("/:code", Index.Path)

	api := Router.Group("v1")
	{
		api.POST("/create", Index.Create)
		api.POST("/multicreate", Index.MultiCreate)
		api.POST("/query", Index.Query)
	}
}

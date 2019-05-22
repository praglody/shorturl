package app

import (
	"github.com/gin-gonic/gin"
	"github.com/orca-zhang/lrucache"
	"shorturl/lib/dao"
)

var lc *lrucache.LRUCache

func init() {
	lc = lrucache.New(1000)
}

func Run() {
	r := gin.Default()
	r.GET("/:path", func(context *gin.Context) {
		path := context.Param("path")
		if len(path) != 6 {
			context.AbortWithStatus(204)
			return
		}

		url, ok := lc.Get(path)
		if ok {
			context.Header("Location", url.(string))
			context.AbortWithStatus(302)
			return
		}

		shorturl, err := dao.GetRow(path)
		if err != nil {
			context.AbortWithStatus(204)
			return
		}

		lc.Put(path, shorturl.Url)
		context.Header("Location", shorturl.Url)
		context.AbortWithStatus(302)
	})

	r.Run("0.0.0.0:8080")
}

package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl/commons"
)

type Controller struct{}

func (*Controller) success(c *gin.Context, data map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": commons.Success,
		"msg":  "ok",
		"data": data,
	})
}

func (*Controller) failed(c *gin.Context, code int, msg string) {
	logs.Error("failed, code: %d, msg: %s", code, msg)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": "",
	})
}

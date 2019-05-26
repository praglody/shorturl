package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl/models"
)

type Controller struct{}

func (*Controller) success(c *gin.Context, data map[string]interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": models.Success,
		"msg":  "ok",
		"data": data,
	})
}

func (*Controller) failed(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": "",
	})
	c.Abort()
}

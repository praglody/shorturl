package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl/models"
	"shorturl/services"
	"strings"
)

var Index *IndexController

type IndexController struct {
	Controller
}

func init() {
	Index = &IndexController{}
}

func (i *IndexController) Create(c *gin.Context) {
	url := c.PostForm("url")
	logs.Info("incoming create url request, url: " + url)
	if url == "" {
		logs.Error("url is error, url: " + url)
		i.failed(c, models.ParamsError, "参数错误")
		c.Abort()
		return
	}
	if !strings.HasPrefix(url, "http") {
		logs.Error("url is error, url: " + url)
		i.failed(c, models.ParamsError, "请输入合法的url，以http开头！")
		c.Abort()
		return
	}

	code, err := services.UrlService{}.GenCode(url)
	if err != nil {
		logs.Error("gen code failed, error: " + err.Error())
		i.failed(c, models.Failed, "请求出错")
		c.Abort()
		return
	} else {
		logs.Info("get code: " + code + " for url: " + url)
		i.success(c, gin.H{
			"code": models.Conf.AppUrl + code,
		})
		return
	}
}

func (i *IndexController) Query(c *gin.Context) {
	code := c.PostForm("code")
	logs.Info("incoming query, code: " + code)
	if len(code) < 3 || len(code) > 6 {
		i.failed(c, models.ParamsError, "参数错误")
	}
	url, err := services.UrlService{}.RecCode(code)
	if err != nil {
		i.failed(c, models.NotFound, err.Error())
		c.Abort()
		return
	} else {
		i.success(c, gin.H{
			"url": url,
		})
		return
	}
}

func (i *IndexController) Path(c *gin.Context) {
	code := c.Param("code")
	logs.Info("incoming query, code: " + code)
	if len(code) < 3 || len(code) > 6 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	url, err := services.UrlService{}.RecCode(code)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	logs.Info("querying...get url: " + url + " for code: " + code)
	c.Header("Location", url)
	c.AbortWithStatus(302)
	return
}

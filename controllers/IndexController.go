package controllers

import (
	"github.com/asaskevich/govalidator"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"shorturl/models"
	"shorturl/services"
	"strings"
)

var Index = &IndexController{}

type IndexController struct {
	Controller
}

func (i *IndexController) Create(c *gin.Context) {
	lUrl := c.PostForm("url")
	logs.Info("incoming create url request, url: " + lUrl)
	if lUrl == "" {
		logs.Error("url is empty, url: " + lUrl)
		i.failed(c, models.ParamsError, "参数错误")
		return
	}

	if ok := govalidator.IsURL(lUrl); !ok {
		logs.Error("url is invalid, url: " + lUrl)
		i.failed(c, models.ParamsError, "请输入合法的url，以http开头")
		return
	}

	code, err := services.UrlService{}.GenCode(lUrl)
	if err != nil {
		logs.Error("gen code failed, error: " + err.Error())
		i.failed(c, models.Failed, "请求出错")
		return
	} else {
		logs.Info("[create]: " + lUrl + " => " + code)
		i.success(c, gin.H{
			"code": models.Conf.AppUrl + code,
		})
		return
	}
}

func (i *IndexController) Query(c *gin.Context) {
	sUrl := c.PostForm("sUrl")

	parse, err := url.Parse(sUrl)
	if err != nil {
		i.failed(c, models.ParamsError, err.Error())
		return
	}
	code := strings.Trim(parse.Path, "/")
	if len(code) < 3 || len(code) > 6 {
		i.failed(c, models.ParamsError, "参数错误")
		return
	}
	lUrl, err := services.UrlService{}.RecCode(code)
	if err != nil {
		i.failed(c, models.NotFound, err.Error())
		return
	} else {
		i.success(c, gin.H{
			"url": lUrl,
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
	lUrl, err := services.UrlService{}.RecCode(code)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	logs.Info("[query]: " + code + " => " + lUrl)
	c.Header("Location", lUrl)
	c.AbortWithStatus(302)
	return
}

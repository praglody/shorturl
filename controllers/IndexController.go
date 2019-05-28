package controllers

import (
	"encoding/json"
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

type CreateRequest struct {
	Url []string
}

type result struct {
	url  string
	code string
}

//单个生成短网址
func (i *IndexController) Create(c *gin.Context) {
	lUrl := c.PostForm("url")
	logs.Info("incoming create url request, url: " + lUrl)
	if lUrl == "" {
		logs.Info("url is empty, url: " + lUrl)
		i.failed(c, models.ParamsError, "参数错误")
		return
	}

	if ok := govalidator.IsURL(lUrl); !ok {
		logs.Info("url is invalid, url: " + lUrl)
		i.failed(c, models.ParamsError, "无效的url")
		return
	}

	shortUrl, err := services.UrlService{}.GenShortUrl(lUrl)
	if err != nil {
		logs.Error("gen shortUrl failed, error: " + err.Error())
		i.failed(c, models.Failed, "请求出错")
		return
	} else {
		logs.Info("[create]: " + lUrl + " => " + shortUrl)
		i.success(c, gin.H{
			"url": shortUrl,
		})
		return
	}
}

//批量生成短网址
func (i *IndexController) MultiCreate(c *gin.Context) {
	var urls CreateRequest
	err := c.ShouldBindJSON(&urls)
	if err != nil {
		i.failed(c, models.ParamsError, "参数错误")
		return
	}
	if len(urls.Url) == 0 {
		i.failed(c, models.ParamsError, "url不能为空")
		return
	}
	if len(urls.Url) > 50 {
		i.failed(c, models.ParamsError, "最多可同时生成50个")
		return
	}

	str, _ := json.Marshal(urls.Url)
	logs.Info("incoming multicreate url request, url: " + string(str))

	var cCode = make(chan result)
	for _, v := range urls.Url {
		go func(lUrl string) {
			if ok := govalidator.IsURL(lUrl); !ok {
				logs.Info("url is invalid, url: " + lUrl)
				cCode <- result{lUrl, "url is not valid"}
				return
			}
			shortUrl, err := services.UrlService{}.GenShortUrl(lUrl)
			if err != nil {
				logs.Error("gen shortUrl failed, error: " + err.Error())
				cCode <- result{lUrl, err.Error()}
			} else {
				logs.Info("[create]: " + lUrl + " => " + shortUrl)
				cCode <- result{lUrl, shortUrl}
			}
		}(v)
	}

	var results = make(map[string]interface{})
	for {
		res := <-cCode
		results[res.url] = res.code
		if len(results) == len(urls.Url) {
			close(cCode)
			i.success(c, gin.H{"urls": results})
			return
		}
	}
}

func (i *IndexController) Query(c *gin.Context) {
	sUrl := c.PostForm("url")

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
	lUrl, err := services.UrlService{}.RestoreUrl(code)
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
	lUrl, err := services.UrlService{}.RestoreUrl(code)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	logs.Info("[query]: " + code + " => " + lUrl)
	c.Header("Location", lUrl)
	c.AbortWithStatus(302)
	return
}

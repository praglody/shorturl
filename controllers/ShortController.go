package controllers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"shorturl/app"
	"shorturl/commons"
	"strings"
)

var Index = &ShortController{
	urlCode: app.NewService(),
}

type ShortController struct {
	Controller
	urlCode app.Service
}

type CreateRequest struct {
	Urls []string
}

type result struct {
	url  string
	code string
}

//单个生成短网址
func (s *ShortController) Create(c *gin.Context) {
	userId := c.GetInt("userId")
	lUrl := c.PostForm("url")
	logs.Info("incoming create url request, url: " + lUrl)
	if lUrl == "" {
		logs.Info("url is empty, url: " + lUrl)
		s.failed(c, commons.ParamsError, "参数错误")
		return
	}

	if ok := govalidator.IsURL(lUrl); !ok {
		logs.Info("url is invalid, url: " + lUrl)
		s.failed(c, commons.ParamsError, "无效的url")
		return
	}
	shortUrl, err := s.urlCode.GenShortUrl(lUrl, userId)
	if err != nil {
		logs.Error("gen shortUrl failed, error: " + err.Error())
		s.failed(c, commons.Failed, "请求出错")
		return
	} else {
		logs.Info("[create]: " + lUrl + " => " + shortUrl)
		s.success(c, gin.H{
			"url": shortUrl,
		})
		return
	}
}

//批量生成短网址
func (s *ShortController) MultiCreate(c *gin.Context) {
	var request CreateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		s.failed(c, commons.ParamsError, "参数错误")
		return
	}
	if len(request.Urls) == 0 {
		s.failed(c, commons.ParamsError, "url不能为空")
		return
	}
	if len(request.Urls) > 50 {
		s.failed(c, commons.ParamsError, "最多可同时生成50个")
		return
	}

	str, _ := json.Marshal(request.Urls)
	logs.Info("incoming multicreate url request, url: " + string(str))
	userId := c.GetInt("userId")
	resultCode := make(chan result)
	for _, v := range request.Urls {
		go func(lUrl string) {
			if ok := govalidator.IsURL(lUrl); !ok {
				logs.Info("url is invalid, url: " + lUrl)
				resultCode <- result{lUrl, "url is not valid"}
				return
			}
			shortUrl, err := s.urlCode.GenShortUrl(lUrl, userId)
			if err != nil {
				logs.Error("gen shortUrl failed, error: " + err.Error())
				resultCode <- result{lUrl, err.Error()}
			} else {
				logs.Info("[create]: " + lUrl + " => " + shortUrl)
				resultCode <- result{lUrl, shortUrl}
			}
		}(v)
	}

	var results = make(map[string]interface{})
	var count = 1
	for {
		res := <-resultCode
		results[res.url] = res.code
		if count == len(request.Urls) {
			close(resultCode)
			s.success(c, gin.H{"urls": results})
			return
		}
		count++
	}
}

func (s *ShortController) Query(c *gin.Context) {
	sUrl := c.PostForm("url")
	parse, err := url.Parse(sUrl)
	if err != nil {
		s.failed(c, commons.ParamsError, err.Error())
		return
	}
	code := strings.Trim(parse.Path, "/")
	lUrl, err := s.urlCode.RestoreUrl(code)
	if err != nil {
		s.failed(c, commons.NotFound, err.Error())
		return
	} else {
		s.success(c, gin.H{
			"url": lUrl,
		})
		return
	}
}

func (s *ShortController) Path(c *gin.Context) {
	code := c.Param("code")
	logs.Info("incoming query, code: " + code)
	lUrl, err := s.urlCode.RestoreUrl(code)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	logs.Info("[query]: " + code + " => " + lUrl)
	if !strings.HasPrefix(lUrl, "http") {
		lUrl = "http://" + lUrl
	}
	c.Header("Location", lUrl)
	c.AbortWithStatus(302)
	return
}

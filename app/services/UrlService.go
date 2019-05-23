package services

import (
	"errors"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"shorturl/models"
	"time"
)

type UrlService struct{}

var bm cache.Cache

func init() {
	bm, _ = cache.NewCache("memory", `{"interval":60}`)
}

func (UrlService) GenCode(url string) (code string, err error) {

	var urlCode models.UrlCode

	existed := urlCode.GetCode(url)

	if existed.Code != "" {
		return existed.Code, nil
	}

	code = getValidCode()
	if code == "" {
		return code, errors.New("gen code failed")
	}

	logs.Info("add new short url, code: " + code)
	err = urlCode.AddUrl(url, code)
	if err != nil {
		return "", err
	}
	_ = bm.Put(code, url, time.Hour)

	return code, nil
}

func (UrlService) RecCode(code string) (string, error) {
	//cache
	url := cache.GetString(bm.Get(code))
	if url == "" {
		var urlCode models.UrlCode
		result := urlCode.GetUrl(code)
		if result.Url == "" {
			return "", errors.New("code not existed")
		}
		url = result.Url
		_ = bm.Put(code, url, time.Hour)
	}

	return url, nil
}

func getValidCode() string {
	var code string
	var urlCode models.UrlCode

	code = genRandomCode()

	for i := 0; i < 5; i++ {
		url := urlCode.GetUrl(code)
		if url.Url == "" {
			return code
		}
	}
	return ""
}

func genRandomCode() string {
	bytes := []byte("0123456789qazwsxedcrfvtgbyhnujmiklopQAZWSXEDCRFVTGBYHNUJMIKLOP")

	var str []byte
	for i := 0; i < 6; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		str = append(str, bytes[r.Intn(len(bytes))])
	}

	return string(str)
}

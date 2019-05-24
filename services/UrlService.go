package services

import (
	"errors"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"shorturl/commons"
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

	uc := urlCode.GetByUrl(url)

	var id int
	if uc.Code != "" {
		return code, nil
	} else if uc.Id != 0 && uc.Code == "" {
		id = uc.Id
	} else {
		id = urlCode.AddUrl(url)
	}

	if id == 0 {
		return "", errors.New("get id failed")
	}
	code = TransToCode(id)
	if code == "" {
		return code, errors.New("gen code failed")
	}

	logs.Info("add new short url, code: " + code)
	err = urlCode.UpdateCode(id, code)
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
		result := urlCode.GetByCode(code)
		if result.Url == "" {
			return "", errors.New("code not existed")
		}
		url = result.Url
		_ = bm.Put(code, url, time.Hour)
	}

	return url, nil
}

//把这个数字转换成62进制
func TransToCode(id int) string {
	bytes := []byte("0lv12NUJ3789qazwegbyhnujmipQAZWsxSXEDCR4kt56FVTGBYHMIodcrfKLOP")

	var code string

	for m := id; m > 0; m = m / 62 {
		n := m % 62
		code += string(bytes[n])
		if m < 62 {
			break
		}
	}
	return commons.ReverseString(code)
}

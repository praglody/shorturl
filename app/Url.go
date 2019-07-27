package app

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/hashicorp/golang-lru"
	"shorturl/commons"
	"shorturl/configs"
)

type Service interface {
	GenShortUrl(string, int) (string, error)
	RestoreUrl(string) (string, error)
}

type service struct {
	cCache   *lru.Cache
	uCache   *lru.Cache
	urlModel *UrlCode
}

func NewService() Service {
	var (
		cCache, _ = lru.New(10000)
		uCache, _ = lru.New(10000)
	)
	return &service{
		cCache:   cCache,
		uCache:   uCache,
		urlModel: &UrlCode{},
	}
}

// 生成短链接
func (s *service) GenShortUrl(url string, userId int) (shortUrl string, err error) {
	var shortCode string
	var urlMd5 = commons.MD5(url)
	if code, ok := s.uCache.Get(urlMd5); ok {
		shortCode = code.(string)
	} else {
		result := s.urlModel.GetByUrl(url)
		if result.Code != "" {
			shortCode = result.Code
		} else {
			var id = 0
			if result.Id != 0 {
				id = result.Id
			} else {
				id = s.urlModel.AddUrl(url, userId)
			}
			if id == 0 {
				return "", errors.New("get id failed")
			}
			shortCode = TransToCode(id)
			if shortCode == "" {
				return "", errors.New("gen code failed")
			}
			logs.Info("add new short url, code: " + shortCode)
			err = s.urlModel.UpdateCode(id, shortCode)
			if err != nil {
				return "", err
			}
		}
	}
	//add cache
	go func() {
		s.cCache.Add(shortCode, url)
		s.uCache.Add(urlMd5, shortCode)
	}()
	return configs.Conf.String("APP_URL") + shortCode, nil
}

// 还原url
func (s *service) RestoreUrl(code string) (string, error) {
	//get from cache
	var url string
	if c, ok := s.cCache.Get(code); ok {
		url = c.(string)
	} else {
		result := s.urlModel.GetByCode(code)
		if result.Url == "" {
			return "", errors.New("code not existed")
		} else {
			url = result.Url
		}
	}
	//add cache
	go func() {
		s.cCache.Add(code, url)
	}()
	//add click
	go func() {
		addClick <- code
	}()

	return url, nil
}

//把数字转换成62进制
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
	return code
}

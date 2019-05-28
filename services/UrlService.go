package services

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hashicorp/golang-lru"
	"github.com/jinzhu/gorm"
	"os"
	"shorturl/models"
	"time"
)

type UrlService struct{}

//cache code
var cCache, _ = lru.New(10000)

//cache url
var uCache, _ = lru.New(10000)

//db model
var urlCode = &models.UrlCode{}

//click channel
var addClick = make(chan string)

func init() {
	go Clicker()
}

//异步统计点击数
func Clicker() {
	go func() {
		var clicks = make(map[string]int, 100)
		for {
			c := <-addClick
			if c == "save" {
				go saveClicks(clicks, false)
				clicks = make(map[string]int, 100)
			} else if c == "shutdown" {
				go saveClicks(clicks, false)
			} else {
				clicks[c]++
				if len(clicks) > 1000 {
					go saveClicks(clicks, false)
					clicks = make(map[string]int, 100)
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			<-ticker.C
			addClick <- "save"
		}
	}()
}

func Shutdown() {
	addClick <- "shutdown"
}

func saveClicks(clicks map[string]int, shutdown bool) {
	for code, c := range clicks {
		var uc models.UrlCode
		models.DB.Where("code = ?", code).Find(&uc).UpdateColumn("click", gorm.Expr("click + ?", c))
		logs.Info(fmt.Sprintf("add %d click on %s", c, code))
	}
	if shutdown {
		os.Exit(1)
	}
}

func (UrlService) GenShortUrl(url string) (shortUrl string, err error) {
	var shortCode string

	var urlMd5 = models.MD5(url)

	if code, ok := uCache.Get(urlMd5); ok {
		shortCode = code.(string)
	} else {
		result := urlCode.GetByUrl(url)
		if result.Code != "" {
			shortCode = result.Code
		} else {
			var id = 0
			if result.Id != 0 {
				id = result.Id
			} else {
				id = urlCode.AddUrl(url)
			}
			if id == 0 {
				return "", errors.New("get id failed")
			}
			shortCode = TransToCode(id)
			if shortCode == "" {
				return "", errors.New("gen code failed")
			}
			logs.Info("add new short url, code: " + shortCode)
			err = urlCode.UpdateCode(id, shortCode)
			if err != nil {
				return "", err
			}
		}
	}
	//add cache
	go func() {
		cCache.Add(shortCode, url)
		uCache.Add(urlMd5, shortCode)
	}()

	return models.Conf.AppUrl + shortCode, nil
}

func (UrlService) RestoreUrl(code string) (string, error) {
	//get from cache
	var url string
	if c, ok := cCache.Get(code); ok {
		url = c.(string)
	} else {
		result := urlCode.GetByCode(code)
		if result.Url == "" {
			return "", errors.New("code not existed")
		} else {
			url = result.Url
		}
	}

	//add cache
	go func() {
		cCache.Add(code, url)
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

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
var cCode, _ = lru.New(10000)

//cache url
var cUrl, _ = lru.New(10000)

//db model
var urlCode = &models.UrlCode{}

//click channel
var ClickC = make(chan string)

func init() {
	//点击数异步入库
	go Clicker()
}

func Clicker() {
	go func() {
		var clicks = make(map[string]int, 100)
		for {
			c := <-ClickC
			if c == "shutdown" {
				go saveAndShutdown(clicks)
			} else if c == "save" {
				go saveToDB(clicks)
				clicks = make(map[string]int, 100)
			} else {
				clicks[c]++
				if len(clicks) > 1000 {
					go saveToDB(clicks)
					clicks = make(map[string]int, 100)
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			<-ticker.C
			ClickC <- "save"
		}
	}()
}

func saveToDB(clicks map[string]int) {
	for code, c := range clicks {
		var uc models.UrlCode
		models.DB.Where("code = ?", code).Find(&uc).UpdateColumn("click", gorm.Expr("click + ?", c))
		logs.Info(fmt.Sprintf("add %d click on %s", c, code))
	}
}

func saveAndShutdown(clicks map[string]int) {
	saveToDB(clicks)
	os.Exit(1)
}

func (UrlService) GenCode(url string) (code string, err error) {
	//get from cache
	if code, ok := cUrl.Get(models.MD5(url)); ok {
		return code.(string), nil
	}
	uc := urlCode.GetByUrl(url)
	var id int
	if uc.Code != "" {
		cUrl.Add(models.MD5(url), uc.Code)
		return uc.Code, nil
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
		return "", errors.New("gen code failed")
	}

	logs.Info("add new short url, code: " + code)
	err = urlCode.UpdateCode(id, code)
	if err != nil {
		return "", err
	}
	//cache
	go func() {
		cCode.Add(code, url)
		cUrl.Add(models.MD5(url), code)
	}()

	return code, nil
}

func (UrlService) RecCode(code string) (string, error) {
	//get from cache
	var url string
	c, _ := cCode.Get(code)
	if c != nil {
		url = c.(string)
	}
	if url == "" {
		result := urlCode.GetByCode(code)
		if result.Url == "" {
			return "", errors.New("code not existed")
		}
		url = result.Url
		//add cache
		go func() {
			cCode.Add(code, url)
		}()
	}

	go func() {
		ClickC <- code
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

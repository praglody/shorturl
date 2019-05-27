package services

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hashicorp/golang-lru"
	"shorturl/models"
	"sync"
	"time"
)

type UrlService struct{}

type Clicks struct {
	sync.RWMutex
	count map[string]int
}

//cache code
var cCode *lru.Cache

//cache url
var cUrl *lru.Cache

//db model
var urlCode *models.UrlCode

//click
var clicks *Clicks

func init() {
	cCode, _ = lru.New(10000)
	cUrl, _ = lru.New(10000)
	urlCode = &models.UrlCode{}

	startClicker()
}

//save click to db
func startClicker() {
	clicks = &Clicks{count: make(map[string]int)}
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			if len(clicks.count) > 0 {
				SaveClick()
			}
			<-ticker.C
		}
	}()
}

func SaveClick() {
	clicks.Lock()
	data := clicks.count
	clicks = &Clicks{count: make(map[string]int)} //create a new map
	for code, c := range data {
		logs.Info(fmt.Sprintf("add %d click on %s", c, code))
		models.DB.Where("code = ?", code).Find(&urlCode)
		urlCode.Click = urlCode.Click + c
		models.DB.Save(&urlCode)
	}
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
	cCode.Add(code, url)
	cUrl.Add(models.MD5(url), code)
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
		cCode.Add(code, url)
	}

	clicks.Lock()
	clicks.count[code]++
	clicks.Unlock()
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

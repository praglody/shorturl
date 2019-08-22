package app

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"shorturl/commons"
	"time"
)

type UrlCode struct {
	Id        int `gorm:"primary_key"`
	MD5       string
	Code      string
	Url       string
	Click     int
	UserId    int
	CreatedAt int
}

func (UrlCode) AddUrl(url string, userId int) (int, error) {
	var uc UrlCode
	uc.Url = url
	uc.Code = ""
	uc.MD5 = commons.MD5(url)
	uc.UserId = userId
	uc.CreatedAt = int(time.Now().Unix())
	err := commons.DB.Create(&uc).Error
	if err != nil {
		return 0, err
	}
	return uc.Id, nil
}

func (UrlCode) GetByUrl(url string) UrlCode {
	var result UrlCode
	commons.DB.Where("md5 = ?", commons.MD5(url)).Find(&result)
	return result
}

func (UrlCode) GetByCode(code string) UrlCode {
	var uc UrlCode
	commons.DB.Where("code = ?", code).First(&uc)
	return uc
}

func (UrlCode) UpdateCode(id int, code string) error {
	commons.DB.Table("url_codes").Where("id = ?", id).Update("code", code)
	if commons.DB.Error != nil {
		return commons.DB.Error
	}
	return nil
}

func (UrlCode) SaveClicks(clicks map[string]int) {
	for code, c := range clicks {
		var uc UrlCode
		commons.DB.Where("code = ?", code).Find(&uc).UpdateColumn("click", gorm.Expr("click + ?", c))
		logs.Info(fmt.Sprintf("add %d click on %s", c, code))
	}
}

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"shorturl/commons"
	"time"
)

type UrlCode struct {
	Id        int `gorm:"primary_key"`
	MD5       string
	Code      string
	Url       string
	CreatedAt int
}

//自定义表名
func (UrlCode) TableName() string {
	return "tbl_url_code"
}

func (uc UrlCode) GetCode(url string) UrlCode {

	md5 := commons.MD5(url)

	DB.Where("md5 = ?", md5).First(&uc)

	return uc
}

func (uc UrlCode) GetUrl(code string) UrlCode {

	DB.Where("code = ?", code).First(&uc)

	return uc
}

func (uc UrlCode) AddUrl(url string, code string) error {

	uc.Url = url
	uc.Code = code
	uc.MD5 = commons.MD5(url)
	uc.CreatedAt = int(time.Now().Unix())

	DB.Create(&uc)

	if DB.Error != nil {
		return DB.Error
	}

	return nil
}

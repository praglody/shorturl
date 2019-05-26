package models

import (
	_ "github.com/go-sql-driver/mysql"
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

func (uc UrlCode) AddUrl(url string) int {
	uc.Url = url
	uc.Code = ""
	uc.MD5 = MD5(url)
	uc.CreatedAt = int(time.Now().Unix())
	DB.Create(&uc)
	return uc.Id
}

func (uc UrlCode) GetByUrl(url string) UrlCode {
	var result UrlCode
	DB.Where("md5 = ?", MD5(url)).Find(&result)
	return result
}

func (uc UrlCode) GetByCode(code string) UrlCode {
	DB.Where("code = ?", code).First(&uc)
	return uc
}

func (uc UrlCode) UpdateCode(id int, code string) error {
	DB.Find(&uc, id)
	uc.Code = code
	DB.Save(&uc)
	if DB.Error != nil {
		return DB.Error
	}
	return nil
}

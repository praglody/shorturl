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
	Click     int
	UserId    int
	CreatedAt int
}

func (UrlCode) AddUrl(url string, userId int) int {
	var uc UrlCode
	uc.Url = url
	uc.Code = ""
	uc.MD5 = MD5(url)
	uc.UserId = userId
	uc.CreatedAt = int(time.Now().Unix())
	DB.Create(&uc)
	return uc.Id
}

func (UrlCode) GetByUrl(url string) UrlCode {
	var result UrlCode
	DB.Where("md5 = ?", MD5(url)).Find(&result)
	return result
}

func (UrlCode) GetByCode(code string) UrlCode {
	var uc UrlCode
	DB.Where("code = ?", code).First(&uc)
	return uc
}

func (UrlCode) UpdateCode(id int, code string) error {
	var uc UrlCode
	DB.Find(&uc, id)
	uc.Code = code
	DB.Save(&uc)
	if DB.Error != nil {
		return DB.Error
	}
	return nil
}

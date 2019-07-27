package app

import (
	"shorturl/commons"
)

type UserModel struct {
	Id        int `gorm:"primary_key"`
	Name      string
	AppId     string
	AppSecret string
	Status    int
}

func (UserModel) GetByAppId(appId string) UserModel {
	var user UserModel
	commons.DB.Where("app_id = ? and status = 0", appId).Find(&user)
	return user
}

func (UserModel) GetUsers() []UserModel {
	var users []UserModel
	commons.DB.Where("status = 0").Find(&users)
	return users
}

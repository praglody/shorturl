package models

type User struct {
	Id        int `gorm:"primary_key"`
	Name      string
	AppId     string
	AppSecret string
	Status    int
}

func (User) GetByAppId(appId string) User {
	var user User
	DB.Where("app_id = ? and status = 0", appId).Find(&user)
	return user
}

func (User) GetUsers() []User {
	var users []User
	DB.Where("status = 0").Find(&users)
	return users
}

package models

var AppConfig *AppConf

type AppConf struct {
	Env     string
	Port    string
	BaseUrl string
	WorkDir string
}

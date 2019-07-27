package configs

import (
	"github.com/astaxie/beego/config"
	"log"
	"os"
)

var Conf config.Configer

func init() {
	//配置文件
	var envFile = ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Panicf("conf file [%s]  not found!", envFile)
	}
	conf, err := config.NewConfig("ini", envFile)
	if err != nil {
		log.Panicf("parse conf file [%s] failed, err: %s", envFile, err.Error())
	}
	Conf = conf
	log.Println("init all config file success")
}

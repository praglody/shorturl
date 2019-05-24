package tests

import (
	"log"
	"os"
	"shorturl/models"
)

//加载框架配置
func init() {
	dir, _ := os.Getwd()
	file := dir + "/../env.ini"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Panicf("conf file [%s]  not found!", file)
	}
	models.Conf.InitConfig(file)
}

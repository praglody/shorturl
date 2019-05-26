package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"shorturl/models"
	"shorturl/routers"
)

func main() {
	if models.Conf.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	err := logs.SetLogger(logs.AdapterFile, `{"filename":"storage/logs/app.log"}`)

	if err != nil {
		log.Fatalln("Log init failed, error: " + err.Error())
	}

	routers.Route(r)

	err = r.Run(":" + models.Conf.AppPort)
	if err != nil {
		log.Fatalln("Server start failed, error: " + err.Error())
	}
}

func init() {
	dir, _ := os.Getwd()
	file := dir + "/.env"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Panicf("conf file [%s]  not found!", file)
	}
	models.Conf.InitConfig(file)
}
